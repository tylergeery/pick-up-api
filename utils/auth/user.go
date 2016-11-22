package auth

import (
    "time"
    "github.com/dgrijalva/jwt-go"
)

func CreateUserToken(userId int64, permissions int) (string, error) {
    claims := jwt.MapClaims{
        "userId": userId,
        "permissions": permissions,
        "ts": time.Now().Unix(),
        "type": "user"}

    return CreateToken(claims)
}

func ValidateUserToken(stringToken string, userId int64) bool {
    var valid bool = true

    claims, _ := ExtractToken(stringToken)
    userIdFromClaims := int64(claims["userId"].(float64))

    if (claims["type"] != "user") {
        valid = false
    }

    if (userIdFromClaims != userId) {
        valid = false
    }

    return valid
}
