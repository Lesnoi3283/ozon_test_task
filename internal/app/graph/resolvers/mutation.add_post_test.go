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

func Test_mutationResolver_AddPost(t *testing.T) {
	type args struct {
		ctx             context.Context
		title           string
		text            string
		commentsAllowed *bool
	}
	type resolverFields struct {
		getPostRepo func(c *gomock.Controller) repository.PostRepo
	}
	tests := []struct {
		name           string
		resolverFields resolverFields
		args           args
		want           *model.AddPostResponse
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
				ctx:             context.Background(),
				title:           "Title",
				text:            "Text",
				commentsAllowed: func() *bool { v := true; return &v }(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Conflict error",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().AddPost(gomock.Any(), gomock.Any()).Return(0, repository.NewErrConflict())
					return pr
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "qwerty"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				title:           "Title",
				text:            "Text",
				commentsAllowed: func() *bool { v := true; return &v }(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Internal server error",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().AddPost(gomock.Any(), gomock.Any()).Return(0, fmt.Errorf("some error"))
					return pr
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "qwerty"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				title:           "Title",
				text:            "Text",
				commentsAllowed: func() *bool { v := true; return &v }(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Ok",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().AddPost(gomock.Any(), gomock.Any()).DoAndReturn(
						func(ctx context.Context, post *models.Post) (int, error) {
							return 42, nil
						},
					)
					return pr
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "qwerty"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				title:           "Title",
				text:            "Text",
				commentsAllowed: func() *bool { v := true; return &v }(),
			},
			want: &model.AddPostResponse{
				Post: &model.Post{
					ID:              "42",
					Title:           "Title",
					Text:            "Text",
					Owner:           &model.User{ID: "1", Username: "qwerty"},
					CommentsAllowed: true,
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
					Logger:   sugar,
					PostRepo: tt.resolverFields.getPostRepo(c),
				},
			}
			got, err := r.AddPost(tt.args.ctx, tt.args.title, tt.args.text, tt.args.commentsAllowed)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddPost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddPost() got = %v, want %v", got, tt.want)
			}
		})
	}
}
