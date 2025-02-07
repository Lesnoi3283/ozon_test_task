package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
	"ozon_test_task/internal/app/models"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, username string, password string) (*model.AuthResponse, error) {
	if len(username) == 0 {
		return nil, fmt.Errorf("username cannot be empty")
	}
	if len(password) == 0 {
		return nil, fmt.Errorf("password cannot be empty")
	}

	r.UserRepo.AddUser(ctx, &models.User{
		Login:        username,
		PasswordHash: auth,
	})
}
