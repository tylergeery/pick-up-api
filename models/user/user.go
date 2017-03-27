package user

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/pick-up-api/utils/auth"
	"github.com/pick-up-api/utils/messaging"
	"github.com/pick-up-api/utils/types"
	"github.com/pick-up-api/utils/validation"
	"golang.org/x/crypto/bcrypt"
)

const BCryptHashCost = 10

type User struct {
	Id           int64           `json:"id" db:"id"`
	Email        string          `json:"email" db:"email"`
	password     string          `json:"-" db:"password"` // omit
	Name         string          `json:"name" db:"name"`
	FacebookId   types.NullInt64 `json:"facebook_id,omitempty" db:"facebook_id"`
	Active       int             `json:"active" db:"is_active"`
	RefreshToken string          `json:"refresh_token,omitempty" db:"refresh_token"`
	AccessToken  string          `json:"refresh_token,omitempty" db:"access_token"`
	CreatedAt    string          `json:"created_at,omitempty" db:"-"`
	UpdatedAt    string          `json:"updated_at,omitempty" db:"-"`
}

/**
 * Save the new user in the DB
 */
func (u *User) Save() (int64, error) {
	columns, values := u.GetUserColumnStringAndValues(false, true)
	columns += ", password"
	values = append(values, u.GetPassword())

	id, err := UserInsert(columns, values)
	u.Id = id

	return id, err
}

/**
 * Update the user in the db
 */
func (u *User) Update() error {
	columns, values := u.GetUserColumnStringAndValues(false, false)

	return UserUpdateValues(u.Id, columns, values, true)
}

/**
 * Set the password for the current User
 */
func (u *User) SetPassword(pw string) {
	u.password = pw
}

/**
 * Get the user's password
 */
func (u *User) GetPassword() string {
	return u.password
}

/**
 * Get the saveable property and value stubs for a db insert/update
 */
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

/**
 * Build a new user from a user map
 */
func (u *User) Build(userPostData map[string][]string) error {
	var err error

	log.Println(userPostData)
	for k, v := range userPostData {
		switch k {
		case "name":
			name := v[0]

			if !validation.IsNonEmptyString(name) {
				err = errors.New(messaging.USER_NAME_EMPTY)
			} else {
				u.Name = name
			}
		case "password":
			var valid bool
			pw := v[0]

			if valid, err = validation.IsValidPassword(pw); valid {
				hash, hashError := bcrypt.GenerateFromPassword([]byte(pw), BCryptHashCost)

				if hashError != nil {
					err = errors.New(messaging.USER_PASSWORD_INVALID)
				} else {
					u.SetPassword(string(hash))
				}
			}
		case "email":
			email := v[0]

			if !validation.IsValidEmail(email) {
				err = errors.New(messaging.USER_EMAIL_INVALID)
			} else {
				u.Email = email
			}
		default:
			// only add properties that are user supplied,
			// all else can be ignored here
		}

		// return err as soon as one arises
		if err != nil {
			break
		}
	}

	return err
}

/**
 * Add a refresh token to a newly created User
 */
func (u *User) AddRefreshToken() {
	tokenString, err := auth.CreateUserToken(u.Id, 1, "refresh")

	if err == nil {
		u.RefreshToken = tokenString
	}
}

/**
 * Adds a shorter lived access token for a User
 */
func (u *User) AddAccessToken() {
	tokenString, err := auth.CreateUserToken(u.Id, 1, "access")

	if err == nil {
		u.AccessToken = tokenString
	}
}
