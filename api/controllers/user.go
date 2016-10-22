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
        message := fmt.Sprintf("User with id %s could not be found", userId)
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

    if err == nil {
        response.Success(w, user)
    } else {
        response.Fail(w, http.StatusBadRequest, "User could not be created")
    }
}

/**
 * Update a user
 */
 func UserUpdate(w http.ResponseWriter, r *http.Request) {
     r.ParseForm()

     user, err := models.UserUpdateProfile(r.Form)

     if err == nil {
         response.Success(w, user)
     } else {
         response.Fail(w, http.StatusBadRequest, "User could not be created")
     }
 }

/**
 * Delete a user
 */
func UserDelete(w http.ResponseWriter, r *http.Request) {
    var err error

    r.ParseForm()

    if val, exists := r.Form["userId"]; exists {
        err = models.UserDeleteProfile(val[0])
    }

    if err == nil {
        response.Success(w, struct{message string}{"User successfully removed"})
    } else {
        response.Fail(w, http.StatusBadRequest, "User could not be be deleted")
    }
}
