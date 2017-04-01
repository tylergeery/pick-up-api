package controllers

import (
	"net/http"
	"strings"

	"github.com/pick-up-api/models/user"
	"github.com/pick-up-api/utils/auth"
	"github.com/pick-up-api/utils/messaging"
	"github.com/pick-up-api/utils/response"
)

/**
 * Refresh a user access token
 */
func RefreshToken(w http.ResponseWriter, r *http.Request) {
	bearer := r.Header.Get("Authorization")
	token := strings.TrimPrefix(bearer, "Bearer ")

	// validate refresh token
	userId, success := auth.GetUserIdFromRefreshToken(token)

	// TODO: check for blacklisted token

	if success {
		user, err := user.UserGetById(userId)

		if err == nil {
			// create a new token
			user.AddAccessToken()

			// emit success responser
			response.Success(w, user)
		} else {
			response.Fail(w, http.StatusServiceUnavailable, err.Error())
		}

	} else {
		// emit failure response
		response.Fail(w, http.StatusBadRequest, messaging.AUTH_INVALID_REFRESH_TOKEN)
	}
}
