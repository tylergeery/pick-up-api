package router

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pick-up-api/controllers"
	"github.com/pick-up-api/middleware"
)

func GetRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/hello/{name}", hello).Methods("GET")

	// User API
	router.HandleFunc("/user/create", controllers.UserCreate).Methods("POST")
	router.Handle("/user/update", auth(controllers.UserUpdate)).Methods("POST")
	router.Handle("/user/delete", auth(controllers.UserDelete)).Methods("POST")
	router.HandleFunc("/user/{userId}", controllers.UserProfile).Methods("GET")

	// Events API

	return router
}

func auth(cb func(w http.ResponseWriter, r *http.Request)) http.Handler {
	handler := http.HandlerFunc(cb)

	return middleware.SetUser(handler)
}

func hello(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Fprintf(w, "Hello, %s!", vars["name"])
}
