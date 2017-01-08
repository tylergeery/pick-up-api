package auth

import (
	"fmt"
	"log"

	jwt "github.com/dgrijalva/jwt-go"
)

var signingKey []byte = []byte(PICKUP_AUTH_TOKEN_SECRET)

func CreateToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(signingKey)

	if err != nil {
		log.Print("JWT CreateToken Err: %s", err)
	}

	return tokenString, err
}

func ExtractToken(tokenString string) (jwt.MapClaims, bool) {
	var claims jwt.MapClaims
	var ok bool

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return signingKey, nil
	})

	if err == nil && token.Claims != nil {
		claims, ok = token.Claims.(jwt.MapClaims)
	}

	return claims, ok
}
