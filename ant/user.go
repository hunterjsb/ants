package ant

import "context"

type User struct {
	Name     string    `json:"name"`
	Colonies []*Colony `json:"colonies"`
}

// Define a key type that's unexported to other packages.
type keyType int

// Define constants for your keys.
const (
	userKey keyType = iota // Unique key for user-related data.
)

// NewContextWithUser adds a user ID to the context.
func (u User) AddToContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, userKey, u)
}

// UserFromContext retrieves a user ID from the context.
func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userKey).(User)
	return user, ok
}
