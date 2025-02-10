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
)

func Test_queryResolver_Posts(t *testing.T) {
	type args struct {
		ctx   context.Context
		limit *int32
		after *string
	}
	type resolverFields struct {
		cfg         cfg.Cfg
		getPostRepo func(c *gomock.Controller) repository.PostRepo
	}
	tests := []struct {
		name           string
		resolverFields resolverFields
		args           args
		want           *model.PostConnection
		wantErr        bool
	}{
		{
			name: "after is not a number",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{},
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					return mocks.NewMockPostRepo(c)
				},
			},
			args: args{
				ctx:   context.Background(),
				limit: nil,
				after: func() *string {
					v := "abc"
					return &v
				}(),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "db err",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					DefaultPostsLimit: 10,
				},
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPosts(gomock.Any(), 10, 0).Return(nil, false, fmt.Errorf("db error"))
					return pr
				},
			},
			args: args{
				ctx:   context.Background(),
				limit: nil,
				after: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "cfg value",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					DefaultPostsLimit: 10,
					MaxPostsLimit:     10,
				},
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPosts(gomock.Any(), 10, 0).Return([]*models.Post{}, false, nil)
					return pr
				},
			},
			args: args{
				ctx:   context.Background(),
				limit: nil,
				after: nil,
			},
			want: &model.PostConnection{
				Edges: []*model.PostEdge{},
				PageInfo: &model.PageInfo{
					StartCursor: func() *string { s := ""; return &s }(),
					EndCursor:   func() *string { s := ""; return &s }(),
					HasNextPage: false,
				},
			},
			wantErr: false,
		},
		{
			name: "limit provided",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					DefaultPostsLimit: 5,
					MaxPostsLimit:     10,
				},
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPosts(gomock.Any(), 3, 0).Return([]*models.Post{}, false, nil)
					return pr
				},
			},
			args: args{
				ctx: context.Background(),
				limit: func() *int32 {
					v := int32(3)
					return &v
				}(),
				after: nil,
			},
			want: &model.PostConnection{
				Edges: []*model.PostEdge{},
				PageInfo: &model.PageInfo{
					StartCursor: func() *string { s := ""; return &s }(),
					EndCursor:   func() *string { s := ""; return &s }(),
					HasNextPage: false,
				},
			},
			wantErr: false,
		},
		{
			name: "after provided",
			resolverFields: resolverFields{
				cfg: cfg.Cfg{
					DefaultPostsLimit: 5,
					MaxPostsLimit:     10,
				},
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPosts(gomock.Any(), 5, 2).Return([]*models.Post{}, false, nil)
					return pr
				},
			},
			args: args{
				ctx:   context.Background(),
				limit: nil,
				after: func() *string {
					v := "2"
					return &v
				}(),
			},
			want: &model.PostConnection{
				Edges: []*model.PostEdge{},
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
					DefaultPostsLimit: 10,
					MaxPostsLimit:     10,
				},
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPosts(gomock.Any(), 10, 0).Return([]*models.Post{
						{
							ID:              1,
							Title:           "title1",
							Text:            "text1",
							CommentsAllowed: true,
						},
						{
							ID:              2,
							Title:           "title2",
							Text:            "text2",
							CommentsAllowed: false,
						},
					}, true, nil)
					return pr
				},
			},
			args: args{
				ctx:   context.Background(),
				limit: nil,
				after: nil,
			},
			want: &model.PostConnection{
				Edges: []*model.PostEdge{
					{
						Cursor: "1",
						Node: &model.Post{
							ID:              "1",
							Title:           "title1",
							Text:            "text1",
							CommentsAllowed: true,
						},
					},
					{
						Cursor: "2",
						Node: &model.Post{
							ID:              "2",
							Title:           "title2",
							Text:            "text2",
							CommentsAllowed: false,
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
			r := &queryResolver{
				Resolver: &Resolver{
					Cfg:      tt.resolverFields.cfg,
					PostRepo: tt.resolverFields.getPostRepo(c),
				},
			}
			got, err := r.Posts(tt.args.ctx, tt.args.limit, tt.args.after)
			if (err != nil) != tt.wantErr {
				t.Errorf("Posts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Posts() got = %v, want %v", got, tt.want)
			}
		})
	}
}
