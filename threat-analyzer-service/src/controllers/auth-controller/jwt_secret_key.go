package authcontroller

import (
	"log"
	"os"
)

func JWTSecretKey() []byte {
	key := os.Getenv("JWT_SECRET_KEY")
	if key == "" {
		log.Fatal("JWT_SECRET_KEY is not set in environment variables")
	}
	if len(key) < 32 {
		log.Fatal("JWT_SECRET_KEY must be at least 32 characters long for HS256")
	}
	return []byte(key)
}
