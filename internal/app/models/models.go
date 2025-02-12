package models

import "time"

type Post struct {
	ID              int
	Owner           User
	Title           string
	Text            string
	CommentsAllowed bool
}

type Comment struct {
	ID        int
	Owner     User
	PostID    int //zero if comment is a sub-comment
	ParentID  int //zero if comment doesnt have parent.
	Text      string
	CreatedAt time.Time
}

type User struct {
	ID           int
	Login        string
	PasswordHash string
	PasswordSalt string
}
