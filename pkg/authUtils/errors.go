package authUtils

import "errors"

var errJWTIsNotValid = errors.New("JWT token is not valid")

func NewErrJWTIsNotValid() error {
	return errJWTIsNotValid
}
