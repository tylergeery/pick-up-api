package auth

import (
    "errors"
    "log"
    "github.com/dgrijalva/jwt-go"
)

func CreateToken(claims jwt.MapClaims) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Sign and get the complete encoded token as a string using the secret
    tokenString, err := token.SignedString(PICKUP_AUTH_TOKEN_SECRET)

    //TODO: log errors

    return tokenString, err
}

func ExtractToken(tokenString string) (jwt.MapClaims, bool) {
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Don't forget to validate the alg is what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
        }

        return PICKUP_AUTH_TOKEN_SECRET, nil
    })

    if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
        return claims, ok
    }
}
