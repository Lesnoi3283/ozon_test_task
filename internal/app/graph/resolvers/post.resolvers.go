package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
	"strconv"
)

type PostResolver struct{ *Resolver }

func (p *PostResolver) Comments(ctx context.Context, obj *model.Post, limit *int, after *string) (*model.CommentConnection, error) {
	//data prepare
	id, err := strconv.Atoi(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("postID is not an int")
	}
	if limit == nil {
		limit = new(int)
	}
	var afterInt int
	if after == nil {
		afterInt = 0
	} else {
		afterInt, err = strconv.Atoi(*after)
		if err != nil {
			return nil, fmt.Errorf("after is not convertable to int")
		}
	}

	//get data
	comments, hasNextPage, err := p.CommentRepo.GetCommentsByPostID(ctx, id, *limit, afterInt)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by post id: %w", err)
	}

	var startCursor string
	var endCursor string
	if len(comments) > 0 {
		startCursor = strconv.Itoa(comments[0].ID)
		endCursor = strconv.Itoa(comments[len(comments)-1].ID)
	}

	//prepare answer
	edges := make([]*model.CommentEdge, len(comments))
	for i, comment := range comments {
		edges[i] = &model.CommentEdge{
			Cursor: strconv.Itoa(comment.ID),
			Node: &model.Comment{
				ID:        strconv.Itoa(comment.ID),
				Text:      comment.Text,
				CreatedAt: comment.CreatedAt.String(),
			},
		}
	}

	// return the answer
	return &model.CommentConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			StartCursor: &startCursor,
			EndCursor:   &endCursor,
			HasNextPage: hasNextPage,
		},
	}, nil
}
