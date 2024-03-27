package midllewares

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var BearerToken = os.Getenv("TOKEN")

func Authenticate(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized: Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokens := strings.Split(authHeader, " ")
		if len(tokens) != 2 || tokens[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenValue := tokens[1]

		if tokenValue != BearerToken {
			fmt.Println(tokenValue, BearerToken)
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		h.ServeHTTP(w, r)
	})
}
