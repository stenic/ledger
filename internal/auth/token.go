package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/stenic/ledger/internal/pkg/client"
	"github.com/stenic/ledger/internal/pkg/users"
)

const tokenTTL = 15 * time.Minute
const localIssuer = "ledger"

// Create a struct that will be encoded to a JWT.
// We add jwt.RegisteredClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func GenerateToken(username string, callbacks ...func(*Claims)) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix milliseconds
			Issuer:    localIssuer,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   username,
			Audience:  jwt.ClaimStrings{localIssuer},
			ID:        uuid.New().String(),
		},
	}

	for _, cb := range callbacks {
		cb(&claims)
	}

	// Create the JWT string
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(getJWTSecret())
}

func RefreshToken(tknStr string) (string, error) {
	claims, err := ValidateToken(tknStr)
	if err != nil {
		return "", err
	}

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second {
		return "", fmt.Errorf(http.StatusText(http.StatusBadRequest))
	}

	// Now, create a new token for the current use, with a renewed expiration time
	claims.NotBefore = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(tokenTTL))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func ValidateToken(token string) (*Claims, error) {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there's a problem with the signing method")
		}
		return getJWTSecret(), nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, err
		}
		return nil, err
	}
	if !tkn.Valid {
		return nil, fmt.Errorf(http.StatusText(http.StatusUnauthorized))
	}

	switch claims.Issuer {
	case localIssuer:
		// Check if the user still exists
		if users.FindByUsername(claims.Username) == nil {
			return nil, fmt.Errorf("unknown user")
		}
	case client.LocalClientIssuer:
		// Check if the client still exists
		if len(client.FindByID(claims.RegisteredClaims.ID)) != 1 {
			return nil, fmt.Errorf("unknown client token - %s", claims.RegisteredClaims.ID)
		}
		if err := client.UpdateLastUsageByID(claims.RegisteredClaims.ID); err != nil {
			logrus.Warn("failed to update last_usage: " + err.Error())
		}
	default:
		return nil, fmt.Errorf("invalid issuer")
	}

	return claims, nil
}
