package auth

import (
	"os"

	"github.com/sirupsen/logrus"
)

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		logrus.Warn("Please set the JWT_SECRET environment variable")
		return []byte("notSecure")
	}

	return []byte(secret)
}
