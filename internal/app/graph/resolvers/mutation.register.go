package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
)

// Register is the resolver for the register field.
func (r *mutationResolver) Register(ctx context.Context, username string, password string) (*model.AuthResponse, error) {
	panic(fmt.Errorf("not implemented: Register - register"))
}
