package service

import (
	"context"
	"net/http"

	"github.com/rl404/hibiki/internal/domain/author/entity"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/utils"
)

type author struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// GetAuthorsRequest is get authors request model.
type GetAuthorsRequest struct {
	Name  string `validate:"omitempty,gte=3" mod:"trim,lcase"`
	Page  int    `validate:"required,gte=1" mod:"default=1"`
	Limit int    `validate:"required,gte=-1" mod:"default=20"`
}

// GetAuthors to get author list.
func (s *service) GetAuthors(ctx context.Context, data GetAuthorsRequest) ([]author, *pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, errors.Wrap(ctx, err)
	}

	authors, total, code, err := s.author.GetAll(ctx, entity.GetAllRequest{
		Name:  data.Name,
		Page:  data.Page,
		Limit: data.Limit,
	})
	if err != nil {
		return nil, nil, code, errors.Wrap(ctx, err)
	}

	res := make([]author, len(authors))
	for i, a := range authors {
		res[i] = author{
			ID:        a.ID,
			FirstName: a.FirstName,
			LastName:  a.LastName,
		}
	}

	return res, &pagination{
		Page:  data.Page,
		Limit: data.Limit,
		Total: total,
	}, http.StatusOK, nil
}
