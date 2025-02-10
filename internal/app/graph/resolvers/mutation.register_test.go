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
	"reflect"
	"testing"
)

func Test_mutationResolver_Register(t *testing.T) {
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
			name: "DB error",
			resolverFields: resolverFields{
				getUserRepo: func(c *gomock.Controller) repository.UserRepo {
					ur := mocks.NewMockUserRepo(c)
					ur.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(0, fmt.Errorf("db error"))
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
			name: "Ok",
			resolverFields: resolverFields{
				getUserRepo: func(c *gomock.Controller) repository.UserRepo {
					ur := mocks.NewMockUserRepo(c)
					ur.EXPECT().AddUser(gomock.Any(), gomock.Any()).Return(123, nil)
					return ur
				},
				getJWTManager: func(c *gomock.Controller) middlewares.JWTManager {
					jm := mwmocks.NewMockJWTManager(c)
					jm.EXPECT().BuildNewJWTString(123).Return("token123", nil)
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
			got, err := r.Register(tt.args.ctx, tt.args.username, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Register() got = %v, want %v", got, tt.want)
			}
		})
	}
}
