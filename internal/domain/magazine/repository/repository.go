package repository

import (
	"context"

	"github.com/rl404/hibiki/internal/domain/magazine/entity"
)

// Repository contains functions for magazine domain.
type Repository interface {
	BatchUpdate(ctx context.Context, data []entity.Magazine) (int, error)
	GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.Magazine, int, int, error)
}
