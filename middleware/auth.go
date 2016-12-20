package middleware

import (
    //"errors"
    //"log"
    "net/http"
    "github.com/pick-up-api/utils/response"
    //"github.com/gorilla/mux"
    //"github.com/pick-up-api/utils/auth"
)

func IsAuthorized(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //TODO
        // Read authorization header
        response.Fail(w, 401, "Forbidden");
        return
    })
}
