package pickUpValidationTests

import (
    "unicode/utf8"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/pick-up-api/utils/auth"
)

func TestUserTokens(t *testing.T) {
    // test strings
    userId := 1234
    userIdTwo := 5678123412341234
    permissions := 1
    permissionsTwo := 3

    tokenString, tokenErr := auth.CreateUserToken(userId, permissions)
    tokenStringTwo, tokenErrTwo := auth.CreateUserToken(userIdTwo, permissionsTwo)

    // assert create user tokens
    assert.True(t, tokenErr == nil, "There was an error creating user token")
    assert.True(t, utf8.RuneCountInString(tokenString) > 100, "Token is not of valid length")
    assert.True(t, tokenErrTwo == nil, "There was an error creating user token")
    assert.True(t, utf8.RuneCountInString(tokenStringTwo) > 100, "Token is not of valid length")

    // validate user tokens
    assert.True(t, auth.ValidateUserToken(tokenString, userId), "Token could not validate user")
    assert.True(t, auth.ValidateUserToken(tokenStringTwo, userIdTwo), "Token could not validate user")

    // assert invalid
    assert.False(t, auth.ValidateUserToken(tokenStringTwo, userId), "Token should have not validated user")
    assert.False(t, auth.ValidateUserToken(tokenString, userIdTwo), "Token should have not validated user")
    assert.False(t, auth.ValidateUserToken(tokenString, 0), "Token should have not validated user")
}
