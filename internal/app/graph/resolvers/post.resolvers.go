package resolvers

import (
	"context"
	"ozon_test_task/internal/app/graph/model"
)

type PostResolver struct{ *Resolver }

func (p *PostResolver) Comments(ctx context.Context, obj *model.Post) (*model.User, error) {
	//owner, err := p.
	//todo
}
