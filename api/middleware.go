package api

import (
	"context"
	"log"
	"net/http"
	"strings"

	"ants/db"
)

// Define our struct
type authenticationMiddleware struct {
}

var ctx = context.Background()

// Middleware function, which will be called for each request
func (amw *authenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Split the Authorization header to extract the token
		splitToken := strings.Split(authHeader, " ")
		if len(splitToken) != 2 || splitToken[0] != "Bearer" {
			// The Authorization header is not in the expected format
			http.Error(w, "Invalid Authorization Token", http.StatusUnauthorized)
			return
		}

		// Extracted token without "Bearer"
		token := splitToken[1]

		if user, err := db.Redis.Get(ctx, "token:"+token+":user").Result(); err == nil {
			log.Printf("Authenticated user %s\n", user)
			// Pass down the request to the next middleware (or final handler)
			next.ServeHTTP(w, r)
		} else {
			// Token not found in Redis or other error occurred
			log.Printf("Error: %s\n", err)
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}
