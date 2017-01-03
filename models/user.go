package models

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/pick-up-api/utils/auth"
	"github.com/pick-up-api/utils/resources"
	"github.com/pick-up-api/utils/validation"
	"golang.org/x/crypto/bcrypt"
)

const BCryptHashCost = 10

type User struct {
	Id       int64  `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	password string `json:"-" db:"password"` // omit
	Name     string `json:"name" db:"name"`
	Active   int64  `json:"active" db:"is_active"`
	Token    string `json:"token,omitempty" db:"-"`
}

func (u *User) Save() (int64, error) {
	var id int64

	columns, values := u.GetUserColumnStringAndValues(false, true)
	values = append(values, u.GetPassword())
	db := resources.DB()
	query := fmt.Sprintf(
		`   INSERT INTO users (%s, password, created_at)
            VALUES (%s, NOW()) RETURNING id
        `, columns, resources.SqlStub(len(values)))

	err := db.QueryRow(query, values...).Scan(&id)

	u.Id = id

	return id, err
}

func (u *User) Update() error {
	columns, values := u.GetUserColumnStringAndValues(false, false)
	db := resources.DB()
	query := fmt.Sprintf(
		`   UPDATE users
            SET (%s) = (%s)
            WHERE id = %d
        `, columns, resources.SqlStub(len(values)), u.Id)

	_, err := db.Exec(query, values...)

	return err
}

func (u *User) SetPassword(pw string) {
	u.password = pw
}

func (u *User) GetPassword() string {
	return u.password
}

func (u *User) GetUserColumnStringAndValues(includeID, includePW bool) (string, []interface{}) {
	var columns string
	var values []interface{}

	userElem := reflect.ValueOf(u).Elem()
	userType := userElem.Type()

	// Iterate over fields
	for i := 0; i < userElem.NumField(); i++ {
		columnName := userType.Field(i).Tag.Get("db")

		if !includeID && columnName == "id" {
			continue
		}

		if columnName == "password" || columnName == "-" {
			continue
		}

		columns += fmt.Sprintf("%s, ", userType.Field(i).Tag.Get("db"))
		values = append(values, userElem.Field(i).Interface())
	}

	columns = strings.TrimRight(columns, ", ")

	return columns, values
}

func (u *User) Build(userPostData map[string][]string) error {
	var err error

	for k, v := range userPostData {
		switch k {
		case "name":
			name := v[0]

			if !validation.IsNonEmptyString(name) {
				err = errors.New("Name cannot be empty.")
			} else {
				u.Name = name
			}
		case "password":
			var valid bool
			pw := v[0]

			if valid, err = validation.IsValidPassword(pw); valid {
				hash, hashError := bcrypt.GenerateFromPassword([]byte(pw), BCryptHashCost)

				if hashError != nil {
					err = errors.New("Password is invalid.")
				} else {
					u.SetPassword(string(hash))
				}
			}
		case "email":
			email := v[0]

			if !validation.IsValidEmail(email) {
				err = errors.New("Email is invalid.")
			} else {
				u.Email = email
			}
		default:
			// TODO, probably ignore
		}

		// return err as soon as one arises
		if err != nil {
			break
		}
	}

	return err
}

func (u *User) AddToken() {
	tokenString, err := auth.CreateUserToken(u.Id, 1)

	if err == nil {
		u.Token = tokenString
	}
}

func UserGetById(id int64) (User, error) {
	var user User
	db := resources.DB()

	rows, err := db.Query("SELECT id, email, name, is_active FROM users WHERE id = $1", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id, active int64
		var email, name string

		if err := rows.Scan(&id, &email, &name, &active); err != nil {
			log.Fatal(err)
		}

		user = User{}
		user.Id = id
		user.Email = email
		user.password = "secret"
		user.Name = name
		user.Active = active
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if user.Id == 0 && err == nil {
		err = errors.New("User could not be found")
	}

	if user.Active != 1 {
		err = errors.New("User does not have an active account")
	}

	return user, err
}

func UserCreateProfile(userPostData map[string][]string) (User, error) {
	var user User

	err := user.Build(userPostData)
	user.Active = 1

	if err == nil {
		_, err = user.Save()
	}

	return user, err
}

func UserUpdateProfile(userId int64, userPostData map[string][]string) (User, error) {
	user, err := UserGetById(userId)

	if err == nil {
		err = user.Build(userPostData)
	}
	if err == nil {
		err = user.Update()
	}

	return user, err
}

func UserDeleteProfile(userId string) error {
	db := resources.DB()
	_, err := db.Query(
		`   UPDATE users
            SET is_active = 0, updated_at = NOW()
            WHERE id = $1
        `, userId)

	return err
}
