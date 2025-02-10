package middlewares

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"ozon_test_task/internal/app/graph/repository"
)

//go:generate mockgen -source auth.go -destination mocks/mock_mw.go -package mwmocks

type UserContextKeyType string

const UserContextKey UserContextKeyType = "userContextKey"

const AuthHeaderName = "Authorization"

type JWTManager interface {
	BuildNewJWTString(userID int) (string, error)
	GetUserID(token string) (int, error)
}

// GetAuthMiddleware - returns authentication middleware func.
// It will just try to get user using JWT, but it won`t deny access for unauthorised users.
// After successful authentication it will add "*models.User" to a context using UserContextKey as a key.
// If user is unauthorised - nothing would be in context.
func GetAuthMiddleware(manager JWTManager, userRepo repository.UserRepo, logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(AuthHeaderName)
			if authHeader == "" {
				logger.Debugf("Auth header is empty")
				next.ServeHTTP(w, r)
				return
			}

			userID, err := manager.GetUserID(authHeader)
			if err != nil {
				logger.Debugf("cant get userID from jwt, err: %v", err)
				next.ServeHTTP(w, r)
				return
			}
			user, err := userRepo.GetUserByID(r.Context(), userID)
			if err != nil {
				logger.Debugf("cant get user from db, err: %v", err)
				next.ServeHTTP(w, r)
				return
			}

			logger.Debugf("success auth!")
			ctx := context.WithValue(r.Context(), UserContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
