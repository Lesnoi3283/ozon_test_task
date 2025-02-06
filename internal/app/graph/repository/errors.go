package repository

import "errors"

var errConflict = errors.New("conflict")

func NewErrConflict() error {
	return errConflict
}
