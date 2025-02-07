package resolvers

import (
	"ozon_test_task/internal/app/graph"
)

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation() graph.MutationResolver {
	return &mutationResolver{r}
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct {
	*Resolver
}

type queryResolver struct{ *Resolver }
