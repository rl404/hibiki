package service

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/domain/user_manga/entity"
	"github.com/rl404/hibiki/internal/utils"
)

type userManga struct {
	MangaID   int64         `json:"manga_id"`
	Title     string        `json:"title"`
	Status    entity.Status `json:"status" swaggertype:"string"`
	Score     int           `json:"score"`
	Chapter   int           `json:"chapter"`
	Volume    int           `json:"volume"`
	Tags      []string      `json:"tags"`
	Comment   string        `json:"comment"`
	UpdatedAt time.Time     `json:"updated_at"`
}

// GetUserMangaRequest is get user manga request model.
type GetUserMangaRequest struct {
	Username string `validate:"required" mod:"trim,lcase"`
	Page     int    `validate:"required,gte=1" mod:"default=1"`
	Limit    int    `validate:"required,gte=-1" mod:"default=20"`
}

// GetUserManga to get user manga.
func (s *service) GetUserManga(ctx context.Context, data GetUserMangaRequest) ([]userManga, *pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	userMangas, cnt, code, err := s.userManga.Get(ctx, entity.GetUserMangaRequest{
		Username: data.Username,
		Page:     data.Page,
		Limit:    data.Limit,
	})
	if err != nil {
		return nil, nil, code, stack.Wrap(ctx, err)
	}

	if cnt == 0 {
		// Queue to parse.
		if err := s.publisher.PublishParseUserManga(ctx, data.Username); err != nil {
			return nil, nil, http.StatusInternalServerError, stack.Wrap(ctx, err)
		}
		return nil, nil, http.StatusAccepted, nil
	}

	res := make([]userManga, len(userMangas))
	for i, ua := range userMangas {
		res[i] = userManga{
			MangaID:   ua.MangaID,
			Title:     ua.Title,
			Status:    ua.Status,
			Score:     ua.Score,
			Chapter:   ua.Chapter,
			Volume:    ua.Volume,
			Tags:      ua.Tags,
			Comment:   ua.Comment,
			UpdatedAt: ua.UpdatedAt,
		}
	}

	return res, &pagination{
		Page:  data.Page,
		Limit: data.Limit,
		Total: cnt,
	}, http.StatusOK, nil
}
