package resolvers

import (
	"context"
	"fmt"
	"go.uber.org/mock/gomock"
	"ozon_test_task/cfg"
	"ozon_test_task/internal/app/graph/model"
	"ozon_test_task/internal/app/graph/repository"
	"ozon_test_task/internal/app/graph/repository/mocks"
	"ozon_test_task/internal/app/models"
	"reflect"
	"testing"
	"time"
)

func Test_postResolver_Comments(t *testing.T) {
	type args struct {
		ctx   context.Context
		obj   *model.Post
		limit *int32
		after *string
	}
	type resolverFields struct {
		cfg            cfg.Cfg
		getCommentRepo func(c *gomock.Controller) repository.CommentRepo
	}
	tests := []struct {
		name           string
		resolverFields resolverFields
		args           args
		want           *model.CommentConnection
		wantErr        bool
	}{
		{
			name: "postID is not int",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					return mocks.NewMockCommentRepo(c)
				},
			},
			args: args{
				ctx: context.Background(),
				obj: &model.Post{ID: "abc"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "after is not convertable to int",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					return mocks.NewMockCommentRepo(c)
				},
			},
			args: args{
				ctx: context.Background(),
				obj: &model.Post{ID: "10"},
				after: func() *string {
					v := "abc"
					return &v
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "DB err",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					DefaultCommentsLimit: 10,
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					cr := mocks.NewMockCommentRepo(c)
					cr.EXPECT().GetCommentsByPostID(gomock.Any(), 10, 10, 0).Return(nil, false, fmt.Errorf("db error"))
					return cr
				},
			},
			args: args{
				ctx:   context.Background(),
				obj:   &model.Post{ID: "10"},
				limit: nil,
				after: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "cfg values",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					DefaultCommentsLimit: 10,
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					cr := mocks.NewMockCommentRepo(c)
					cr.EXPECT().GetCommentsByPostID(gomock.Any(), 10, 10, 0).Return([]*models.Comment{}, false, nil)
					return cr
				},
			},
			args: args{
				ctx:   context.Background(),
				obj:   &model.Post{ID: "10"},
				limit: nil,
				after: nil,
			},
			want: &model.CommentConnection{
				Edges: []*model.CommentEdge{},
				PageInfo: &model.PageInfo{
					StartCursor: func() *string { s := ""; return &s }(),
					EndCursor:   func() *string { s := ""; return &s }(),
					HasNextPage: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Ok",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					DefaultCommentsLimit: 10,
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					cr := mocks.NewMockCommentRepo(c)
					cr.EXPECT().GetCommentsByPostID(gomock.Any(), 10, 10, 0).Return([]*models.Comment{
						{
							ID:        11,
							Text:      "comment1",
							CreatedAt: time.Date(2020, 10, 31, 0, 0, 0, 0, time.UTC),
							Owner: models.User{
								ID:    1,
								Login: "user1",
							},
						},
						{
							ID:        12,
							Text:      "comment2",
							CreatedAt: time.Date(2020, 10, 31, 0, 0, 1, 0, time.UTC),
							Owner: models.User{
								ID:    2,
								Login: "user2",
							},
						},
					}, true, nil)
					return cr
				},
			},
			args: args{
				ctx:   context.Background(),
				obj:   &model.Post{ID: "10"},
				limit: nil,
				after: nil,
			},
			want: &model.CommentConnection{
				Edges: []*model.CommentEdge{
					{
						Cursor: "11",
						Node: &model.Comment{
							ID:        "11",
							Text:      "comment1",
							CreatedAt: time.Date(2020, 10, 31, 0, 0, 0, 0, time.UTC).String(),
							Owner: &model.User{
								ID:       "1",
								Username: "user1",
							},
						},
					},
					{
						Cursor: "12",
						Node: &model.Comment{
							ID:        "12",
							Text:      "comment2",
							CreatedAt: time.Date(2020, 10, 31, 0, 0, 1, 0, time.UTC).String(),
							Owner: &model.User{
								ID:       "2",
								Username: "user2",
							},
						},
					},
				},
				PageInfo: &model.PageInfo{
					StartCursor: func() *string { s := "11"; return &s }(),
					EndCursor:   func() *string { s := "12"; return &s }(),
					HasNextPage: true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			p := &postResolver{
				Resolver: &Resolver{
					Cfg:         tt.resolverFields.cfg,
					CommentRepo: tt.resolverFields.getCommentRepo(c),
				},
			}
			got, err := p.Comments(tt.args.ctx, tt.args.obj, tt.args.limit, tt.args.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("Comments() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Comments() got = %v, want %v", got, tt.want)
			}
		})
	}
}
