package auth

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestUserTokens(t *testing.T) {
	// test strings
	const userId int64 = 1234
	const userIdTwo int64 = 5678123412341234
	permissions := 1
	permissionsTwo := 3

	tokenString, tokenErr := CreateUserToken(userId, permissions, "access")
	tokenStringTwo, tokenErrTwo := CreateUserToken(userIdTwo, permissionsTwo, "refresh")

	// assert create user tokens
	assert.True(t, tokenErr == nil, "There was an error creating user token")
	assert.True(t, utf8.RuneCountInString(tokenString) > 100, "Token is not of valid length")
	assert.True(t, tokenErrTwo == nil, "There was an error creating user token")
	assert.True(t, utf8.RuneCountInString(tokenStringTwo) > 100, "Token is not of valid length")

	// validate user tokens
	assert.True(t, ValidateUserToken(tokenString, userId), "Token could not validate user")
	assert.True(t, ValidateUserToken(tokenStringTwo, userIdTwo), "Token could not validate user")

	// assert invalid
	assert.False(t, ValidateUserToken(tokenStringTwo, userId), "Token should have not validated user")
	assert.False(t, ValidateUserToken(tokenString, userIdTwo), "Token should have not validated user")
	assert.False(t, ValidateUserToken(tokenString, 0), "Token should have not validated user")
}
