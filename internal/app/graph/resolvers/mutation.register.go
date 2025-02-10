package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
	"ozon_test_task/internal/app/models"
	"ozon_test_task/pkg/authUtils"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, username string, password string) (*model.AuthResponse, error) {
	//check data
	if len(username) == 0 {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if len(password) == 0 {
		return nil, fmt.Errorf("password cannot be empty")
	}

	//gen password salt and hash
	salt, err := authUtils.GenPasswordSalt()
	if err != nil {
		return nil, fmt.Errorf("internal server error")
	}

	passwordHash := authUtils.HashPassword(password, salt)

	//add user
	id, err := r.UserRepo.AddUser(ctx, &models.User{
		Login:        username,
		PasswordHash: passwordHash,
		PasswordSalt: salt,
	})
	if err != nil {
		return nil, fmt.Errorf("internal server error")
	}

	//make jwt
	jwt, err := r.JWTManager.BuildNewJWTString(id)
	if err != nil {
		return nil, fmt.Errorf("user created, auth error")
	}

	return &model.AuthResponse{
		Token: jwt,
		Error: "",
	}, nil
}
