package validation

import (
    "errors"
    "github.com/asaskevich/govalidator"
)

/**
 * Minimum password length for new users
 */
const MinPasswordLength = 8

/**
 * Checks for non empty string
 */
func IsNonEmptyString(str string) bool {
    return (str != "")
}

/**
 * Checks for valid email
 */
func IsValidEmail(email string) bool {
    return (govalidator.IsEmail(email))
}

/**
 * Checks for valid password
 *
 * Requirements:
 *  - Must be at least MinPasswordLength characters
 */
func IsValidPassword(pw string) (bool, error) {
    var err error

    if !IsNonEmptyString(pw) {
        return false, errors.New("Password is blank.")
    }

    if !IsStringOfLength(pw, MinPasswordLength) {
        return false, errors.New("Password is too short.")
    }

    return true, err
}

/**
 * Checks if the string is of valid length (in characters)
 */
func IsStringOfLength(str string, length int) bool {
    runes := []rune(str)
    runeCount := len(runes)

    return (runeCount >= length)
}
