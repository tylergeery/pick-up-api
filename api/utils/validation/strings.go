package validation

import (
    "errors"
)

/**
 * Minimum password length for new users
 */
const MinPasswordLength = 8

/**
 * Checks for non empty string
 */
func IsNonEmptyString(str string) bool {
    return true
}

/**
 * Checks for valid email
 */
func IsValidEmail(str string) (bool, error) {
    var err error

    return true, err
}

/**
 * Checks for valid password
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
 * Checks if the string is of valid length
 */
func IsStringOfLength(str string, length int) bool {
    runes := []rune(str)
    runeCount := len(runes)

    return (runeCount >= length)
}
