package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.False(t, IsNonEmptyString(nilString), "This string is nil and empty")
	assert.False(t, IsNonEmptyString(emptyString), "This string is empty")

	// Test Non-Empty strings
	assert.True(t, IsNonEmptyString(spaceString), "This string is not empty")
	assert.True(t, IsNonEmptyString(spaceWithTextString), "This string is not empty")
	assert.True(t, IsNonEmptyString(textString), "This string is not empty")
	assert.True(t, IsNonEmptyString(unicodeString), "This string is not empty")
}

func TestIsValidEmail(t *testing.T) {
	var nilString string
	var unicodeString string = "ўСЋР»РїР°РЅС‹"
	var fakeEmail string = "tyler.geery@"
	var validEmail string = "test@yahoo.com"
	var validNonComEmail string = "reply.test@elastic.co"

	// Test Invalid Emails
	assert.False(t, IsValidEmail(nilString), "This string is nil and empty")
	assert.False(t, IsValidEmail(unicodeString), "This string is not a valid email")
	assert.False(t, IsValidEmail(fakeEmail), "This string is not a valid email")

	// Test Valid Emails
	assert.True(t, IsValidEmail(validEmail), "This string is a valid .com email")
	assert.True(t, IsValidEmail(validNonComEmail), "This string is a valid .co email")
}

func TestIsValidPassword(t *testing.T) {
	var nilPW string
	var shortPW string = "test"
	var unicodeShortPW string = "ўСЋР»Р"
	var validPW string = "testing!#@#123"
	var unicodeValidPW string = "ўСЋР»РїР°РЅС‹"

	// Execute password validations
	nilFailure, nilError := IsValidPassword(nilPW)
	shortFailure, shortError := IsValidPassword(shortPW)
	unicodeShortFailure, unicodeShortError := IsValidPassword(unicodeShortPW)
	validSuccess, validErr := IsValidPassword(validPW)
	unicodeValidSuccess, unicodeValidErr := IsValidPassword(unicodeValidPW)

	// Test Invalid Passwords
	assert.False(t, nilFailure, "This password is nil and empty")
	assert.False(t, shortFailure, "This password is not of valid length")
	assert.False(t, unicodeShortFailure, "This password is not of valid length")

	// Test Appropriate Errors
	assert.Equal(t, "Password is blank.", nilError.Error(), "Expected an empty password error")
	assert.Equal(t, "Password is too short.", shortError.Error(), "Expected a password too short error")
	assert.Equal(t, "Password is too short.", unicodeShortError.Error(), "Expected a password too short error")

	// Test Valid Emails
	assert.True(t, validSuccess, "This is a valid password")
	assert.True(t, unicodeValidSuccess, "This is a valid unicode password")

	// Test Errors are nil
	assert.True(t, validErr == nil, "Valid password should have no associated error")
	assert.True(t, unicodeValidErr == nil, "Unicode valid password should have no error")
}

func TestIsStringOfLength(t *testing.T) {
	var nilString string
	var shortString string = "test"
	var unicodeString string = "ўСЋР»Р"

	// Test Invalid Calls
	assert.False(t, IsStringOfLength(nilString, 1), "Nil string does not have length of 1")
	assert.False(t, IsStringOfLength(shortString, 5), "Short string does not have length of 5")
	assert.False(t, IsStringOfLength(unicodeString, 7), "Unicode string does not have length of 7")

	// Test Valid Calls
	for i := 0; i < 5; i++ {
		assert.True(t, IsStringOfLength(shortString, i), "Short string does have this length")
	}
	for i := 0; i < 7; i++ {
		assert.True(t, IsStringOfLength(unicodeString, i), "Unicode string does have this length")
	}
}
