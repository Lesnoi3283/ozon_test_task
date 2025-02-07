package resolvers

import (
	"ozon_test_task/cfg"
	"ozon_test_task/internal/app/graph/repository"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserRepo    repository.UserRepo
	PostRepo    repository.PostRepo
	CommentRepo repository.CommentRepo
	Cfg         cfg.Cfg
}
