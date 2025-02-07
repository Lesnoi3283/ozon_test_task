package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
	"strconv"
)

type PostResolver struct{ *Resolver }

func (p *PostResolver) Comments(ctx context.Context, obj *model.Post) (*model.CommentConnection, error) {
	id, err := strconv.Atoi(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("postID is not an int")
	}
	comments, hasNextPage, err := p.CommentRepo.GetCommentsByPostID(ctx, id, p.Cfg.DefaultCommentsLimit, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get comments by post id: %w", err)
	}

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

	var startCursor string
	var endCursor string
	if len(comments) > 0 {
		startCursor = strconv.Itoa(comments[0].ID)
		endCursor = strconv.Itoa(comments[len(comments)-1].ID)
	}

	return &model.CommentConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			StartCursor: &startCursor,
			EndCursor:   &endCursor,
			HasNextPage: hasNextPage,
		},
	}, nil
}

//todo написать отдельный резолвер для реплаев, возвращать только 1 уровень реплаев (без сабреплаев).
