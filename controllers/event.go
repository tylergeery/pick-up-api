package controllers

import (
	"net/http"
	"strconv"

	"github.com/pick-up-api/models/event"
	"github.com/pick-up-api/utils/response"
)

/**
 * Get an event as JSON
 */
func EventCreate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventId, _ := strconv.ParseInt(vars["eventId"], 10, 64)

	event, err := event.EventGetById(userId)

	if err == nil {
		response.Success(w, event)
	} else {
		response.Fail(w, http.StatusNotFound, err.Error())
	}
}
