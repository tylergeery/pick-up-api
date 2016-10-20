package connections

import (
    "log"
    "database/sql"
	_ "github.com/lib/pq"
)

/**
 * Get a SQL DB connection
 */
func DB() *sql.DB {
    db, err := sql.Open("postgres", "postgres://raccoon:pickEmUp@192.168.99.100/pickup?sslmode=disable")

    if err != nil {
        log.Fatal(err)
    }

    return db
}