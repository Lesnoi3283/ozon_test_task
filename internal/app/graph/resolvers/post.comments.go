package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
	"strconv"
)

//type postResolver struct{ *Resolver }

func (p *postResolver) Comments(ctx context.Context, obj *model.Post, limit *int32, after *string) (*model.CommentConnection, error) {
	//data prepare
	id, err := strconv.Atoi(obj.ID)
	if err != nil {
		p.Logger.Debugf("cant convert postID to int, err: %v", err)
		return nil, fmt.Errorf("postID is not an int")
	}
	limitInt := 0
	if limit == nil {
		limitInt = p.Cfg.DefaultCommentsLimit
	} else {
		limitInt = int(*limit)
		if limitInt > p.Cfg.MaxCommentsLimit {
			limitInt = p.Cfg.MaxCommentsLimit
		}
	}
	var afterInt int
	if after == nil {
		afterInt = 0
	} else {
		afterInt, err = strconv.Atoi(*after)
		if err != nil {
			p.Logger.Debugf("cant convert after to int, err: %v", err)
			return nil, fmt.Errorf("after is not convertable to int")
		}
	}

	//get data
	comments, hasNextPage, err := p.CommentRepo.GetCommentsByPostID(ctx, id, limitInt, afterInt)
	if err != nil {
		p.Logger.Debugf("cant get comments from db, err: %v", err)
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
				Owner: &model.User{
					ID:       strconv.Itoa(comment.Owner.ID),
					Username: comment.Owner.Login,
				},
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
