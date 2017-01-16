package services

import (
	"errors"
	"log"

	"github.com/pick-up-api/models"
	"github.com/pick-up-api/utils/messaging"
	"github.com/pick-up-api/utils/resources"
	"github.com/pick-up-api/utils/types"
)

func UserGetById(id int64) (models.User, error) {
	var user models.User
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

		user = models.User{}
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

func UserGetByEmail(email string) (models.User, error) {
	var user models.User
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

		user = models.User{}
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

func UserCreateProfile(userPostData map[string][]string) (models.User, error) {
	var user models.User
	var userWithEmail models.User

	err := user.Build(userPostData)
	user.Active = 1

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

func UserUpdateProfile(userId int64, userPostData map[string][]string) (models.User, error) {
	user, err := UserGetById(userId)

	if err == nil {
		err = user.Build(userPostData)
	}
	if err == nil {
		err = user.Update()
	}

	return user, err
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
