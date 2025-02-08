package resolvers

import (
	"context"
	"errors"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
	"ozon_test_task/internal/app/graph/repository"
	"ozon_test_task/pkg/authUtils"
)

// Auth is the resolver for the auth field.
// Returns "user not found" even if user was found but password is incorrect - due to secure reasons.
func (r *mutationResolver) Auth(ctx context.Context, username string, password string) (*model.AuthResponse, error) {
	//pre-check data
	if len(username) == 0 {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if len(password) == 0 {
		return nil, fmt.Errorf("password cannot be empty")
	}

	//auth
	user, err := r.UserRepo.GetUserByLoginWithCred(ctx, username)
	if err != nil {
		if errors.Is(err, repository.NewErrNotFound()) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user by login: %w", err)
	}

	if authUtils.CheckPassword(password, user.PasswordHash, user.PasswordSalt) {
		jwt, err := r.JWTManager.BuildNewJWTString(user.ID)
		if err != nil {
			return nil, fmt.Errorf("internal server error")
		}
		return &model.AuthResponse{
			Token: jwt,
			Error: "",
		}, nil
	}

	return nil, fmt.Errorf("user not found")
}
