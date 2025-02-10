package resolvers

import (
	"context"
	"fmt"
	"go.uber.org/mock/gomock"
	"ozon_test_task/internal/app/graph/model"
	"ozon_test_task/internal/app/graph/repository"
	"ozon_test_task/internal/app/graph/repository/mocks"
	"ozon_test_task/internal/app/models"
	"reflect"
	"testing"
)

func Test_queryResolver_Post(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
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
			name: "ID is not convertable to int",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					return mocks.NewMockPostRepo(c)
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "abc",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Post not found",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPostByID(gomock.Any(), 10).Return(nil, repository.NewErrNotFound())
					return pr
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "10",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "DB error",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPostByID(gomock.Any(), 10).Return(nil, fmt.Errorf("db error"))
					return pr
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "10",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Ok",
			resolverFields: resolverFields{
				getPostRepo: func(c *gomock.Controller) repository.PostRepo {
					pr := mocks.NewMockPostRepo(c)
					pr.EXPECT().GetPostByID(gomock.Any(), 10).Return(&models.Post{
						ID:              10,
						Title:           "Test Title",
						Text:            "Test Text",
						CommentsAllowed: true,
						Owner: models.User{
							ID:    5,
							Login: "ownerUser",
						},
					}, nil)
					return pr
				},
			},
			args: args{
				ctx: context.Background(),
				id:  "10",
			},
			want: &model.Post{
				ID:    "10",
				Title: "Test Title",
				Text:  "Test Text",
				Owner: &model.User{
					ID:       "5",
					Username: "ownerUser",
				},
				CommentsAllowed: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := gomock.NewController(t)
			r := &queryResolver{
				Resolver: &Resolver{
					PostRepo: tt.resolverFields.getPostRepo(c),
				},
			}
			got, err := r.Post(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Post() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Post() got = %v, want %v", got, tt.want)
			}
		})
	}
}
