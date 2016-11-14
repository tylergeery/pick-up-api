package pickUpValidationTests

import (
    "fmt"
    "utf8"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/pick-up-api/utils/auth"
)

func TestUserTokens(t *testing.T) {
    // test strings
    userId := 1234
    permissions := 1

    tokenString, tokenErr := auth.CreateUserToken()

    assert.True(t, tokenErr == nil, "There was an error creating user token")
    assert.True(t, utf8.RuneCountInString(tokenString) > 10, "Token is not of valid length")
}
