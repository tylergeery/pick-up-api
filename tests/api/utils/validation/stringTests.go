package validationTests

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
    utfString := "ўСЋР»РїР°РЅС‹"

    // test Empty strings
    assert.False(t, validation.IsNonEmptyString(nilString), "This string is nil and empty")
    assert.False(t, validation.IsNonEmptyString(emptyString), "This string is empty")

    // Test Non-Empty strings
    assert.True(t, validation.IsNonEmptyString(spaceString), "This string is not empty")
    assert.True(t, validation.IsNonEmptyString(spaceWithTextString), "This string is not empty")
    assert.True(t, validation.IsNonEmptyString(textString), "This string is not empty")
    assert.True(t, validation.IsNonEmptyString(utfString), "This string is not empty")
}

func TestIsValidEmail(t *testing.T) {
    // TODO
}

func TestIsValidPassword(t *testing.T) {
    // TODO
}

func TestIsStringOfLength(t *testing.T) {
    // TODO
}
