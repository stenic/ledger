package auth

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LedgerValidator interface {
	ValidateToken(authString string) error
	GetJWTMiddleware() gin.HandlerFunc
}

type ledgerValidator struct {
	IssuerURL string
	Audience  []string
}

func New(cfg ApiSecurityOptions) LedgerValidator {
	l := ledgerValidator{
		IssuerURL: cfg.IssuerURL,
		Audience:  cfg.Audience,
	}

	return l
}

func (l ledgerValidator) ValidateToken(authString string) error {
	if len(authString) < 8 && strings.ToLower(authString[:6]) != "bearer" {
		return fmt.Errorf("not a bearer auth string")
	}

	// Validate local token
	logrus.Trace("Checking Local")
	_, err := ValidateToken(authString[7:])
	if err == nil {
		return nil
	}

	// OIDC validation
	jwtOIDCValidator := getOidcValidatorFunc(l.IssuerURL, l.Audience)
	if jwtOIDCValidator == nil {
		return fmt.Errorf("token is invalid")
	}

	logrus.Trace("Checking OIDC")
	_, err = jwtOIDCValidator.ValidateToken(context.Background(), authString[7:])
	return err
}

func (l ledgerValidator) GetJWTMiddleware() gin.HandlerFunc {
	var jwtOidcMiddleware *jwtmiddleware.JWTMiddleware
	if l.IssuerURL != "" {
		jwtOidcMiddleware = getOidcValidator(l.IssuerURL, l.Audience)
	}

	return func(c *gin.Context) {
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			c.Request = r
			c.Next()
		}

		if jwtOidcMiddleware != nil {
			logrus.Debug("Checking oidc token")
			jwtOidcMiddleware.CheckJWT(handler).ServeHTTP(c.Writer, c.Request)
		}

		// Continue with local auth if OIDC is not detected
		if c.Request.Context().Value(jwtmiddleware.ContextKey{}) == nil {
			logrus.Debug("Checking local token")
			checkLocalJWT(handler).ServeHTTP(c.Writer, c.Request)
		}
	}
}

func TokenFromContext(ctx context.Context) *UserInfo {
	if os.Getenv("AUTH_INSECURE") == "yes" {
		logrus.Error("Skipping auth, hope you are in dev mode.")
		return &UserInfo{
			Username: "Insecure user",
		}
	}
	if raw, valid := ctx.Value(localJwtClaims{}).(*Claims); valid {
		return &UserInfo{
			Username: raw.Username,
		}
	}

	if raw, valid := ctx.Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims); valid {
		cc, ok := raw.CustomClaims.(*CustomClaims)
		if !ok {
			return nil
		}
		return &UserInfo{
			Username: cc.PreferredUsername,
		}
	}

	return nil
}
