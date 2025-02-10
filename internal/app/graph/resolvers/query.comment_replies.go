package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
	"strconv"
)

// CommentReplies is the resolver for the commentReplies field.
func (r *queryResolver) CommentReplies(ctx context.Context, commentID string, limit *int32, after *string) (*model.CommentConnection, error) {
	//check input
	limitInt := 0
	if limit != nil {
		limitInt = int(*limit)
	} else {
		limitInt = r.Cfg.DefaultCommentsLimit
		if limitInt > r.Cfg.MaxCommentsLimit {
			limitInt = r.Cfg.MaxCommentsLimit
		}
	}

	commentIDInt, err := strconv.Atoi(commentID)
	if err != nil {
		r.Logger.Debugf("cant convert commentID to int: %v", err)
		return nil, fmt.Errorf("commentID is not int")
	}

	afterInt := 0
	if after != nil {
		afterInt, err = strconv.Atoi(*after)
		if err != nil {
			r.Logger.Debugf("cant convert after to int: %v", err)
			return nil, fmt.Errorf("after is not int")
		}
	}

	//get replies
	replays, hasNextPage, err := r.CommentRepo.GetReplaysByCommentID(ctx, commentIDInt, limitInt, afterInt)
	if err != nil {
		r.Logger.Debugf("failed to get replays from db: %v", err)
		return nil, fmt.Errorf("interal server error")
	}

	edges := make([]*model.CommentEdge, len(replays))
	for i, replay := range replays {
		edges[i] = &model.CommentEdge{
			Cursor: strconv.Itoa(replay.ID),
			Node: &model.Comment{
				ID: strconv.Itoa(replay.ID),
				Owner: &model.User{
					ID:       strconv.Itoa(replay.Owner.ID),
					Username: replay.Owner.Login,
				},
				Text:      replay.Text,
				CreatedAt: replay.CreatedAt.String(),
				Replies:   nil,
			},
		}
	}

	startCoursor := ""
	endCoursor := ""

	if len(edges) > 0 {
		startCoursor = edges[0].Node.ID
		endCoursor = edges[len(edges)-1].Node.ID
	}

	return &model.CommentConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			StartCursor: &startCoursor,
			EndCursor:   &endCoursor,
			HasNextPage: hasNextPage,
		},
	}, nil
}
