package pickUpValidationTests

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/pick-up-api/utils/validation"
)

func TestIsNonEmptyString(t *testing.T) {
    // test strings
    var nilString string
    emptyString := ""
    spaceString := " "
    spaceWithTextString := " Farmer John"
    textString := "the great potato war"
    unicodeString := "ўСЋР»РїР°РЅС‹"

    // test Empty strings
    assert.False(t, validation.IsNonEmptyString(nilString), "This string is nil and empty")
    assert.False(t, validation.IsNonEmptyString(emptyString), "This string is empty")

    // Test Non-Empty strings
    assert.True(t, validation.IsNonEmptyString(spaceString), "This string is not empty")
    assert.True(t, validation.IsNonEmptyString(spaceWithTextString), "This string is not empty")
    assert.True(t, validation.IsNonEmptyString(textString), "This string is not empty")
    assert.True(t, validation.IsNonEmptyString(unicodeString), "This string is not empty")
}

func TestIsValidEmail(t *testing.T) {
    var nilString string
    var unicodeString string = "ўСЋР»РїР°РЅС‹"
    var fakeEmail string = "tyler.geery@"
    var validEmail string = "test@yahoo.com"
    var validNonComEmail string = "reply.test@elastic.co"

    // Test Invalid Emails
    assert.False(t, validation.IsValidEmail(nilString), "This string is nil and empty")
    assert.False(t, validation.IsValidEmail(unicodeString), "This string is not a valid email")
    assert.False(t, validation.IsValidEmail(fakeEmail), "This string is not a valid email")

    // Test Valid Emails
    assert.True(t, validation.IsValidEmail(validEmail), "This string is a valid .com email")
    assert.True(t, validation.IsValidEmail(validNonComEmail), "This string is a valid .co email")
}

func TestIsValidPassword(t *testing.T) {
    var nilPW string
    var shortPW string = "test"
    var unicodeShortPW string = "ўСЋР»Р"
    var validPW string = "testing!#@#123"
    var unicodeValidPW string = "ўСЋР»РїР°РЅС‹"

    // Execute password validations
    nilFailure, nilError := validation.IsValidPassword(nilPW)
    shortFailure, shortError := validation.IsValidPassword(shortPW)
    unicodeShortFailure, unicodeShortError := validation.IsValidPassword(unicodeShortPW)
    validSuccess, validErr := validation.IsValidPassword(validPW)
    unicodeValidSuccess, unicodeValidErr := validation.IsValidPassword(unicodeValidPW)

    // Test Invalid Passwords
    assert.False(t, nilFailure, "This password is nil and empty")
    assert.False(t, shortFailure, "This password is not of valid length")
    assert.False(t, unicodeShortFailure, "This password is not of valid length")

    // Test Appropriate Errors
    assert.Equal(t, validation.PasswordEmpty, nilError.Error(), "Expected an empty password error")
    assert.Equal(t, validation.PasswordShort, shortError.Error(), "Expected a password too short error")
    assert.Equal(t, validation.PasswordShort, unicodeShortError.Error(), "Expected a password too short error")

    // Test Valid Emails
    assert.True(t, validSuccess, "This is a valid password")
    assert.True(t, unicodeValidSuccess, "This is a valid unicode password")

    // Test Errors are nil
    assert.True(t, validErr == nil, "Valid password should have no associated error")
    assert.True(t, unicodeValidErr == nil, "Unicode valid password should have no error")
}

func TestIsStringOfLength(t *testing.T) {
    // TODO
}
