package controllers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pick-up-api/models"
	"github.com/pick-up-api/utils/auth"
	"github.com/pick-up-api/utils/response"
)

/**
 * Get a user profile as JSON
 */
func UserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId, _ := strconv.ParseInt(vars["userId"], 10, 64)

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

	_, emailExists := r.Form["email"]
	_, pwExists := r.Form["password"]

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
	requestId := r.Context().Value(auth.USER_ID_KEY).(int64)

	r.ParseForm()

	if userIdArray, exists := r.Form["userId"]; exists {
		userId, _ := strconv.ParseInt(userIdArray[0], 10, 64)

		if userId == requestId {
			user, err = models.UserUpdateProfile(userId, r.Form)
		} else {
			log.Println("Auth ID: ", requestId)
			err = errors.New("You are not authorized to update this user")
		}
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
		response.Success(w, struct{ Message string }{"User successfully removed"})
	} else {
		response.Fail(w, http.StatusBadRequest, err.Error())
	}
}
