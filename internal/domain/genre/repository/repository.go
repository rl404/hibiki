package repository

import (
	"context"

	"github.com/rl404/hibiki/internal/domain/genre/entity"
)

// Repository contains functions for genre domain.
type Repository interface {
	BatchUpdate(ctx context.Context, data []entity.Genre) (int, error)
	GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.Genre, int, int, error)
}
