package resolvers

import (
	"context"
	"fmt"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"ozon_test_task/internal/app/graph/model"
	"ozon_test_task/internal/app/middlewares"
	"ozon_test_task/internal/app/models"
	"strconv"
	"time"
)

// AddReplay is the resolver for the addReplay field.
func (r *mutationResolver) AddReplay(ctx context.Context, parentCommentID string, text string) (*model.AddReplayResponse, error) {
	user, ok := ctx.Value(middlewares.UserContextKey).(*models.User)
	if !ok {
		r.Logger.Debugf("cant get user from ctx")
		return nil, gqlerror.Errorf("Not authorized")
	}

	parentIDInt, err := strconv.Atoi(parentCommentID)
	if err != nil {
		r.Logger.Debugf("parentCommentID is not an int, parsing err: %v", err)
		return nil, fmt.Errorf("parentCommentID is not int")
	}

	if len(text) > r.Cfg.MaxCommentTextLength {
		r.Logger.Debugf("max comment length exceeded, current len is \"%v\", max len is \"%v\"", len(text), r.Cfg.MaxCommentTextLength)
		return nil, fmt.Errorf("replay text too long, max lenght: %d", r.Cfg.MaxCommentTextLength)
	}

	comment := &models.Comment{
		Owner:     *user,
		PostID:    0, // zero means comment is a sub-comment.
		ParentID:  parentIDInt,
		Text:      text,
		CreatedAt: time.Now(),
	}

	id, err := r.CommentRepo.AddComment(ctx, comment)
	if err != nil {
		r.Logger.Debugf("cant add comment to a db, err: %v", err)
		return nil, fmt.Errorf("internal server error")
	}

	return &model.AddReplayResponse{
		Comment: &model.Comment{
			ID: strconv.Itoa(id),
			Owner: &model.User{
				ID:       strconv.Itoa(comment.Owner.ID),
				Username: comment.Owner.Login,
			},
			Text:      comment.Text,
			CreatedAt: comment.CreatedAt.String(),
			Replies:   nil,
		},
		Error: "",
	}, nil
}
