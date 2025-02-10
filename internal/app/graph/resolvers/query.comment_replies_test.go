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
	"ozon_test_task/internal/app/models"
	"reflect"
	"testing"
	"time"
)

func Test_queryResolver_CommentReplies(t *testing.T) {
	type args struct {
		ctx       context.Context
		commentID string
		limit     *int32
		after     *string
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
			name: "commentID is not int",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					return mocks.NewMockCommentRepo(c)
				},
			},
			args: args{
				ctx:       context.Background(),
				commentID: "abc",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "after is not int",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					return mocks.NewMockCommentRepo(c)
				},
			},
			args: args{
				ctx:       context.Background(),
				commentID: "10",
				limit:     nil,
				after:     func() *string { v := "abc"; return &v }(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "DB error",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					DefaultCommentsLimit: 10,
					MaxCommentsLimit:     20,
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					cr := mocks.NewMockCommentRepo(c)
					cr.EXPECT().GetReplaysByCommentID(gomock.Any(), 10, 10, 0).Return(nil, false, fmt.Errorf("db error"))
					return cr
				},
			},
			args: args{
				ctx:       context.Background(),
				commentID: "10",
				limit:     nil,
				after:     nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Cfg value",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					DefaultCommentsLimit: 5,
					MaxCommentsLimit:     10,
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					cr := mocks.NewMockCommentRepo(c)
					cr.EXPECT().GetReplaysByCommentID(gomock.Any(), 10, 5, 0).Return([]*models.Comment{}, false, nil)
					return cr
				},
			},
			args: args{
				ctx:       context.Background(),
				commentID: "10",
				limit:     nil,
				after:     nil,
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
					MaxCommentsLimit:     10,
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					cr := mocks.NewMockCommentRepo(c)
					cr.EXPECT().GetReplaysByCommentID(gomock.Any(), 10, 10, 0).Return([]*models.Comment{
						{
							ID: 1,
							Owner: models.User{
								ID:    10,
								Login: "user10",
							},
							Text:      "reply1",
							CreatedAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC),
						},
						{
							ID: 2,
							Owner: models.User{
								ID:    11,
								Login: "user11",
							},
							Text:      "reply2",
							CreatedAt: time.Date(2020, 12, 1, 0, 0, 1, 0, time.UTC),
						},
					}, true, nil)
					return cr
				},
			},
			args: args{
				ctx:       context.Background(),
				commentID: "10",
				limit:     nil,
				after:     nil,
			},
			want: &model.CommentConnection{
				Edges: []*model.CommentEdge{
					{
						Cursor: "1",
						Node: &model.Comment{
							ID:        "1",
							Text:      "reply1",
							CreatedAt: time.Date(2020, 12, 1, 0, 0, 0, 0, time.UTC).String(),
							Owner: &model.User{
								ID:       "10",
								Username: "user10",
							},
							Replies: nil,
						},
					},
					{
						Cursor: "2",
						Node: &model.Comment{
							ID:        "2",
							Text:      "reply2",
							CreatedAt: time.Date(2020, 12, 1, 0, 0, 1, 0, time.UTC).String(),
							Owner: &model.User{
								ID:       "11",
								Username: "user11",
							},
							Replies: nil,
						},
					},
				},
				PageInfo: &model.PageInfo{
					StartCursor: func() *string { s := "1"; return &s }(),
					EndCursor:   func() *string { s := "2"; return &s }(),
					HasNextPage: true,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			logger := zaptest.NewLogger(t)
			sugar := logger.Sugar()
			r := &queryResolver{
				Resolver: &Resolver{
					Logger:      sugar,
					Cfg:         tt.resolverFields.cfg,
					CommentRepo: tt.resolverFields.getCommentRepo(c),
				},
			}
			got, err := r.CommentReplies(tt.args.ctx, tt.args.commentID, tt.args.limit, tt.args.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommentReplies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommentReplies() got = %v, want %v", got, tt.want)
			}
		})
	}
}
