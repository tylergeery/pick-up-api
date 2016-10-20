package models

import (
    "log"
    "errors"
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

func (u *User) SetPassword(pw string) {
    u.password = pw
}

func (u *User) GetPassword() string {
    return u.password
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

    if err == nil && user.Id == 0 {
        err = errors.New("User could not be found")
    }

    return user, err
}

func UserCreateProfile(userPostData map[string][]string) (User, error) {
    var id int64
    var user User
    var err error

    for k, v := range userPostData {
        switch k {
        case "name":
            name := v[0]

            if !validation.IsNonEmptyString(name) {
                errors.New("Name cannot be empty.")
                break
            } else {
                user.Name = name
            }
        case "password":
            pw := v[0]

            if !validation.IsValidEmail(pw) {
                errors.New("Email is not valid.")
                break
            } else {
                user.SetPassword(pw)
            }
        case "email":
            email := v[0]

            if !validation.IsValidEmail(email) {
                errors.New("Email is not valid.")
                break
            } else {
                user.Email = email
            }
        default:
            // TODO, probably ignore
        }
    }

    if err == nil {
        id, err = user.Save()

        if err == nil {
            user.Id = id
        }
    }

    return user, err
}
