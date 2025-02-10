package resolvers

import (
	"context"
	"fmt"
	"go.uber.org/mock/gomock"
	"ozon_test_task/cfg"
	"ozon_test_task/internal/app/graph/model"
	"ozon_test_task/internal/app/graph/repository"
	"ozon_test_task/internal/app/graph/repository/mocks"
	"ozon_test_task/internal/app/middlewares"
	"ozon_test_task/internal/app/models"
	"reflect"
	"testing"
)

func Test_mutationResolver_AddReplay(t *testing.T) {
	type args struct {
		ctx             context.Context
		parentCommentID string
		text            string
	}
	type resolverFields struct {
		cfg            cfg.Cfg
		getCommentRepo func(c *gomock.Controller) repository.CommentRepo
	}
	tests := []struct {
		name           string
		resolverFields resolverFields
		args           args
		want           *model.AddReplayResponse
		wantErr        bool
	}{
		{
			name: "Not authorized",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					return mocks.NewMockCommentRepo(c)
				},
			},
			args: args{
				ctx:             context.Background(),
				parentCommentID: "10",
				text:            "Hello",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "parentCommentID is not int",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					return mocks.NewMockCommentRepo(c)
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "qwerty"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				parentCommentID: "abc",
				text:            "Hello",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Replay text too long",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					MaxCommentTextLength: 5,
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
				parentCommentID: "10",
				text:            "long_text_here",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Internal server error",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					MaxCommentTextLength: 100,
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					cr := mocks.NewMockCommentRepo(c)
					cr.EXPECT().AddComment(gomock.Any(), gomock.Any()).Return(0, fmt.Errorf("some db error"))
					return cr
				},
			},
			args: args{
				ctx: func() context.Context {
					user := &models.User{ID: 1, Login: "qwerty"}
					return context.WithValue(context.Background(), middlewares.UserContextKey, user)
				}(),
				parentCommentID: "10",
				text:            "Hello",
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
				parentCommentID: "10",
				text:            "Hello",
			},
			want: &model.AddReplayResponse{
				Comment: &model.Comment{
					ID: "123",
					Owner: &model.User{
						ID:       "1",
						Username: "qwerty",
					},
					Text:      "Hello",
					CreatedAt: "",
					Replies:   nil,
				},
				Error: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			r := &mutationResolver{
				Resolver: &Resolver{
					Cfg:         tt.resolverFields.cfg,
					CommentRepo: tt.resolverFields.getCommentRepo(c),
				},
			}
			got, err := r.AddReplay(tt.args.ctx, tt.args.parentCommentID, tt.args.text)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddReplay() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != nil && got.Comment != nil {
				// Ignore CreatedAt in tests
				got.Comment.CreatedAt = ""
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddReplay() got = %v, want %v", got, tt.want)
			}
		})
	}
}
