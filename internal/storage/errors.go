package storage

import (
	"errors"
)

var (
	ErrConflict   = errors.New("shortened URL already exists")
	ErrNotExist   = errors.New("shortened URL does not exist")
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

type errDeletedURL struct {
	err error
}

func (err *errDeletedURL) Error() string {
	return err.err.Error()
}

func (err *errDeletedURL) Unwrap() error {
	return err.err
}

func (err *errDeletedURL) IsErrDeletedURL() bool {
	return true
}

func newErrDeletedURL(err error) error {
	return &errDeletedURL{
		err: err,
	}
}
