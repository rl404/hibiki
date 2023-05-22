package repository

import (
	"context"

	"github.com/rl404/hibiki/internal/domain/author/entity"
)

// Repository contains functions for author domain.
type Repository interface {
	BatchUpdate(ctx context.Context, data []entity.Author) (int, error)
	GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.Author, int, int, error)
}
