package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"ants/ant"
	"ants/db"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		log.Println(r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

type AuthenticationMiddleware struct{}

// Middleware enhances HTTP requests with authentication logic.
func (amw *AuthenticationMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		// Ensure the Authorization header is well-formed.
		token, err := extractToken(authHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// Verify the token.
		ctx := r.Context()
		// Assuming verifyToken returns the username of the authenticated user:
		userName, err := verifyToken(ctx, token)
		if err != nil {
			http.Error(w, "Forbidden: invalid token", http.StatusForbidden)
			return
		}

		// Create a User object (assuming you have a struct like this)
		user := ant.User{Name: userName} // TODO load colonies

		// Attach the User to the context
		ctxWithUser := user.AddToContext(ctx)

		// Create a new request with the updated context
		rWithUser := r.WithContext(ctxWithUser)

		// Continue to the next handler, using the new request
		next.ServeHTTP(w, rWithUser)
	})
}

// extractToken parses the Authorization header to extract the bearer token.
func extractToken(authHeader string) (string, error) {
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 || splitToken[0] != "Bearer" {
		return "", fmt.Errorf("invalid Authorization token")
	}
	return splitToken[1], nil
}

// verifyToken checks the validity of the token and retrieves the associated user.
func verifyToken(ctx context.Context, token string) (string, error) {
	user, err := db.Redis.Get(ctx, "token:"+token+":user").Result()
	if err != nil {
		return "", err // Consider wrapping this error to distinguish between Redis errors and not found errors.
	}
	return user, nil
}
