package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
)

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	panic(fmt.Errorf("not implemented: Post - post"))
}
