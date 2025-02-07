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
	// GetPosts returns "limit" amount of posts or less, after "after" post`s id (comment with "after" id won`t be selected).
	// Also returns hasNextPage true if it`s exists more comments in database after last selected one.
	GetPosts(ctx context.Context, limit int, after int) (posts []*models.Post, hasNextPage bool, err error)
}

type CommentRepo interface {
	// AddComment adds a new comment to a storage and returns it`s ID.
	AddComment(ctx context.Context, comment *models.Comment) (int, error)
	// GetCommentsByPostID returns "limit" amount of comments or less, after "after" comment`s id (comment with "after" id won`t be selected).
	// Also returns hasNextPage true if it`s exists more comments in database after last selected one.
	GetCommentsByPostID(ctx context.Context, postID int, limit int, after int) (comments []*models.Comment, hasNextPage bool, err error)
	// GetReplaysByCommentID returns "limit" amount of comments (replays) or less, after "after" comment`s id (comment with "after" id won`t be selected).
	// Also returns hasNextPage true if it`s exists more comments in database after last selected one.
	GetReplaysByCommentID(ctx context.Context, commentID int, limit int, after int) (replays []*models.Comment, hasNextPage bool, err error)
}

type UserRepo interface {
	AddUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
}
