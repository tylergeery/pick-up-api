package resources

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
	tx   PtxInterface
)

/**
 * Get a SQL DB connection
 */
func DB() *sql.DB {
	var err error

	// handles singleton like functionality
	// uses Mutex lock under the hood for atomic functionality
	once.Do(func() {
		db, err = sql.Open(
			"postgres",
			fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
				os.Getenv("DB_USER"), os.Getenv("DB_PASS"),
				os.Getenv("DB_HOST"), os.Getenv("DB_TABLE")))

		if err != nil {
			log.Fatal(err)
		}
	})

	return db
}

func TX() PtxInterface {
	if tx == nil {
		tx = &Ptx{tx: getTx()}
	}

	return tx
}

func getTx() *sql.Tx {
	tx, err := DB().Begin()

	if err != nil {
		log.Fatal(err)
	}

	return tx
}

func SetTestTXInterface() {
	tx = &PtxMock{tx: getTx()}
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
