package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
)

// CommentReplies is the resolver for the commentReplies field.
func (r *queryResolver) CommentReplies(ctx context.Context, commentID string, first *int32, after *string) (*model.CommentConnection, error) {
	panic(fmt.Errorf("not implemented: CommentReplies - commentReplies"))
}
