package main

import (
	"errors"
	"log"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

// hmacSigningSecret a global var for our signing secret.
var hmacSigningSecret []byte

// JWTInit must be called to initialize our signing secret. Can only
// be called during app initialization because not thread-safe.
func initializeJWT(signingSecret string) {
	hmacSigningSecret = []byte(signingSecret)
}

// ensureInit ensures that the caller initialized a signing secret.
// Be sure to call this before dealing with anything that requires
// signing secret to be present.
func ensureInit() {
	if len(hmacSigningSecret) == 0 {
		log.Fatal("You must call initializeJWT with a valid signing secret")
	}
}

// UserClaims is a struct for storing our users' claims.
type UserClaims struct {
	UserID    int       `json:"userID"`
	ExpiresAt time.Time `json:"expiresAt"`
	jwt.StandardClaims
}

// generateToken generates a token with our custom claims and returns a signed string.
func generateToken(userID int) (string, error) {
	ensureInit()
	// Create the Claims
	claims := UserClaims{
		userID,
		time.Now().Add(48 * time.Hour),
		jwt.StandardClaims{},
	}
	// Create a token and return a signed string
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(hmacSigningSecret)
}

// parseToken parses a signed string and returns our custom claims.
func parseToken(tokenString string) (*UserClaims, error) {
	ensureInit()
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Make sure we are using the correct signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Wrong signing method")
		}
		return hmacSigningSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("Token is invalid.")
}
