package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type localJwtClaims struct{}

type ApiSecurityOptions struct {
	IssuerURL string   `json:"authority"`
	ClientID  string   `json:"client_id"`
	Audience  []string `json:"-"`
}

type CustomClaims struct {
	PreferredUsername string `json:"preferred_username"`
	EmailVerified     bool   `json:"email_verified"`
}

func (c *CustomClaims) Validate(ctx context.Context) error {
	return nil
}

func checkLocalJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(auth) > 8 && strings.ToLower(auth[:6]) == "bearer" {
			token := auth[7:]
			claims, err := ValidateToken(token)
			if err == nil {
				r = r.WithContext(context.WithValue(r.Context(), localJwtClaims{}, claims))
			} else {
				logrus.Debug(err.Error())
			}
		}
		next.ServeHTTP(w, r)
	})
}

type UserInfo struct {
	Username string
}
