package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
)

// AddComment is the resolver for the addComment field.
func (r *mutationResolver) AddComment(ctx context.Context, postID string, text string) (*model.AddCommentResponse, error) {
	panic(fmt.Errorf("not implemented: AddComment - addComment"))
}
