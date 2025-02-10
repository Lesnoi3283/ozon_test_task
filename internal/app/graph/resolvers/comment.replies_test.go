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

func Test_commentResolver_Replies(t *testing.T) {

	//type fields struct {
	//	Resolver *Resolver
	//}
	type args struct {
		ctx   context.Context
		obj   *model.Comment
		limit *int32
		after *string
	}
	type resolverFields struct {
		cfg            cfg.Cfg
		getCommentRepo func(c *gomock.Controller) repository.CommentRepo
	}
	tests := []struct {
		name string
		//fields  fields
		resolverFields resolverFields
		args           args
		want           *model.CommentConnection
		wantErr        bool
	}{
		{
			name: "Cfg values",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					DefaultCommentsLimit: 20,
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					cr := mocks.NewMockCommentRepo(c)
					cr.EXPECT().GetReplaysByCommentID(gomock.Any(), 10, 20, 0).Return([]*models.Comment{
						{
							ID: 6,
							Owner: models.User{
								ID:    1,
								Login: "qwerty",
							},
							PostID:    5,
							ParentID:  10,
							Text:      "Hello",
							CreatedAt: time.Date(2020, 10, 30, 0, 0, 0, 0, time.UTC),
						},
						{
							ID: 7,
							Owner: models.User{
								ID:    2,
								Login: "ytrewq",
							},
							PostID:    5,
							ParentID:  10,
							Text:      "Hi",
							CreatedAt: time.Date(2020, 10, 30, 0, 0, 1, 0, time.UTC),
						},
					}, false, nil)
					return cr
				},
			},
			args: args{
				ctx: context.Background(),
				obj: &model.Comment{
					ID:        "10",
					Owner:     nil,
					Text:      "",
					CreatedAt: "",
					Replies:   nil,
				},
				limit: nil,
				after: nil,
			},
			want: &model.CommentConnection{
				Edges: []*model.CommentEdge{
					{
						Cursor: "6",
						Node: &model.Comment{
							ID: "6",
							Owner: &model.User{
								ID:       "1",
								Username: "qwerty",
							},
							Text:      "Hello",
							CreatedAt: time.Date(2020, 10, 30, 0, 0, 0, 0, time.UTC).String(),
							Replies:   nil,
						},
					},
					{
						Cursor: "7",
						Node: &model.Comment{
							ID: "7",
							Owner: &model.User{
								ID:       "2",
								Username: "ytrewq",
							},
							Text:      "Hi",
							CreatedAt: time.Date(2020, 10, 30, 0, 0, 1, 0, time.UTC).String(),
							Replies:   nil,
						},
					},
				},
				PageInfo: &model.PageInfo{
					StartCursor: func() *string {
						v := "6"
						return &v
					}(),
					EndCursor: func() *string {
						v := "7"
						return &v
					}(),
					HasNextPage: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Requested limit is bigger then cfg.max",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					MaxCommentsLimit: 2,
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					cr := mocks.NewMockCommentRepo(c)
					cr.EXPECT().GetReplaysByCommentID(gomock.Any(), 10, 2, 0).Return([]*models.Comment{
						{
							ID: 6,
							Owner: models.User{
								ID:    1,
								Login: "qwerty",
							},
							PostID:    5,
							ParentID:  10,
							Text:      "Hello",
							CreatedAt: time.Date(2020, 10, 30, 0, 0, 0, 0, time.UTC),
						},
						{
							ID: 7,
							Owner: models.User{
								ID:    2,
								Login: "ytrewq",
							},
							PostID:    5,
							ParentID:  10,
							Text:      "Hi",
							CreatedAt: time.Date(2020, 10, 30, 0, 0, 1, 0, time.UTC),
						},
					}, false, nil)
					return cr
				},
			},
			args: args{
				ctx: context.Background(),
				obj: &model.Comment{
					ID:        "10",
					Owner:     nil,
					Text:      "",
					CreatedAt: "",
					Replies:   nil,
				},
				limit: func() *int32 {
					v := int32(20)
					return &v
				}(),
				after: nil,
			},
			want: &model.CommentConnection{
				Edges: []*model.CommentEdge{
					{
						Cursor: "6",
						Node: &model.Comment{
							ID: "6",
							Owner: &model.User{
								ID:       "1",
								Username: "qwerty",
							},
							Text:      "Hello",
							CreatedAt: time.Date(2020, 10, 30, 0, 0, 0, 0, time.UTC).String(),
							Replies:   nil,
						},
					},
					{
						Cursor: "7",
						Node: &model.Comment{
							ID: "7",
							Owner: &model.User{
								ID:       "2",
								Username: "ytrewq",
							},
							Text:      "Hi",
							CreatedAt: time.Date(2020, 10, 30, 0, 0, 1, 0, time.UTC).String(),
							Replies:   nil,
						},
					},
				},
				PageInfo: &model.PageInfo{
					StartCursor: func() *string {
						v := "6"
						return &v
					}(),
					EndCursor: func() *string {
						v := "7"
						return &v
					}(),
					HasNextPage: false,
				},
			},
			wantErr: false,
		},
		{
			name: "DB err",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					MaxCommentsLimit: 200,
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					cr := mocks.NewMockCommentRepo(c)
					cr.EXPECT().GetReplaysByCommentID(gomock.Any(), 10, 20, 5).Return(nil, false, fmt.Errorf("Some test err"))
					return cr
				},
			},
			args: args{
				ctx: context.Background(),
				obj: &model.Comment{
					ID:        "10",
					Owner:     nil,
					Text:      "",
					CreatedAt: "",
					Replies:   nil,
				},
				limit: func() *int32 {
					v := int32(20)
					return &v
				}(),
				after: func() *string {
					v := "5"
					return &v
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Ok",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					MaxCommentsLimit: 200,
				},
				getCommentRepo: func(c *gomock.Controller) repository.CommentRepo {
					cr := mocks.NewMockCommentRepo(c)
					cr.EXPECT().GetReplaysByCommentID(gomock.Any(), 10, 20, 5).Return([]*models.Comment{
						{
							ID: 6,
							Owner: models.User{
								ID:    1,
								Login: "qwerty",
							},
							PostID:    5,
							ParentID:  10,
							Text:      "Hello",
							CreatedAt: time.Date(2020, 10, 30, 0, 0, 0, 0, time.UTC),
						},
						{
							ID: 7,
							Owner: models.User{
								ID:    2,
								Login: "ytrewq",
							},
							PostID:    5,
							ParentID:  10,
							Text:      "Hi",
							CreatedAt: time.Date(2020, 10, 30, 0, 0, 1, 0, time.UTC),
						},
					}, false, nil)
					return cr
				},
			},
			args: args{
				ctx: context.Background(),
				obj: &model.Comment{
					ID:        "10",
					Owner:     nil,
					Text:      "",
					CreatedAt: "",
					Replies:   nil,
				},
				limit: func() *int32 {
					v := int32(20)
					return &v
				}(),
				after: func() *string {
					v := "5"
					return &v
				}(),
			},
			want: &model.CommentConnection{
				Edges: []*model.CommentEdge{
					{
						Cursor: "6",
						Node: &model.Comment{
							ID: "6",
							Owner: &model.User{
								ID:       "1",
								Username: "qwerty",
							},
							Text:      "Hello",
							CreatedAt: time.Date(2020, 10, 30, 0, 0, 0, 0, time.UTC).String(),
							Replies:   nil,
						},
					},
					{
						Cursor: "7",
						Node: &model.Comment{
							ID: "7",
							Owner: &model.User{
								ID:       "2",
								Username: "ytrewq",
							},
							Text:      "Hi",
							CreatedAt: time.Date(2020, 10, 30, 0, 0, 1, 0, time.UTC).String(),
							Replies:   nil,
						},
					},
				},
				PageInfo: &model.PageInfo{
					StartCursor: func() *string {
						v := "6"
						return &v
					}(),
					EndCursor: func() *string {
						v := "7"
						return &v
					}(),
					HasNextPage: false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			r := &commentResolver{
				Resolver: &Resolver{
					Cfg:         tt.resolverFields.cfg,
					CommentRepo: tt.resolverFields.getCommentRepo(c),
				},
			}
			got, err := r.Replies(tt.args.ctx, tt.args.obj, tt.args.limit, tt.args.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("Replies() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Replies() got = %v, want %v", got, tt.want)
			}
		})
	}
}
