package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"ozon_test_task/internal/app/graph/repository"
	"ozon_test_task/internal/app/models"
)

// RepoPG is a PostgreSQL repository that implements PostRepo, CommentRepo, and UserRepo interfaces.
type RepoPG struct {
	DB *sql.DB
}

// NewRepoPG returns a new RepoPG.
func NewRepoPG(db *sql.DB) *RepoPG {
	return &RepoPG{DB: db}
}

// InitDB creates the necessary tables in the PostgreSQL database if they do not already exist.
func InitDB(db *sql.DB) error {
	ctx := context.Background()

	usersTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login VARCHAR(255) NOT NULL UNIQUE,
		password_hash VARCHAR(255) NOT NULL,
	    password_salt VARCHAR(255) NOT NULL,
	);`
	if _, err := db.ExecContext(ctx, usersTableQuery); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	postsTableQuery := `
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		owner_id INTEGER NOT NULL,
		title VARCHAR(255) NOT NULL,
		text TEXT NOT NULL,
		commentsallowed BOOLEAN NOT NULL DEFAULT TRUE,
		FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
	);`
	if _, err := db.ExecContext(ctx, postsTableQuery); err != nil {
		return fmt.Errorf("failed to create posts table: %w", err)
	}

	commentsTableQuery := `
	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL PRIMARY KEY,
		owner_id INTEGER NOT NULL,
		post_id INTEGER NOT NULL,
		parent_id INTEGER NOT NULL DEFAULT 0,
		text TEXT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE,
		FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE
	);`
	if _, err := db.ExecContext(ctx, commentsTableQuery); err != nil {
		return fmt.Errorf("failed to create comments table: %w", err)
	}

	return nil
}

// AddPost adds a new post to the database and returns its generated ID.
func (r *RepoPG) AddPost(ctx context.Context, post *models.Post) (int, error) {
	var id int
	query := `
		INSERT INTO posts (owner_id, title, text, commentsallowed)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	err := r.DB.QueryRowContext(ctx, query, post.Owner.ID, post.Title, post.Text, post.CommentsAllowed).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add post: %w", err)
	}
	return id, nil
}

// SetCommentsAllowed updates the commentsallowed flag for a given post.
func (r *RepoPG) SetCommentsAllowed(ctx context.Context, postID int, commentsAllowed bool) error {
	query := `UPDATE posts SET commentsallowed = $1 WHERE id = $2`
	result, err := r.DB.ExecContext(ctx, query, commentsAllowed, postID)
	if err != nil {
		return fmt.Errorf("failed to update comments allowed: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetPostByID returns a post by its ID.
func (r *RepoPG) GetPostByID(ctx context.Context, postID int) (*models.Post, error) {
	query := `
		SELECT p.id, p.title, p.text, p.commentsallowed,
		       u.id, u.login
		FROM posts p
		JOIN users u ON p.owner_id = u.id
		WHERE p.id = $1`
	row := r.DB.QueryRowContext(ctx, query, postID)

	var p models.Post
	var u models.User
	if err := row.Scan(&p.ID, &p.Title, &p.Text, &p.CommentsAllowed, &u.ID, &u.Login); err != nil {
		return nil, fmt.Errorf("failed to get post by ID: %w", err)
	}
	p.Owner = u
	return &p, nil
}

// GetPosts retrieves a list of posts with pagination.
func (r *RepoPG) GetPosts(ctx context.Context, limit int, after int) (posts []*models.Post, hasNextPage bool, err error) {
	limitPlusOne := limit + 1
	query := `
		SELECT p.id, p.title, p.text, p.commentsallowed,
		       u.id, u.login
		FROM posts p
		JOIN users u ON p.owner_id = u.id
		WHERE p.id > $1
		ORDER BY p.id
		LIMIT $2`
	rows, err := r.DB.QueryContext(ctx, query, after, limitPlusOne)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get posts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var p models.Post
		var u models.User
		if err := rows.Scan(&p.ID, &p.Title, &p.Text, &p.CommentsAllowed, &u.ID, &u.Login); err != nil {
			return nil, false, fmt.Errorf("failed to scan post: %w", err)
		}
		p.Owner = u
		posts = append(posts, &p)
	}
	if err := rows.Err(); err != nil {
		return nil, false, fmt.Errorf("rows error: %w", err)
	}
	if len(posts) > limit {
		hasNextPage = true
		posts = posts[:limit]
	}
	return posts, hasNextPage, nil
}

// AddComment adds a new comment to the database and returns its ID.
func (r *RepoPG) AddComment(ctx context.Context, comment *models.Comment) (int, error) {
	var id int
	query := `
		INSERT INTO comments (owner_id, post_id, parent_id, text, created_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	err := r.DB.QueryRowContext(ctx, query, comment.Owner.ID, comment.PostID, comment.ParentID, comment.Text, comment.CreatedAt).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to add comment: %w", err)
	}
	return id, nil
}

// GetCommentsByPostID returns top-level comments (without a parent or their sub-comments) for a given post.
func (r *RepoPG) GetCommentsByPostID(ctx context.Context, postID int, limit int, after int) (comments []*models.Comment, hasNextPage bool, err error) {
	limitPlusOne := limit + 1
	query := `
		SELECT c.id, c.post_id, c.parent_id, c.text, c.created_at,
		       u.id, u.login
		FROM comments c
		JOIN users u ON c.owner_id = u.id
		WHERE c.post_id = $1 AND c.parent_id = 0 AND c.id > $2
		ORDER BY c.id
		LIMIT $3`
	rows, err := r.DB.QueryContext(ctx, query, postID, after, limitPlusOne)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get comments by post ID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Comment
		var u models.User
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParentID, &c.Text, &c.CreatedAt, &u.ID, &u.Login); err != nil {
			return nil, false, fmt.Errorf("failed to scan comment: %w", err)
		}
		c.Owner = u
		comments = append(comments, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, false, fmt.Errorf("rows error: %w", err)
	}
	if len(comments) > limit {
		hasNextPage = true
		comments = comments[:limit]
	}
	return comments, hasNextPage, nil
}

// GetReplaysByCommentID gets replies for a given comment.
func (r *RepoPG) GetReplaysByCommentID(ctx context.Context, commentID int, limit int, after int) (replies []*models.Comment, hasNextPage bool, err error) {
	limitPlusOne := limit + 1
	query := `
		SELECT c.id, c.post_id, c.parent_id, c.text, c.created_at,
		       u.id, u.login
		FROM comments c
		JOIN users u ON c.owner_id = u.id
		WHERE c.parent_id = $1 AND c.id > $2
		ORDER BY c.id
		LIMIT $3`
	rows, err := r.DB.QueryContext(ctx, query, commentID, after, limitPlusOne)
	if err != nil {
		return nil, false, fmt.Errorf("failed to get replies by comment ID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var c models.Comment
		var u models.User
		if err := rows.Scan(&c.ID, &c.PostID, &c.ParentID, &c.Text, &c.CreatedAt, &u.ID, &u.Login); err != nil {
			return nil, false, fmt.Errorf("failed to scan reply: %w", err)
		}
		c.Owner = u
		replies = append(replies, &c)
	}
	if err := rows.Err(); err != nil {
		return nil, false, fmt.Errorf("rows error: %w", err)
	}
	if len(replies) > limit {
		hasNextPage = true
		replies = replies[:limit]
	}
	return replies, hasNextPage, nil
}

// AddUser adds a new user to the database and returns the user's ID.
func (r *RepoPG) AddUser(ctx context.Context, user *models.User) (int, error) {
	userID := 0
	query := `
		INSERT INTO users (login, password_hash, password_salt)
		VALUES ($1, $2, $3)
		RETURNING id`
	err := r.DB.QueryRowContext(ctx, query, user.Login, user.PasswordHash, user.PasswordSalt).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("failed to add user: %w", err)
	}
	return userID, nil
}

// GetUserByID returns a user from the database by its ID without password hash and salt.
func (r *RepoPG) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	query := `SELECT id, login FROM users WHERE id = $1`
	row := r.DB.QueryRowContext(ctx, query, userID)

	var u models.User
	if err := row.Scan(&u.ID, &u.Login); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NewErrNotFound()
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &u, nil
}

// GetUserByIDWithCred returns a user by its login with credentials (password_hash and password_salt).
func (r *RepoPG) GetUserByLoginWithCred(ctx context.Context, login string) (*models.User, error) {
	query := `SELECT id, login, password_hash, password_salt FROM users WHERE login = $1`
	row := r.DB.QueryRowContext(ctx, query, login)

	var u models.User
	if err := row.Scan(&u.ID, &u.Login, &u.PasswordHash, &u.PasswordSalt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, repository.NewErrNotFound()
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &u, nil
}
