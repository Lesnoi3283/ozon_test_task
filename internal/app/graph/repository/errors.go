package repository

import "errors"

var errConflict = errors.New("conflict")

func NewErrConflict() error {
	return errConflict
}

var errNotFound = errors.New("not found")

func NewErrNotFound() error {
	return errNotFound
}
