package resolvers

import (
	"context"
	"fmt"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	"ozon_test_task/internal/app/graph/model"
	"ozon_test_task/internal/app/graph/repository"
	"ozon_test_task/internal/app/graph/repository/mocks"
	"ozon_test_task/internal/app/middlewares"
	"ozon_test_task/internal/app/models"
	"reflect"
	"testing"
)

func Test_mutationResolver_SetCommentsAllowed(t *testing.T) {
	type args struct {
		ctx     context.Context
		postID  string
		allowed bool
	}
	type resolverFields struct {
		getPostRepo func(c *gomock.Controller) repository.PostRepo
	}
	tests := []struct {
		name           string
		resolverFields resolverFields
		args           args
		want           *model.Post
		wantErr        bool
	}{
		{
			name: "Not authorized",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					return mocks.NewMockPostRepo(c)
				},
			},
			args: args{
				ctx:     context.Background(),
				postID:  "10",
				allowed: true,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "postID is not int",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					return mocks.NewMockPostRepo(c)
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "user1"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				postID:  "abc",
				allowed: true,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "post not found",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPostByID(gomock.Any(), 10).Return(nil, repository.NewErrNotFound())
					return pr
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "user1"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				postID:  "10",
				allowed: true,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "db err get post",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPostByID(gomock.Any(), 10).Return(nil, fmt.Errorf("db error"))
					return pr
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "user1"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				postID:  "10",
				allowed: true,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "not an owner",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPostByID(gomock.Any(), 10).Return(&models.Post{
						ID:    10,
						Title: "title",
						Text:  "text",
						Owner: models.User{ID: 2, Login: "another_user"},
					}, nil)
					return pr
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "user1"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				postID:  "10",
				allowed: true,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "db err set allowed",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPostByID(gomock.Any(), 10).Return(&models.Post{
						ID:    10,
						Title: "title",
						Text:  "text",
						Owner: models.User{ID: 1, Login: "user1"},
					}, nil)
					pr.EXPECT().SetCommentsAllowed(gomock.Any(), 10, false).Return(fmt.Errorf("db error"))
					return pr
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "user1"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				postID:  "10",
				allowed: false,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ok",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPostByID(gomock.Any(), 10).Return(&models.Post{
						ID:    10,
						Title: "title",
						Text:  "text",
						Owner: models.User{ID: 1, Login: "user1"},
					}, nil)
					pr.EXPECT().SetCommentsAllowed(gomock.Any(), 10, true).Return(nil)
					return pr
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "user1"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				postID:  "10",
				allowed: true,
			},
			want: &model.Post{
				ID:              "10",
				Title:           "title",
				Text:            "text",
				Owner:           &model.User{ID: "1", Username: "user1"},
				CommentsAllowed: true,
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
					Logger:   sugar,
					PostRepo: tt.resolverFields.getPostRepo(c),
				},
			}
			got, err := r.SetCommentsAllowed(tt.args.ctx, tt.args.postID, tt.args.allowed)
			if (err != nil) != tt.wantErr {
				t.Errorf("SetCommentsAllowed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetCommentsAllowed() got = %v, want %v", got, tt.want)
			}
		})
	}
}
