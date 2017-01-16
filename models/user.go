package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/pick-up-api/utils/auth"
	"github.com/pick-up-api/utils/resources"
	"github.com/pick-up-api/utils/types"
	"github.com/pick-up-api/utils/validation"
	"golang.org/x/crypto/bcrypt"
)

const BCryptHashCost = 10

type User struct {
	Id         int64           `json:"id" db:"id"`
	Email      string          `json:"email" db:"email"`
	password   string          `json:"-" db:"password"` // omit
	Name       string          `json:"name" db:"name"`
	FacebookId types.NullInt64 `json:"facebook_id,omitempty" db:"facebook_id"`
	Active     int             `json:"active" db:"is_active"`
	Token      string          `json:"token,omitempty" db:"-"`
	CreatedAt  string          `json:"created_at,omitempty" db:"-"`
	UpdatedAt  string          `json:"updated_at,omitempty" db:"-"`
}

func (u *User) Save() (int64, error) {
	var id int64

	columns, values := u.GetUserColumnStringAndValues(false, true)
	values = append(values, u.GetPassword())
	tx := resources.TX()
	defer tx.Rollback()
	query := fmt.Sprintf(
		`   INSERT INTO users (%s, password, created_at)
            VALUES (%s, NOW()) RETURNING id
        `, columns, resources.SqlStub(len(values)))

	err := tx.QueryRow(query, values...).Scan(&id)

	if err == nil {
		u.Id = id
		tx.Commit()
	}

	return id, err
}

func (u *User) Update() error {
	columns, values := u.GetUserColumnStringAndValues(false, false)
	tx := resources.TX()
	query := fmt.Sprintf(
		`   UPDATE users
            SET (%s, updated_at) = (%s, NOW())
            WHERE id = %d
        `, columns, resources.SqlStub(len(values)), u.Id)

	_, err := tx.Exec(query, values...)

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
