package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/mux"

    "database/sql"

	_ "github.com/lib/pq"
)

func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/hello/{name}", Hello)
    router.HandleFunc("/user/{userId}", UserProfile)

    log.Fatal(http.ListenAndServe(":3001", router))
}

func Hello(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fmt.Fprintf(w, "Hello, %s!", vars["name"])
}

func UserProfile(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userId := vars["userId"]
    fmt.Fprintln(w, "Finding user with id:", userId)

    db, err := sql.Open("postgres", "postgres://raccoon:password@192.168.99.100/pqgotest?sslmode=disable")
	if err != nil {
        fmt.Fprintln(w, "Found error: ", err)
		log.Fatal(err)
	}

	rows, err := db.Query("SELECT * FROM users WHERE id = $id", userId)
    fmt.Fprintln(w, "Found rows: ", rows)
}
