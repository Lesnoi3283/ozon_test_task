package models

type Post struct {
	ID              int
	OwnerID         int
	Title           string
	Text            string
	CommentsAllowed bool
}

type Comment struct {
	ID       int
	PostID   int
	ParentID int //zero if comment doesnt have parent.
	Text     string
}

type User struct {
	ID           int
	Login        string
	PasswordHash string
}
