package service

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/domain/magazine/entity"
	"github.com/rl404/hibiki/internal/utils"
)

// GetMagazinesRequest is get magazines request model.
type GetMagazinesRequest struct {
	Name  string `validate:"omitempty,gte=3" mod:"trim,lcase"`
	Page  int    `validate:"required,gte=1" mod:"default=1"`
	Limit int    `validate:"required,gte=-1" mod:"default=20"`
}

// GetMagazines to get magazine list.
func (s *service) GetMagazines(ctx context.Context, data GetMagazinesRequest) ([]magazine, *pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	magazines, total, code, err := s.magazine.GetAll(ctx, entity.GetAllRequest{
		Name:  data.Name,
		Page:  data.Page,
		Limit: data.Limit,
	})
	if err != nil {
		return nil, nil, code, stack.Wrap(ctx, err)
	}

	res := make([]magazine, len(magazines))
	for i, a := range magazines {
		res[i] = magazine{
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

// GetMagazineByID to get magazine by id.
func (s *service) GetMagazineByID(ctx context.Context, id int64) (*magazine, int, error) {
	a, code, err := s.magazine.GetByID(ctx, id)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}
	return &magazine{
		ID:   a.ID,
		Name: a.Name,
	}, http.StatusOK, nil
}
