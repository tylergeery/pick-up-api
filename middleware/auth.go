package middleware

import (
    "log"
    "strconv"
    "strings"
    "context"
    "net/http"
    "github.com/pick-up-api/utils/auth"
)

func SetUser(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        bearer := r.Header.Get("Authorization")
        token := strings.TrimPrefix(bearer, "Bearer ")
        userId, success := auth.GetUserIdFromToken(token)

        log.Println("Auth Token: " + token)
        log.Println("User id: " + strconv.FormatInt(userId, 10))

        if (success) {
            ctx := context.WithValue(r.Context(), auth.USER_ID_KEY, userId)
            next.ServeHTTP(w, r.WithContext(ctx))
        } else {
            next.ServeHTTP(w, r)
        }
    })
}
