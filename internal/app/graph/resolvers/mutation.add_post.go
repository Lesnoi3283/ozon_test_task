package resolvers

import (
	"context"
	"errors"
	"fmt"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"ozon_test_task/internal/app/graph/model"
	"ozon_test_task/internal/app/graph/repository"
	"ozon_test_task/internal/app/middlewares"
	"ozon_test_task/internal/app/models"
	"strconv"
)

// AddPost is the resolver for the addPost field.
func (r *mutationResolver) AddPost(ctx context.Context, title string, text string, commentsAllowed *bool) (*model.AddPostResponse, error) {
	user, ok := ctx.Value(middlewares.UserContextKey).(*models.User)
	if !ok {
		return nil, gqlerror.Errorf("Not authorized")
	}

	newPost := &models.Post{
		Owner:           *user,
		Title:           title,
		Text:            text,
		CommentsAllowed: *commentsAllowed,
	}

	postID, err := r.PostRepo.AddPost(ctx, newPost)
	if err != nil {
		if errors.Is(err, repository.NewErrConflict()) {
			return nil, fmt.Errorf("conflict")
		}
		return nil, fmt.Errorf("internal server error")
	}

	return &model.AddPostResponse{
		Post: &model.Post{
			ID:    strconv.Itoa(postID),
			Title: newPost.Title,
			Text:  newPost.Text,
			Owner: &model.User{
				ID:       strconv.Itoa(user.ID),
				Username: user.Login,
			},
			CommentsAllowed: newPost.CommentsAllowed,
		},
		Error: "",
	}, nil
}
