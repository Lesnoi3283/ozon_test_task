package resolvers

import (
	"context"
	"fmt"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	"ozon_test_task/cfg"
	"ozon_test_task/internal/app/graph/model"
	"ozon_test_task/internal/app/graph/repository"
	"ozon_test_task/internal/app/graph/repository/mocks"
	"ozon_test_task/internal/app/middlewares"
	"ozon_test_task/internal/app/models"
	"reflect"
	"testing"
)

func Test_mutationResolver_AddComment(t *testing.T) {
	type args struct {
		ctx    context.Context
		postID string
		text   string
	}
	type resolverFields struct {
		cfg            cfg.Cfg
		getPostRepo    func(c *gomock.Controller) repository.PostRepo
		getCommentRepo func(c *gomock.Controller) repository.CommentRepo
	}
	tests := []struct {
		name           string
		resolverFields resolverFields
		args           args
		want           *model.AddCommentResponse
		wantErr        bool
	}{
		{
			name: "Not authorized",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{},
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					return mocks.NewMockPostRepo(c)
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					return mocks.NewMockCommentRepo(c)
				},
			},
			args: args{
				ctx:    context.Background(),
				postID: "1",
				text:   "Hello",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "postID is not int",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{},
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					return mocks.NewMockPostRepo(c)
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					return mocks.NewMockCommentRepo(c)
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "qwerty"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				postID: "abc",
				text:   "Hello",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Comment text too long",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					MaxCommentTextLength: 10,
				},
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					return mocks.NewMockPostRepo(c)
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					return mocks.NewMockCommentRepo(c)
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "qwerty"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				postID: "1",
				text:   "01234567890",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Post not found",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					MaxCommentTextLength: 100,
				},
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPostByID(gomock.Any(), 1).Return(nil, fmt.Errorf("not found"))
					return pr
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					return mocks.NewMockCommentRepo(c)
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "qwerty"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				postID: "1",
				text:   "Hello",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Comment not allowed",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					MaxCommentTextLength: 100,
				},
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPostByID(gomock.Any(), 1).Return(&models.Post{
						ID:              1,
						CommentsAllowed: false,
					}, nil)
					return pr
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					return mocks.NewMockCommentRepo(c)
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "qwerty"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				postID: "1",
				text:   "Hello",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Failed to create comment",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					MaxCommentTextLength: 100,
				},
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPostByID(gomock.Any(), 1).Return(&models.Post{
						ID:              1,
						CommentsAllowed: true,
					}, nil)
					return pr
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					cr := mocks.NewMockCommentRepo(c)
					cr.EXPECT().AddComment(gomock.Any(), gomock.Any()).Return(0, fmt.Errorf("db error"))
					return cr
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "qwerty"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				postID: "1",
				text:   "Hello",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Ok",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					MaxCommentTextLength: 100,
				},
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPostByID(gomock.Any(), 1).Return(&models.Post{
						ID:              1,
						CommentsAllowed: true,
					}, nil)
					return pr
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					cr := mocks.NewMockCommentRepo(c)
					cr.EXPECT().AddComment(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, comment *models.Comment) (int, error) {
							return 123, nil
						},
					)
					return cr
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "qwerty"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				postID: "1",
				text:   "Hello",
			},
			want: &model.AddCommentResponse{
				Comment: &model.Comment{
					ID: "123",
					Owner: &model.User{
						ID:       "1",
						Username: "qwerty",
					},
					Text:      "Hello",
					CreatedAt: "",
				},
				Error: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := zaptest.NewLogger(t)
			sugar := logger.Sugar()
			c := gomock.NewController(t)
			r := &mutationResolver{
				Resolver: &Resolver{
					Logger:      sugar,
					Cfg:         tt.resolverFields.cfg,
					PostRepo:    tt.resolverFields.getPostRepo(c),
					CommentRepo: tt.resolverFields.getCommentRepo(c),
				},
			}
			got, err := r.AddComment(tt.args.ctx, tt.args.postID, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddComment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if got.Comment != nil {
					got.Comment.CreatedAt = ""
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddComment() got = %v, want %v", got, tt.want)
			}
		})
	}
}
