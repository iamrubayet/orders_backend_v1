package utils

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Define a struct to hold the JWT claims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT for the user after login
func GenerateJWT(username string) (string, error) {
	// Set the expiration time for the token (e.g., 5 days)
	expirationTime := time.Now().Add(5 * time.Hour)

	// Create the claims
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "yourAppName", // You can set your application name or another identifier here
		},
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key (you can store this key in an environment variable)
	secretKey := os.Getenv("JWT_SECRET_KEY") // It's better to load the key from the environment for security

	// Generate the signed token string
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Println("Error signing the token:", err)
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates the JWT token
func ValidateToken(tokenString string) (*Claims, error) {
	// Parse the token
	claims := &Claims{}
	secretKey := os.Getenv("JWT_SECRET_KEY")

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Check token expiration
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}
