package models

import (
    "log"
    "database/sql"
	_ "github.com/lib/pq"
)

type User struct {
    id int
    email string
    name string
}

func UserGetById(id string) User {
    var user User
    db, err := sql.Open("postgres", "postgres://raccoon:pickEmUp@192.168.99.100/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT id, email, name FROM users")
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()
    for rows.Next() {
        var id int
        var email string
        var name string
        if err := rows.Scan(&id, &email, &name); err != nil {
            log.Fatal(err)
        }
        user = User{id, email, name}
    }
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }

    return user
}
