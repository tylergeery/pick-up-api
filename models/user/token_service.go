package user

import (
	"fmt"
	"strings"

	"github.com/pick-up-api/utils/resources"
)

func UserInsertRefreshToken(userId int64, tokenString string) error {
	values := []interface{}{
		userId,
		tokenString,
	}

	tx := resources.TX()
	defer tx.Rollback()
	query := fmt.Sprintf(
		`   INSERT INTO user_tokens
            (user_id, refresh_token, created_at)
            VALUES (%s, NOW())
        `, resources.SqlStub(len(values)))

	_, err := tx.Exec(query, values...)

	if err == nil {
		tx.Commit()
	}

	return err
}

func userUpdateTokenValues(userId int64, columns string, values []interface{}) error {
	tx := resources.TX()
	defer tx.Rollback()
	query := fmt.Sprintf(
		`   UPDATE user_tokens
            SET (%s, updated_at) = (%s, NOW())
            WHERE user_id = %d
        `, columns, resources.SqlStub(len(values)), userId)

	_, err := tx.Exec(query, values...)

	if err == nil {
		tx.Commit()
	}

	return err
}

func UserUpdateTokens(userId int64, accessToken string, refreshToken string) error {
	cols := ""
	vals := []interface{}{}

	if accessToken != "" {
		cols += "access_token,"
		vals = append(vals, accessToken)
	}

	if refreshToken != "" {
		cols += "refresh_token,"
		vals = append(vals, refreshToken)
	}

	cols = strings.TrimRight(cols, ", ")

	return userUpdateTokenValues(userId, cols, vals)
}
