package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
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

func JwtHandler(opts ApiSecurityOptions) gin.HandlerFunc {

	var jwtOidcMiddleware *jwtmiddleware.JWTMiddleware
	if opts.IssuerURL != "" {
		jwtOidcMiddleware = getOidcValidator(opts.IssuerURL, opts.Audience)
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

func TokenFromContext(ctx context.Context) *UserInfo {
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
