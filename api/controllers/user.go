package controllers

import (
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
    r.ParseForm()

    user, err := models.UserCreateProfile(r.Form)

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
     r.ParseForm()

     user, err := models.UserUpdateProfile(r.Form)

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

    if val, exists := r.Form["userId"]; exists {
        err = models.UserDeleteProfile(val[0])
    }

    if err == nil {
        response.Success(w, struct{Message string}{"User successfully removed"})
    } else {
        response.Fail(w, http.StatusBadRequest, err.Error())
    }
}
