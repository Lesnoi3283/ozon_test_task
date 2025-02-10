package resolvers

import (
	"go.uber.org/zap"
	"ozon_test_task/cfg"
	"ozon_test_task/internal/app/graph/model"
	"ozon_test_task/internal/app/graph/repository"
	"ozon_test_task/internal/app/middlewares"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserRepo     repository.UserRepo
	PostRepo     repository.PostRepo
	CommentRepo  repository.CommentRepo
	Cfg          cfg.Cfg
	JWTManager   middlewares.JWTManager
	CommentAdded chan *model.Comment
	Logger       *zap.SugaredLogger
}
