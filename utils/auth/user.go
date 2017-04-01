package auth

import (
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

const ACCESS_TOKEN_LIFETIME = 1200

func CreateUserToken(userId int64, permissions int, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"userId":      userId,
		"permissions": permissions,
		"ts":          time.Now().Unix(),
		"type":        "user",
		"tokenType":   tokenType}

	return CreateToken(claims)
}

func ValidateUserToken(stringToken string, userId int64) bool {
	var valid bool = true

	claims, success := ExtractToken(stringToken)

	if !success || claims["userId"] == nil {
		valid = false
	} else {
		userIdFromClaims := int64(claims["userId"].(float64))

		if claims["type"] != "user" {
			valid = false
		}

		if userIdFromClaims != userId {
			valid = false
		}
	}

	return valid
}

/**
 * Validate an access token and return userId on success
 */
func GetUserIdFromAccessToken(stringToken string) (int64, bool) {
	claims, success := ExtractToken(stringToken)

	if !success {
		return 0, false
	}

	tokenTime := int64(claims["ts"].(float64))
	if (tokenTime + ACCESS_TOKEN_LIFETIME) < time.Now().Unix() {
		// access token is expired
		return 0, false
	}

	return getUserIdFromClaims(claims, "access")
}

/**
 * Validate a refresh token and return userId on success
 */
func GetUserIdFromRefreshToken(stringToken string) (int64, bool) {
	claims, success := ExtractToken(stringToken)

	if !success {
		return 0, false
	}

	return getUserIdFromClaims(claims, "refresh")
}

/**
 * Get the user id from claims for a given token
 */
func getUserIdFromClaims(claims jwt.MapClaims, tokenType string) (int64, bool) {
	if claims["userId"] == nil || claims["tokenType"].(string) != tokenType {
		return 0, false
	} else {
		userIdFromClaims := int64(claims["userId"].(float64))

		return userIdFromClaims, true
	}
}
