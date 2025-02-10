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
		r.Logger.Debugf("username is empty")
		return nil, fmt.Errorf("username cannot be empty")
	}
	if len(password) == 0 {
		r.Logger.Debugf("password is empty")
		return nil, fmt.Errorf("password cannot be empty")
	}

	//gen password salt and hash
	salt, err := authUtils.GenPasswordSalt()
	if err != nil {
		r.Logger.Errorf("failed to generate salt: %v", err)
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
		r.Logger.Errorf("failed to add user to a db: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	//make jwt
	jwt, err := r.JWTManager.BuildNewJWTString(id)
	if err != nil {
		r.Logger.Errorf("failed to build jwt string: %v", err)
		return nil, fmt.Errorf("user created, auth error")
	}

	return &model.AuthResponse{
		Token: jwt,
		Error: "",
	}, nil
}
