package database

import (
	"context"
	"errors"
	"fmt"
	"ozon_test_task/internal/app/graph/repository"
	"ozon_test_task/internal/app/models"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

// RepoRedis is a Redis repository that implements PostRepo, CommentRepo, and UserRepo.
type RepoRedis struct {
	client *redis.Client
}

// NewRepoRedis returns a new RepoRedis.
func NewRepoRedis(client *redis.Client) *RepoRedis {
	return &RepoRedis{client: client}
}

// AddPost adds a new post and returns its ID.
func (r *RepoRedis) AddPost(ctx context.Context, post *models.Post) (int, error) {
	// inc counter.
	id64, err := r.client.Incr(ctx, "counter:post").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to generate post id: %w", err)
	}
	postID := int(id64)
	key := fmt.Sprintf("post:%d", postID)

	//h set
	err = r.client.HSet(ctx, key, map[string]interface{}{
		"owner_id":        post.Owner.ID,
		"title":           post.Title,
		"text":            post.Text,
		"commentsallowed": post.CommentsAllowed,
	}).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to add post: %w", err)
	}

	//add
	if err := r.client.ZAdd(ctx, "posts", &redis.Z{Score: float64(postID), Member: postID}).Err(); err != nil {
		return 0, fmt.Errorf("failed to add post to sorted set: %w", err)
	}
	return postID, nil
}

// SetCommentsAllowed updates the "commentsallowed" field for a given post.
func (r *RepoRedis) SetCommentsAllowed(ctx context.Context, postID int, commentsAllowed bool) error {
	key := fmt.Sprintf("post:%d", postID)
	if err := r.client.HSet(ctx, key, "commentsallowed", commentsAllowed).Err(); err != nil {
		return fmt.Errorf("failed to update comments allowed: %w", err)
	}
	return nil
}

// GetPostByID returns a post by its ID.
func (r *RepoRedis) GetPostByID(ctx context.Context, postID int) (*models.Post, error) {
	//get post data
	key := fmt.Sprintf("post:%d", postID)
	m, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get post: %w", err)
	}
	if len(m) == 0 {
		return nil, fmt.Errorf("post not found")
	}
	ownerID, err := strconv.Atoi(m["owner_id"])
	if err != nil {
		return nil, fmt.Errorf("invalid owner_id: %w", err)
	}
	commentsAllowed, err := strconv.ParseBool(m["commentsallowed"])
	if err != nil {
		return nil, fmt.Errorf("invalid commentsallowed: %w", err)
	}
	post := &models.Post{
		ID:              postID,
		Owner:           models.User{ID: ownerID},
		Title:           m["title"],
		Text:            m["text"],
		CommentsAllowed: commentsAllowed,
	}

	//get owner data
	owner, err := r.GetUserByID(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner: %w", err)
	}
	post.Owner = *owner
	return post, nil
}

// GetPosts returns a list of posts.
func (r *RepoRedis) GetPosts(ctx context.Context, limit int, after int) (posts []*models.Post, hasNextPage bool, err error) {
	min := fmt.Sprintf("(%d", after)
	max := "+inf"
	zRangeArgs := &redis.ZRangeBy{
		Min:   min,
		Max:   max,
		Count: int64(limit + 1),
	}
	idStrs, err := r.client.ZRangeByScore(ctx, "posts", zRangeArgs).Result()
	if err != nil {
		return nil, false, fmt.Errorf("failed to get posts from sorted set: %w", err)
	}

	if len(idStrs) > limit {
		hasNextPage = true
		idStrs = idStrs[:limit]
	}
	for _, idStr := range idStrs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, false, fmt.Errorf("invalid post id: %w", err)
		}
		post, err := r.GetPostByID(ctx, id)
		if err != nil {
			return nil, false, fmt.Errorf("failed to get post by id: %w", err)
		}
		posts = append(posts, post)
	}
	return posts, hasNextPage, nil
}

// AddComment adds a new comment to Redis and returns its ID.
func (r *RepoRedis) AddComment(ctx context.Context, comment *models.Comment) (int, error) {
	id64, err := r.client.Incr(ctx, "counter:comment").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to generate comment id: %w", err)
	}
	commentID := int(id64)
	key := fmt.Sprintf("comment:%d", commentID)

	err = r.client.HSet(ctx, key, map[string]interface{}{
		"owner_id":   comment.Owner.ID,
		"post_id":    comment.PostID,
		"parent_id":  comment.ParentID,
		"text":       comment.Text,
		"created_at": comment.CreatedAt.Unix(),
	}).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to add comment: %w", err)
	}

	if comment.ParentID == 0 {
		// Top-level comment for a post
		setKey := fmt.Sprintf("post:%d:comments", comment.PostID)
		if err := r.client.ZAdd(ctx, setKey, &redis.Z{Score: float64(commentID), Member: commentID}).Err(); err != nil {
			return 0, fmt.Errorf("failed to add comment to post sorted set: %w", err)
		}
	} else {
		// Reply
		setKey := fmt.Sprintf("comment:%d:replies", comment.ParentID)
		if err := r.client.ZAdd(ctx, setKey, &redis.Z{Score: float64(commentID), Member: commentID}).Err(); err != nil {
			return 0, fmt.Errorf("failed to add reply to comment sorted set: %w", err)
		}
	}
	return commentID, nil
}

// GetCommentsByPostID retrieves top-level comments (without replays) for a post.
func (r *RepoRedis) GetCommentsByPostID(ctx context.Context, postID int, limit int, after int) (comments []*models.Comment, hasNextPage bool, err error) {
	setKey := fmt.Sprintf("post:%d:comments", postID)
	min := fmt.Sprintf("(%d", after)
	max := "+inf"
	zRangeArgs := &redis.ZRangeBy{
		Min:   min,
		Max:   max,
		Count: int64(limit + 1),
	}
	idStrs, err := r.client.ZRangeByScore(ctx, setKey, zRangeArgs).Result()
	if err != nil {
		return nil, false, fmt.Errorf("failed to get comment ids: %w", err)
	}
	if len(idStrs) > limit {
		hasNextPage = true
		idStrs = idStrs[:limit]
	}
	for _, idStr := range idStrs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, false, fmt.Errorf("invalid comment id: %w", err)
		}
		comment, err := r.getCommentByID(ctx, id)
		if err != nil {
			return nil, false, fmt.Errorf("failed to get comment by id: %w", err)
		}
		comments = append(comments, comment)
	}
	return comments, hasNextPage, nil
}

// GetReplaysByCommentID returns replies for a given comment.
func (r *RepoRedis) GetReplaysByCommentID(ctx context.Context, commentID int, limit int, after int) (replies []*models.Comment, hasNextPage bool, err error) {
	setKey := fmt.Sprintf("comment:%d:replies", commentID)
	min := fmt.Sprintf("(%d", after)
	max := "+inf"
	zRangeArgs := &redis.ZRangeBy{
		Min:   min,
		Max:   max,
		Count: int64(limit + 1),
	}
	idStrs, err := r.client.ZRangeByScore(ctx, setKey, zRangeArgs).Result()
	if err != nil {
		return nil, false, fmt.Errorf("failed to get reply ids: %w", err)
	}
	if len(idStrs) > limit {
		hasNextPage = true
		idStrs = idStrs[:limit]
	}
	for _, idStr := range idStrs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return nil, false, fmt.Errorf("invalid reply id: %w", err)
		}
		reply, err := r.getCommentByID(ctx, id)
		if err != nil {
			return nil, false, fmt.Errorf("failed to get reply by id: %w", err)
		}
		replies = append(replies, reply)
	}
	return replies, hasNextPage, nil
}

// getCommentByID returns a comment by its ID.
func (r *RepoRedis) getCommentByID(ctx context.Context, commentID int) (*models.Comment, error) {
	//get comment data
	key := fmt.Sprintf("comment:%d", commentID)
	m, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get comment: %w", err)
	}
	if len(m) == 0 {
		return nil, fmt.Errorf("comment not found")
	}
	ownerID, err := strconv.Atoi(m["owner_id"])
	if err != nil {
		return nil, fmt.Errorf("invalid owner_id: %w", err)
	}
	postID, err := strconv.Atoi(m["post_id"])
	if err != nil {
		return nil, fmt.Errorf("invalid post_id: %w", err)
	}
	parentID, err := strconv.Atoi(m["parent_id"])
	if err != nil {
		return nil, fmt.Errorf("invalid parent_id: %w", err)
	}
	createdAtUnix, err := strconv.ParseInt(m["created_at"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid created_at: %w", err)
	}
	comment := &models.Comment{
		ID:        commentID,
		Owner:     models.User{ID: ownerID},
		PostID:    postID,
		ParentID:  parentID,
		Text:      m["text"],
		CreatedAt: time.Unix(createdAtUnix, 0),
	}

	//get owner data
	owner, err := r.GetUserByID(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get comment owner: %w", err)
	}
	comment.Owner = *owner
	return comment, nil
}

// AddUser adds a new user to Redis and returns it`s ID.
func (r *RepoRedis) AddUser(ctx context.Context, user *models.User) (int, error) {
	id64, err := r.client.Incr(ctx, "counter:user").Result()
	if err != nil {
		return 0, fmt.Errorf("failed to generate user id: %w", err)
	}
	userID := int(id64)

	//save user
	userKey := fmt.Sprintf("user:%d", userID)
	err = r.client.HSet(ctx, userKey, map[string]interface{}{
		"login":        user.Login,
		"passwordhash": user.PasswordHash,
		"passwordsalt": user.PasswordSalt,
	}).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to add user: %w", err)
	}

	//save login-userID
	loginKey := fmt.Sprintf("login:%s", user.Login)
	err = r.client.Set(ctx, loginKey, userID, 0).Err()
	if err != nil {
		return 0, fmt.Errorf("failed to set login mapping: %w", err)
	}

	return userID, nil
}

// GetUserByID returns a user by its ID without password hash and salt.
func (r *RepoRedis) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	key := fmt.Sprintf("user:%d", userID)
	m, err := r.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if len(m) == 0 {
		return nil, repository.NewErrNotFound()
	}
	user := &models.User{
		ID:    userID,
		Login: m["login"],
	}
	return user, nil
}

// GetUserByLoginWithCred returns a user by its login with credentials (password_hash and password_salt).
func (r *RepoRedis) GetUserByLoginWithCred(ctx context.Context, login string) (*models.User, error) {
	// get userID by login
	loginKey := fmt.Sprintf("login:%s", login)
	userID, err := r.client.Get(ctx, loginKey).Int()
	if errors.Is(err, redis.Nil) {
		return nil, repository.NewErrNotFound()
	} else if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// get user by ID
	userKey := fmt.Sprintf("user:%d", userID)
	m, err := r.client.HGetAll(ctx, userKey).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if len(m) == 0 {
		return nil, repository.NewErrNotFound()
	}

	//return answer
	user := &models.User{
		ID:           userID,
		Login:        m["login"],
		PasswordHash: m["passwordhash"],
		PasswordSalt: m["passwordsalt"],
	}

	return user, nil
}
