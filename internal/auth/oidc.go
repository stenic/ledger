package auth

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	log "github.com/sirupsen/logrus"
)

func getOidcValidator(issuerURLString string, audience []string) *jwtmiddleware.JWTMiddleware {
	issuerURL, err := url.Parse(issuerURLString)
	if err != nil {
		log.Fatal(err)
	}

	provider := jwks.NewCachingProvider(
		issuerURL,
		time.Duration(5*time.Minute),
		// func(p *jwks.Provider) {
		// 	if false {
		// 		return
		// 	}
		// 	tr := &http.Transport{
		// 		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
		// 	}
		// 	p.Client.Transport = tr
		// },
	)

	customClaims := func() validator.CustomClaims {
		return &CustomClaims{}
	}

	jwtValidator, _ := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		audience,
		// validator.WithAllowedClockSkew(30*time.Second),

		validator.WithCustomClaims(customClaims),
	)

	errorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		if errors.Is(err, jwtmiddleware.ErrJWTInvalid) {
			log.Debug(err)
		} else {
			jwtmiddleware.DefaultErrorHandler(w, r, err)
		}
	}

	return jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(errorHandler),
		jwtmiddleware.WithCredentialsOptional(true),
	)
}
