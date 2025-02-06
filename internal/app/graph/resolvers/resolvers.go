package resolvers

import (
	"ozon_test_task/internal/app/graph"
	"ozon_test_task/internal/app/graph/repository"
)

// Mutation returns graph.MutationResolver implementation.
func (r *Resolver) Mutation(postRepo repository.PostRepo, userRepo repository.UserRepo) graph.MutationResolver {
	return &mutationResolver{r}
}

// Query returns graph.QueryResolver implementation.
func (r *Resolver) Query() graph.QueryResolver { return &queryResolver{r} }

type mutationResolver struct {
	*Resolver
}

type queryResolver struct{ *Resolver }
