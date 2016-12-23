package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/pick-up-api/controllers"
    "github.com/pick-up-api/middleware"
)

func main() {

    router := mux.NewRouter().StrictSlash(true)
    router.Handle("/hello/{name}", Hello).Methods("GET")

    // User API
    router.HandleFunc("/user/create", controllers.UserCreate).Methods("POST")
    router.HandleFunc("/user/update", auth(controllers.UserUpdate)).Methods("POST")
    router.HandleFunc("/user/delete", controllers.UserDelete).Methods("POST")
    router.HandleFunc("/user/{userId}", controllers.UserProfile).Methods("GET")

    // Event API

    // auth middleware
    log.Fatal(http.ListenAndServe(":3001", router))
}

func Hello(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    fmt.Fprintf(w, "Hello, %s!", vars["name"])
}

func auth(cb func(w http.ResponseWriter, r *http.Request)) http.Handler {
    handler := http.HandlerFunc(cb)

    return middleware.setUser(handler)
}
