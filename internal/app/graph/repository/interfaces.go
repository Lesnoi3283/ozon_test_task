package repository

import (
	"context"
	"ozon_test_task/internal/app/models"
)

type PostRepo interface {
	// AddPost adds a new post to a storage and returns it`s ID.
	AddPost(ctx context.Context, post *models.Post) (int, error)
	SetCommentsAllowed(ctx context.Context, postID int, commentsAllowed bool) error
	GetPostByID(ctx context.Context, postID int) (*models.Post, error)
	GetPosts(ctx context.Context, limit int, after int) ([]*models.Post, error)
}

type CommentRepo interface {
	// AddComment adds a new comment to a storage and returns it`s ID.
	AddComment(ctx context.Context, comment *models.Comment) (int, error)
	GetCommentsByPostID(ctx context.Context, postID int) ([]*models.Comment, error)
	GetReplaysByCommentID(ctx context.Context, commentID int) ([]*models.Comment, error)
}

type UserRepo interface {
	AddUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
}
