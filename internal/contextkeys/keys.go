package contextkeys

type contextKey struct{}

var (
	ShortURL = contextKey{}
	UserID   = contextKey{}
)
