package resources

import (
    "os"
    "fmt"
    "log"
    "strconv"
    "strings"
    "database/sql"
	_ "github.com/lib/pq"
)

/**
 * Get a SQL DB connection
 */
func DB() *sql.DB {
    db, err := sql.Open(
        "postgres",
        fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
            os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
            os.Getenv("DB_HOST"), os.Getenv("DB_TABLE")))

    if err != nil {
        log.Fatal(err)
    }

    return db
}

/**
 * Get s SQL Stub for variable binding
 */
func SqlStub(length int) string {
    var stub string

    for i := 1; i <= length; i++ {
        stub += fmt.Sprintf("$%s, ", strconv.Itoa(i))
    }

    return strings.TrimRight(stub, ", ")
}
