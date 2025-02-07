package middlewares

import "errors"

var errJWTIsNotValid = errors.New("JWT is not valid")

func NewErrJWTIsNotValid() error {
	return errJWTIsNotValid
}
