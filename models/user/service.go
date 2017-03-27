package user

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/pick-up-api/utils/messaging"
	"github.com/pick-up-api/utils/resources"
	"github.com/pick-up-api/utils/types"
)

func UserGetById(id int64) (User, error) {
	var user User
	tx := resources.TX()

	rows, err := tx.Query("SELECT id, email, name, facebook_id, is_active, created_at, updated_at FROM users WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()
	for rows.Next() {
		var id int64
		var facebook_id types.NullInt64
		var email, name, created_at, updated_at string
		var is_active int

		if err := rows.Scan(&id, &email, &name, &facebook_id, &is_active, &created_at, &updated_at); err != nil {
			log.Fatal(err)
		}

		user = User{}
		user.Id = id
		user.Email = email
		user.Name = name
		user.FacebookId = facebook_id
		user.Active = is_active
		user.CreatedAt = created_at
		user.UpdatedAt = updated_at
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if user.Id == 0 && err == nil {
		err = errors.New(messaging.USER_NOT_FOUND)
	}

	if err == nil && user.Active != 1 {
		err = errors.New(messaging.USER_NOT_ACTIVE)
	}

	return user, err
}

func UserGetByEmail(email string) (User, error) {
	var user User
	tx := resources.TX()

	rows, err := tx.Query("SELECT id, email, name, facebook_id, is_active, created_at, updated_at FROM users WHERE email = $1", email)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int64
		var facebook_id types.NullInt64
		var email, name, created_at, updated_at string
		var is_active int

		if err := rows.Scan(&id, &email, &name, &facebook_id, &is_active, &created_at, &updated_at); err != nil {
			log.Fatal(err)
		}

		user = User{}
		user.Id = id
		user.Email = email
		user.Name = name
		user.FacebookId = facebook_id
		user.Active = is_active
		user.CreatedAt = created_at
		user.UpdatedAt = updated_at
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if user.Id == 0 && err == nil {
		err = errors.New(messaging.USER_NOT_FOUND)
	}

	if err == nil && user.Active != 1 {
		err = errors.New(messaging.USER_NOT_FOUND)
	}

	return user, err
}

func UserCreateProfile(userPostData map[string][]string) (User, error) {
	var user User
	var userWithEmail User

	err := user.Build(userPostData)
	user.Active = 1
	user.AddRefreshToken()

	if err == nil {
		userWithEmail, _ = UserGetByEmail(user.Email)

		if userWithEmail.Id != 0 {
			err = errors.New(messaging.USER_EMAIL_EXISTS)
		}
	}

	if err == nil {
		_, err = user.Save()
	}

	return user, err
}

func UserInsert(columns string, values []interface{}) (int64, error) {
	var id int64

	tx := resources.TX()
	defer tx.Rollback()
	query := fmt.Sprintf(
		`   INSERT INTO users (%s, created_at)
            VALUES (%s, NOW()) RETURNING id
        `, columns, resources.SqlStub(len(values)))

	err := tx.QueryRow(query, values...).Scan(&id)

	if err == nil {
		tx.Commit()
	}

	return id, err
}

func UserUpdateValues(userId int64, columns string, values []interface{}, updated_at bool) error {
	updated_at_col := ""
	updated_at_val := ""
	if updated_at {
		updated_at_col = ", updated_at"
		updated_at_val = ", NOW()"
	}

	tx := resources.TX()
	defer tx.Rollback()
	query := fmt.Sprintf(
		`   UPDATE users
            SET (%s%s) = (%s%s)
            WHERE id = %d
        `, columns, updated_at_col,
		resources.SqlStub(len(values)), updated_at_val,
		userId)

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
	}

	if refreshToken != "" {
		cols += "refresh_token,"
	}

	cols = strings.TrimRight(cols, ", ")

	return UserUpdateValues(userId, cols, vals, false)
}

func UserDeleteProfile(userId int64) error {
	tx := resources.TX()
	defer tx.Rollback()
	stmt, err := tx.Prepare(
		`   UPDATE users
            SET is_active = 0, updated_at = NOW()
            WHERE id = $1;`)

	if err == nil {
		defer stmt.Close()
		_, err = stmt.Exec(userId)
	}

	if err == nil {
		err = tx.Commit()
	}

	return err
}
