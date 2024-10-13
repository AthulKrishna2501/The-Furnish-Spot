package middleware

import (
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("Secret_key")

type Claims struct {
	jwt.RegisteredClaims
	Name string `json:"name"`
	Role string `json:"role"`
}

func GenerateToken(name, role string) (string, error) {
	claims := &Claims{
		Name: name,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // 24 hours expiration
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Fatal("Error creating Token", err)
	}

	return tokenString, nil
}
func ParseToken(tokenString string) (*Claims, error) {

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
