package service

import (
	"context"
	"net/http"

	"github.com/rl404/hibiki/internal/domain/genre/entity"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/utils"
)

// GetGenresRequest is get genres request model.
type GetGenresRequest struct {
	Name  string `validate:"omitempty,gte=3" mod:"trim,lcase"`
	Page  int    `validate:"required,gte=1" mod:"default=1"`
	Limit int    `validate:"required,gte=-1" mod:"default=20"`
}

// GetGenres to get genre list.
func (s *service) GetGenres(ctx context.Context, data GetGenresRequest) ([]genre, *pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	genres, total, code, err := s.genre.GetAll(ctx, entity.GetAllRequest{
		Name:  data.Name,
		Page:  data.Page,
		Limit: data.Limit,
	})
	if err != nil {
		return nil, nil, code, errors.Wrap(ctx, err)
	}

	res := make([]genre, len(genres))
	for i, a := range genres {
		res[i] = genre{
			ID:   a.ID,
			Name: a.Name,
		}
	}

	return res, &pagination{
		Page:  data.Page,
		Limit: data.Limit,
		Total: total,
	}, http.StatusOK, nil
}
