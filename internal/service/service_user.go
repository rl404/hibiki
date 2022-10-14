package service

import (
	"context"
	"net/http"
	"time"

	publisherEntity "github.com/rl404/hibiki/internal/domain/publisher/entity"
	"github.com/rl404/hibiki/internal/domain/user_manga/entity"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/utils"
)

// UserManga is user manga model.
type UserManga struct {
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
func (s *service) GetUserManga(ctx context.Context, data GetUserMangaRequest) ([]UserManga, *Pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	userManga, cnt, code, err := s.userManga.Get(ctx, entity.GetUserMangaRequest{
		Username: data.Username,
		Page:     data.Page,
		Limit:    data.Limit,
	})
	if err != nil {
		return nil, nil, code, errors.Wrap(ctx, err)
	}

	if cnt == 0 {
		// Queue to parse.
		if err := s.publisher.PublishParseUserManga(ctx, publisherEntity.ParseUserMangaRequest{Username: data.Username}); err != nil {
			return nil, nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
		}
		return nil, nil, http.StatusAccepted, nil
	}

	res := make([]UserManga, len(userManga))
	for i, ua := range userManga {
		res[i] = UserManga{
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

	return res, &Pagination{
		Page:  data.Page,
		Limit: data.Limit,
		Total: cnt,
	}, http.StatusOK, nil
}
