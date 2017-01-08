package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/pick-up-api/utils/auth"
)

func SetUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		token := strings.TrimPrefix(bearer, "Bearer ")
		userId, success := auth.GetUserIdFromToken(token)

		if success {
			ctx := context.WithValue(r.Context(), auth.USER_ID_KEY, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
