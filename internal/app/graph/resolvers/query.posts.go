package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
)

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, first *int32, after *string) (*model.PostConnection, error) {
	panic(fmt.Errorf("not implemented: Posts - posts"))
}
