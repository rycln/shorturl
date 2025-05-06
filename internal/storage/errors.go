package storage

import (
	"errors"
)

var (
	ErrConflict   = errors.New("shortened URL already exists")
	ErrNotExist   = errors.New("shortened URL does not exist")
	ErrTimeLimit  = errors.New("time limit exceeded")
	ErrDeletedURL = errors.New("URL removed")
)

type errConflict struct {
	err error
}

func (err *errConflict) Error() string {
	return err.err.Error()
}

func (err *errConflict) Unwrap() error {
	return err.err
}

func (err *errConflict) IsErrConflict() bool {
	return true
}

func newErrConflict(err error) error {
	return &errConflict{
		err: err,
	}
}

type errNotExist struct {
	err error
}

func (err *errNotExist) Error() string {
	return err.err.Error()
}

func (err *errNotExist) Unwrap() error {
	return err.err
}

func (err *errNotExist) IsErrNotExist() bool {
	return true
}

func newErrNotExist(err error) error {
	return &errNotExist{
		err: err,
	}
}
