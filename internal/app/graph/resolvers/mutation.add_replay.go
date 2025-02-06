package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
)

// AddReplay is the resolver for the addReplay field.
func (r *mutationResolver) AddReplay(ctx context.Context, parentCommentID string, text string) (*model.AddReplayResponse, error) {
	panic(fmt.Errorf("not implemented: AddReplay - addReplay"))
}
