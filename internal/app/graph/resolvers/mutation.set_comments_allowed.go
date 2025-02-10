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

// SetCommentsAllowed is the resolver for the setCommentsAllowed field.
func (r *mutationResolver) SetCommentsAllowed(ctx context.Context, postID string, allowed bool) (*model.Post, error) {
	user, ok := ctx.Value(middlewares.UserContextKey).(*models.User)
	if !ok {
		r.Logger.Debugf("cant get user from ctx")
		return nil, gqlerror.Errorf("Not authorized")
	}

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		r.Logger.Debugf("cant convert postID to int, err: %v", err)
		return nil, fmt.Errorf("post id is not int")
	}

	//check if user is owner of this post
	post, err := r.PostRepo.GetPostByID(ctx, postIDInt)
	if err != nil {
		r.Logger.Debugf("cant get post from db, err: %v", err)
		if errors.Is(err, repository.NewErrNotFound()) {
			return nil, gqlerror.Errorf("post not found")
		}
		return nil, fmt.Errorf("internal server error")
	}

	if user.ID != post.Owner.ID {
		r.Logger.Debugf("cant modify this post, user is not an owner. UserID is \"%v\", but ownerID is \"%v\"", user.ID, post.Owner.ID)
		return nil, gqlerror.Errorf("cant modify this post")
	}

	//set allowed
	err = r.PostRepo.SetCommentsAllowed(ctx, postIDInt, allowed)
	if err != nil {
		r.Logger.Debugf("cant set comments allowed, err: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	//return response
	return &model.Post{
		ID:    strconv.Itoa(post.ID),
		Title: post.Title,
		Text:  post.Text,
		Owner: &model.User{
			ID:       strconv.Itoa(user.ID),
			Username: user.Login,
		},
		CommentsAllowed: allowed,
	}, nil
}
