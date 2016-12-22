package middleware

import (
    //"errors"
    "log"
    "net/http"
    //"github.com/gorilla/mux"
    //"github.com/pick-up-api/utils/auth"
)

func IsAuthorized(h http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        bearer := r.Header.Get("Authorization")
        log.Println(bearer)
        h.ServeHTTP(w, r)
    })
}
