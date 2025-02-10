package resolvers

import (
	"context"
	"ozon_test_task/internal/app/graph/model"
	"reflect"
	"testing"
)

func Test_commentResolver_Replies(t *testing.T) {
	type fields struct {
		Resolver *Resolver
	}
	type args struct {
		ctx   context.Context
		obj   *model.Comment
		limit *int32
		after *string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.CommentConnection
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &commentResolver{
				Resolver: tt.fields.Resolver,
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
