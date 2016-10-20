package controllers

import (
    "fmt"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/pick-up-api/models"
    "github.com/pick-up-api/utils/response"
)

/**
 * Get a user profile as JSON
 */
func UserProfile(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userId := vars["userId"]

    user, err := models.UserGetById(userId)

    if err != nil {
        message := fmt.Sprintf("User with id %d could not be found", userId)
        response.Fail(w, http.StatusNotFound, message)
    } else {
        response.Success(w, user)
    }
}

/**
 * Create a user
 */
func UserCreate(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()

    user, err := models.UserCreateProfile(r.Form)

    if err != nil {
        response.Fail(w, http.StatusBadRequest, "User could not be created")
    } else {
        response.Success(w, user)
    }
}
