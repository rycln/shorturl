package contextkeys

type contextKey struct{}

// Package-level context keys for storing common request values.
var (
	// ShortURL is the context key for storing shortened URL value.
	// Used in middleware and handlers to pass parsed URL path.
	ShortURL = contextKey{}

	// UserID is the context key for storing authenticated user ID.
	// Populated by auth middleware after JWT verification.
	UserID = contextKey{}
)
