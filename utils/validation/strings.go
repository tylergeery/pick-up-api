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
 * Password validation errors
 */
const PasswordEmpty string = "Password is blank."
const PasswordShort string = "Password is too short."

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
        return false, errors.New(PasswordEmpty)
    }

    if !IsStringOfLength(pw, MinPasswordLength) {
        return false, errors.New(PasswordShort)
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
