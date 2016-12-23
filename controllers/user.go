package controllers

import (
    "errors"
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

    if err == nil {
        response.Success(w, user)
    } else {
        response.Fail(w, http.StatusNotFound, err.Error())
    }
}

/**
 * Create a user
 */
func UserCreate(w http.ResponseWriter, r *http.Request) {
    var err error
    var user models.User

    r.ParseForm()

    _, emailExists := r.Form["email"];
    _, pwExists := r.Form["password"];

    if emailExists && pwExists {
        user, err = models.UserCreateProfile(r.Form)
        user.AddToken()
    } else {
        if !emailExists {
            err = errors.New("Users require an email")
        } else {
            err = errors.New("Users require a password")
        }
    }

    if err == nil {
        response.Success(w, user)
    } else {
        response.Fail(w, http.StatusBadRequest, err.Error())
    }
}

/**
 * Update a user
 */
 func UserUpdate(w http.ResponseWriter, r *http.Request) {
     var err error
     var user models.User

     r.ParseForm()

     if userId, exists := r.Form["userId"]; exists {
         user, err = models.UserUpdateProfile(userId[0], r.Form)
     } else {
         err = errors.New("User ID not specified")
     }

     if err == nil {
         response.Success(w, user)
     } else {
         response.Fail(w, http.StatusBadRequest, err.Error())
     }
 }

/**
 * Delete a user
 */
func UserDelete(w http.ResponseWriter, r *http.Request) {
    var err error

    r.ParseForm()

    if userId, exists := r.Form["userId"]; exists {
        err = models.UserDeleteProfile(userId[0])
    } else {
        err = errors.New("User ID not specified")
    }

    if err == nil {
        response.Success(w, struct{Message string}{"User successfully removed"})
    } else {
        response.Fail(w, http.StatusBadRequest, err.Error())
    }
}
