package middleware

import (
    "log"
    "strings"
    "context"
    "net/http"
    "github.com/pick-up-api/utils/auth"
    "github.com/pick-up-api/utils/response"
)

func setUser(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        bearer := r.Header.Get("Authorization")
        token := strings.TrimPrefix(bearer, "Bearer ")
        userId, success := auth.GetUserIdFromToken(token)

        log.Println("Auth Token: " + token)
        log.Println("User id: " + userId)

        if (success) {
            ctx := context.WithValue(r.Context(), auth.USER_ID_KEY, userId)
            next.ServeHTTP(w, r.WithContext(ctx))
        } else {
            next.ServeHTTP(w, r)
        }
    })
}
