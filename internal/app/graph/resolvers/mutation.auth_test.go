package resolvers

import (
	"context"
	"fmt"
	"go.uber.org/mock/gomock"
	"ozon_test_task/internal/app/graph/model"
	"ozon_test_task/internal/app/graph/repository"
	"ozon_test_task/internal/app/graph/repository/mocks"
	"ozon_test_task/internal/app/middlewares"
	mwmocks "ozon_test_task/internal/app/middlewares/mocks"
	"ozon_test_task/internal/app/models"
	"ozon_test_task/pkg/authUtils"
	"reflect"
	"testing"
)

func Test_mutationResolver_Auth(t *testing.T) {
	type args struct {
		ctx      context.Context
		username string
		password string
	}
	type resolverFields struct {
		getUserRepo   func(c *gomock.Controller) repository.UserRepo
		getJWTManager func(c *gomock.Controller) middlewares.JWTManager
	}
	tests := []struct {
		name           string
		resolverFields resolverFields
		args           args
		want           *model.AuthResponse
		wantErr        bool
	}{
		{
			name: "Username is empty",
			resolverFields: resolverFields{
				getUserRepo: func(c *gomock.Controller) repository.UserRepo {
					return mocks.NewMockUserRepo(c)
				},
				getJWTManager: func(c *gomock.Controller) middlewares.JWTManager {
					return mwmocks.NewMockJWTManager(c)
				},
			},
			args: args{
				ctx:      context.Background(),
				username: "",
				password: "somepass",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Password is empty",
			resolverFields: resolverFields{
				getUserRepo: func(c *gomock.Controller) repository.UserRepo {
					return mocks.NewMockUserRepo(c)
				},
				getJWTManager: func(c *gomock.Controller) middlewares.JWTManager {
					return mwmocks.NewMockJWTManager(c)
				},
			},
			args: args{
				ctx:      context.Background(),
				username: "someuser",
				password: "",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "User not found in DB",
			resolverFields: resolverFields{
				getUserRepo: func(c *gomock.Controller) repository.UserRepo {
					ur := mocks.NewMockUserRepo(c)
					ur.EXPECT().GetUserByLoginWithCred(gomock.Any(), "notfound").Return(nil, repository.NewErrNotFound())
					return ur
				},
				getJWTManager: func(c *gomock.Controller) middlewares.JWTManager {
					return mwmocks.NewMockJWTManager(c)
				},
			},
			args: args{
				ctx:      context.Background(),
				username: "notfound",
				password: "somepass",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "DB error",
			resolverFields: resolverFields{
				getUserRepo: func(c *gomock.Controller) repository.UserRepo {
					ur := mocks.NewMockUserRepo(c)
					ur.EXPECT().GetUserByLoginWithCred(gomock.Any(), "someuser").Return(nil, fmt.Errorf("db error"))
					return ur
				},
				getJWTManager: func(c *gomock.Controller) middlewares.JWTManager {
					return mwmocks.NewMockJWTManager(c)
				},
			},
			args: args{
				ctx:      context.Background(),
				username: "someuser",
				password: "somepass",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Wrong password",
			resolverFields: resolverFields{
				getUserRepo: func(c *gomock.Controller) repository.UserRepo {
					ur := mocks.NewMockUserRepo(c)
					ur.EXPECT().GetUserByLoginWithCred(gomock.Any(), "someuser").Return(&models.User{
						ID:           1,
						Login:        "someuser",
						PasswordHash: "hash",
						PasswordSalt: "salt",
					}, nil)
					return ur
				},
				getJWTManager: func(c *gomock.Controller) middlewares.JWTManager {
					return mwmocks.NewMockJWTManager(c)
				},
			},
			args: args{
				ctx:      context.Background(),
				username: "someuser",
				password: "wrongpass",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Ok",
			resolverFields: resolverFields{
				getUserRepo: func(c *gomock.Controller) repository.UserRepo {
					ur := mocks.NewMockUserRepo(c)
					passwordHash := authUtils.HashPassword("somepass", "salt")
					ur.EXPECT().GetUserByLoginWithCred(gomock.Any(), "someuser").Return(&models.User{
						ID:           1,
						Login:        "someuser",
						PasswordHash: passwordHash,
						PasswordSalt: "salt",
					}, nil)
					return ur
				},
				getJWTManager: func(c *gomock.Controller) middlewares.JWTManager {
					jm := mwmocks.NewMockJWTManager(c)
					jm.EXPECT().BuildNewJWTString(1).Return("token123", nil)
					return jm
				},
			},
			args: args{
				ctx:      context.Background(),
				username: "someuser",
				password: "somepass",
			},
			want: &model.AuthResponse{
				Token: "token123",
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
					UserRepo:   tt.resolverFields.getUserRepo(c),
					JWTManager: tt.resolverFields.getJWTManager(c),
				},
			}
			got, err := r.Auth(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Auth() got = %v, want %v", got, tt.want)
			}
		})
	}
}
