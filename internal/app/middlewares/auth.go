package middlewares

import "net/http"

type UserContextKeyType string

var UserContextKey UserContextKeyType = "userContextKey"

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//todo...

	})
}
