package models

import (
    "log"
    "errors"
    "reflect"
    "github.com/pick-up-api/utils/connections"
    "github.com/pick-up-api/utils/validation"
)

type User struct {
    Id int64 `json:"id"`
    Email string `json:"email"`
    password string `json:"-"` // omit
    Name string `json:"name"`
}

func (u *User) Save() (int64, error) {
    var id int64

    db := connections.DB()

    err := db.QueryRow(
        `INSERT INTO users (name, email, password)
        VALUES ($1, $2, $3) RETURNING id`,
        u.Name, u.Email, u.GetPassword()).Scan(&id)

    u.Id = id

    return id, err
}

func (u *User) Update() error {
    db := connections.DB()

    _, err := db.Exec(
        `INSERT INTO users (name, email, password)
        VALUES ($1, $2, $3)`,
        u.Name, u.Email, u.GetPassword())

    return err
}

func (u *User) SetPassword(pw string) {
    u.password = pw
}

func (u *User) GetPassword() string {
    return u.password
}

func (u *User) GetUserColumnStringAndValues() (string, []interface{}) {
    var columns string
    var values []interface{}

    userValue := reflect.ValueOf(u)
    userElem := userValue.Elem()

    // Iterate over fields
    for i := 0; i < userElem.NumField(); i++ {
        fieldKey := userElem.Field(i)
        columns += " " + ""
        values = append(values, fieldKey)
    }

    return columns, values
}

func createUser(userPostData map[string][]string) (User, error) {
    var err error
    var user User

    for k, v := range userPostData {
        switch k {
        case "name":
            name := v[0]

            if !validation.IsNonEmptyString(name) {
                err = errors.New("Name cannot be empty.")
                break
            } else {
                user.Name = name
            }
        case "password":
            pw := v[0]

            if !validation.IsValidEmail(pw) {
                err = errors.New("Email is not valid.")
                break
            } else {
                user.SetPassword(pw)
            }
        case "email":
            email := v[0]

            if !validation.IsValidEmail(email) {
                err = errors.New("Email is not valid.")
                break
            } else {
                user.Email = email
            }
        default:
            // TODO, probably ignore
        }
    }

    return user, err
}

func UserGetById(id string) (User, error) {
    var user User
    db := connections.DB()

	rows, err := db.Query("SELECT id, email, name FROM users WHERE id = $1", id)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        var id int64
        var email, name string

        if err := rows.Scan(&id, &email, &name); err != nil {
            log.Fatal(err)
        }

        user = User{id, email, "secret", name}
    }

    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }

    if user.Id == 0 && err == nil {
        err = errors.New("User could not be found")
    }

    return user, err
}

func UserCreateProfile(userPostData map[string][]string) (User, error) {
    var id int64

    user, err := createUser(userPostData)

    if err == nil {
        id, err = user.Save()

        if err == nil {
            user.Id = id
        }
    }

    return user, err
}

func UserUpdateProfile(userPostData map[string][]string) (User, error) {
    user, err := createUser(userPostData)

    if err == nil {
        err = user.Update()
    }

    return user, err
}

func UserDeleteProfile(userId string) error {
    db := connections.DB()

    _, err := db.Query("UPDATE users SET is_active = 0 WHERE id = $1", userId)

    if err != nil {
        log.Fatal(err)
    }

    return err
}
