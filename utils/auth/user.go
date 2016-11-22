package auth

import (
    "time"
    "github.com/dgrijalva/jwt-go"
)

func CreateUserToken(userId int, permissions int) (string, error) {
    claims := jwt.MapClaims{
        "userId": userId,
        "permissions": permissions,
        "ts": time.Now().Unix(),
        "type": "user"}

    return CreateToken(claims)
}

func ValidateUserToken(stringToken string, userId int) bool {
    var valid bool = true

    claims, _ := ExtractToken(stringToken)

    if (claims["type"] != "user") {
        valid = false
    }

    if (claims["userId"] != userId) {
        valid = false
    }

    return valid
}
