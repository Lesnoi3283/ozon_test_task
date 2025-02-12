package authUtils

import (
	"github.com/golang-jwt/jwt/v4"
	"ozon_test_task/internal/app/middlewares"
	"time"
)

const secretKey = "SuperSecretKeyPart"

type JWTHelper struct {
}

func NewJWTHelper() *JWTHelper {
	return &JWTHelper{}
}

type claims struct {
	UserID int
	jwt.RegisteredClaims
}

func (j *JWTHelper) BuildNewJWTString(userID int) (string, error) {
	claims := claims{
		UserID:           userID,
		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72))},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	stringToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return stringToken, nil
}

func (j *JWTHelper) GetUserID(token string) (int, error) {
	claims := claims{}

	tokenGot, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return -1, err
	}
	if !tokenGot.Valid {
		return -1, middlewares.NewErrJWTIsNotValid()
	}

	return claims.UserID, nil
}
