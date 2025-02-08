package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
	"strconv"
)

//todo: rename "first" to "limit"

// Posts is the resolver for the posts field.
func (r *queryResolver) Posts(ctx context.Context, limit *int32, after *string) (*model.PostConnection, error) {
	//prepare input data
	if limit == nil {
		limit = new(int32)
	}
	if after == nil {
		after = new(string)
	}
	afterInt, err := strconv.Atoi(*after)
	if err != nil {
		return nil, fmt.Errorf("after is not a number")
	}

	//get posts
	posts, hasNextPage, err := r.PostRepo.GetPosts(ctx, int(*limit), afterInt)
	if err != nil {
		return nil, fmt.Errorf("cant get posts")
	}

	edges := make([]*model.PostEdge, len(posts))
	for i, post := range posts {
		edges[i] = &model.PostEdge{
			Cursor: strconv.Itoa(post.ID),
			Node: &model.Post{
				ID:              strconv.Itoa(post.ID),
				Title:           post.Title,
				Text:            post.Text,
				CommentsAllowed: post.CommentsAllowed,
			},
		}
	}

	startCursor := ""
	endCursor := ""
	if len(edges) > 0 {
		startCursor = edges[0].Cursor
		endCursor = edges[len(edges)-1].Cursor
	}

	return &model.PostConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			StartCursor: &startCursor,
			EndCursor:   &endCursor,
			HasNextPage: hasNextPage,
		},
	}, nil
}
