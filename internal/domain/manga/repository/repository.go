package repository

import (
	"context"

	"github.com/rl404/hibiki/internal/domain/manga/entity"
)

// Repository contains functions for manga details.
type Repository interface {
	GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.Manga, int, int, error)
	GetByID(ctx context.Context, id int64) (*entity.Manga, int, error)
	Update(ctx context.Context, data entity.Manga) (int, error)
	DeleteByID(ctx context.Context, id int64) (int, error)

	IsOld(ctx context.Context, id int64) (bool, int, error)
	GetMaxID(ctx context.Context) (int64, int, error)
	GetIDs(ctx context.Context) ([]int64, int, error)
	GetOldFinishedIDs(ctx context.Context) ([]int64, int, error)
	GetOldReleasingIDs(ctx context.Context) ([]int64, int, error)
	GetOldNotYetIDs(ctx context.Context) ([]int64, int, error)
}
