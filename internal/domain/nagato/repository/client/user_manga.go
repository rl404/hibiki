package client

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/domain/nagato/entity"
	"github.com/rl404/nagato"
)

// GetUserManga to get user manga.
func (c *Client) GetUserManga(ctx context.Context, data entity.GetUserMangaRequest) ([]nagato.UserManga, int, error) {
	manga, code, err := c.client.GetUserMangaListWithContext(ctx, nagato.GetUserMangaListParam{
		Username: data.Username,
		Status:   data.Status,
		Sort:     data.Sort,
		Limit:    data.Limit,
		Offset:   data.Offset,
	},
		nagato.MangaFieldUserStatus(
			nagato.UserMangaNumTimesReread,
			nagato.UserMangaRereadValue,
			nagato.UserMangaTags,
			nagato.UserMangaComments,
		),
	)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	return manga, http.StatusOK, nil
}
