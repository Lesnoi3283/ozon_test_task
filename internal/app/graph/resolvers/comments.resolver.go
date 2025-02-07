package resolvers

import (
	"context"
	"fmt"
	"ozon_test_task/internal/app/graph/model"
	"strconv"
)

type CommentsResolver struct {
	*Resolver
}

func (r *CommentsResolver) Replies(ctx context.Context, obj *model.Comment, limit *int, after *string) (*model.CommentConnection, error) {
	//data prepare
	id, err := strconv.Atoi(obj.ID)
	if err != nil {
		return nil, fmt.Errorf("commentID is not an int")
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
	replays, hasNextPage, err := r.CommentRepo.GetReplaysByCommentID(ctx, id, *limit, afterInt)
	if err != nil {
		return nil, fmt.Errorf("failed to get replays")
	}

	var startCursor string
	var endCursor string
	if len(replays) > 0 {
		startCursor = strconv.Itoa(replays[0].ID)
		endCursor = strconv.Itoa(replays[len(replays)-1].ID)
	}

	//prepare answer
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
	
	return &model.CommentConnection{
		Edges: edges,
		PageInfo: &model.PageInfo{
			StartCursor: &startCursor,
			EndCursor:   &endCursor,
			HasNextPage: hasNextPage,
		},
	}, nil
}
