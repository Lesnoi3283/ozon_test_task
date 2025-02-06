package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
)

// SetCommentsAllowed is the resolver for the setCommentsAllowed field.
func (r *mutationResolver) SetCommentsAllowed(ctx context.Context, postID string, allowed bool) (*model.Post, error) {
	panic(fmt.Errorf("not implemented: SetCommentsAllowed - setCommentsAllowed"))
}
