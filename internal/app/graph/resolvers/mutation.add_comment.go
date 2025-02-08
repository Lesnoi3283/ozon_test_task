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

// AddComment is the resolver for the addComment field.
func (r *mutationResolver) AddComment(ctx context.Context, postID string, text string) (*model.AddCommentResponse, error) {
	user, ok := ctx.Value(middlewares.UserContextKey).(*models.User)
	if !ok {
		return nil, gqlerror.Errorf("Not authorized")
	}

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		return nil, fmt.Errorf("postID is not int")
	}

	if len(text) > r.Cfg.MaxCommentTextLength {
		return nil, fmt.Errorf("comment text too long, max lenght: %d", r.Cfg.MaxCommentTextLength)
	}

	comment := &models.Comment{
		Owner:     *user,
		Text:      text,
		CreatedAt: time.Now(),
		PostID:    postIDInt,
	}

	//check if comments are allowed
	post, err := r.PostRepo.GetPostByID(ctx, postIDInt)
	if err != nil {
		return nil, fmt.Errorf("post not found")
	}
	if !post.CommentsAllowed {
		return nil, gqlerror.Errorf("Comment is not allowed to this post")
	}

	commentID, err := r.CommentRepo.AddComment(ctx, comment)
	if err != nil {
		return nil, fmt.Errorf("failed to create a comment")
	}

	commentModel := &model.Comment{
		ID: strconv.Itoa(commentID),
		Owner: &model.User{
			ID:       strconv.Itoa(comment.Owner.ID),
			Username: comment.Owner.Login,
		},
		Text:      comment.Text,
		CreatedAt: comment.CreatedAt.String(),
	}

	return &model.AddCommentResponse{
		Comment: commentModel,
		Error:   "",
	}, nil
}
