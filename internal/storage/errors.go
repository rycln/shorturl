package storage

import (
	"errors"
)

var (
	errConflict   = errors.New("shortened URL already exists")
	errNotExist   = errors.New("shortened URL does not exist")
	errDeletedURL = errors.New("URL removed")
)

// conflict represents an error when a resource already exists.
type conflict struct {
	err error
}

// Error returns the string representation of the error.
func (err *conflict) Error() string {
	return err.err.Error()
}

// Unwrap returns the underlying error.
func (err *conflict) Unwrap() error {
	return err.err
}

// IsErrConflict provides type checking capability.
// Always returns true for conflict errors.
func (err *conflict) IsErrConflict() bool {
	return true
}

// newErrConflict constructs a new conflict error.
func newErrConflict(err error) error {
	return &conflict{
		err: err,
	}
}

// notExist represents an error when a resource doesn't exist.
type notExist struct {
	err error
}

// Error returns the string representation of the error.
func (err *notExist) Error() string {
	return err.err.Error()
}

// Unwrap returns the underlying error.
func (err *notExist) Unwrap() error {
	return err.err
}

// IsErrNotExist provides type checking capability.
func (err *notExist) IsErrNotExist() bool {
	return true
}

// newErrNotExist constructs a new notExist error.
func newErrNotExist(err error) error {
	return &notExist{
		err: err,
	}
}

// deletedURL represents an error when accessing a soft-deleted URL.
type deletedURL struct {
	err error
}

// Error returns the string representation of the error.
func (err *deletedURL) Error() string {
	return err.err.Error()
}

// Unwrap returns the underlying error.
func (err *deletedURL) Unwrap() error {
	return err.err
}

// IsErrDeletedURL provides type checking capability.
func (err *deletedURL) IsErrDeletedURL() bool {
	return true
}

// newErrDeletedURL constructs a new deletedURL error.
func newErrDeletedURL(err error) error {
	return &deletedURL{
		err: err,
	}
}
