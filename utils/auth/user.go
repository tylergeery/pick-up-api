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

    claims, success := ExtractToken(stringToken)

    if !success || claims["userId"] == nil {
        valid = false
    } else {
        userIdFromClaims := int64(claims["userId"].(float64))

        if (claims["type"] != "user") {
            valid = false
        }

        if (userIdFromClaims != userId) {
            valid = false
        }
    }

    return valid
}

func GetUserIdFromToken(stringToken string) int64, bool {
    claims, success := ExtractToken(stringToken)

    if !success || claims["userId"] == nil {
        return 0, false
    } else {
        userIdFromClaims := int64(claims["userId"].(float64))

        return userIdFromClaims, true
    }
}
