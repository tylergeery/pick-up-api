package controllers

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/pick-up-api/models"
)

func UserProfile(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userId := vars["userId"]
    fmt.Fprintln(w, "Finding user with id:", userId)

    user := models.UserGetById(userId)
    fmt.Fprintln(w, user)
}
