package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/pick-up-api/controllers"
)

func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc("/hello/{name}", Hello).Methods("GET")

    // User API
    router.HandleFunc("/user/create", controllers.UserCreate).Methods("POST")
    router.HandleFunc("/user/update", controllers.UserUpdate).Methods("POST")
    router.HandleFunc("/user/delete", controllers.UserDelete).Methods("POST")
    router.HandleFunc("/user/{userId}", controllers.UserProfile).Methods("GET")

    // Event API

    log.Fatal(http.ListenAndServe(":3001", router))
}

func Hello(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fmt.Fprintf(w, "Hello, %s!", vars["name"])
}