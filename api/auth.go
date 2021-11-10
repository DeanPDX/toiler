package main

import (
	"errors"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

type JWTClaims struct {
	UserID    int       `json:"userID"`
	ExpiresAt time.Time `json:"expiresAt"`
	jwt.StandardClaims
}

func parseToken(tokenString string) (JWTClaims, error) {
	fmt.Println(tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Wrong signing method")
		}
		return []byte(globalConfig.SigningSecret), nil
	})

	if err != nil {
		return JWTClaims{}, err
	}

	claims, _ := token.Claims.(JWTClaims)
	return claims, nil
}
