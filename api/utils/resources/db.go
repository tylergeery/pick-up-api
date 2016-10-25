package resources

import (
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
            "raccoon", "pickEmUp", "pickup-postgres", "pickup"))

    if err != nil {
        log.Fatal(err)
    }

    return db
}

func SqlStub(length int) string {
    var stub string

    for i := 1; i <= length; i++ {
        stub += fmt.Sprintf("$%s, ", strconv.Itoa(i))
    }

    return strings.TrimRight(stub, ", ")
}
