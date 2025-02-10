package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
	"strconv"
)

// Post is the resolver for the post field.
func (r *queryResolver) Post(ctx context.Context, id string) (*model.Post, error) {
	IDInt, err := strconv.Atoi(id)
	if err != nil {
		r.Logger.Debugf("cant convert ID to int, err: %v", err)
		return nil, fmt.Errorf("ID is not convertable to int")
	}

	post, err := r.PostRepo.GetPostByID(ctx, IDInt)
	if err != nil {
		r.Logger.Debugf("cant get post from db, err: %v", err)
		return nil, err
	}

	return &model.Post{
		ID:    strconv.Itoa(post.ID),
		Title: post.Title,
		Text:  post.Text,
		Owner: &model.User{
			ID:       strconv.Itoa(post.Owner.ID),
			Username: post.Owner.Login,
		},
		CommentsAllowed: post.CommentsAllowed,
	}, nil
}
