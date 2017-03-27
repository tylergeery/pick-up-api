package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pick-up-api/utils/auth"
)

/**
 * Refresh a user access token
 */
func refreshToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bearer := r.Header.Get("Authorization")
	token := strings.TrimPrefix(bearer, "Bearer ")

	// validate refresh token
	userId, success := auth.GetUserIdFromRefreshToken(token)

	if success {
		// create a new token
		// save to db for user
		// emit success responser
		log.Println(userId, vars)
	} else {
		// emit failure response
	}
}
