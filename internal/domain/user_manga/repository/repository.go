package repository

import (
	"context"

	"github.com/rl404/hibiki/internal/domain/user_manga/entity"
)

// Repository contains functions for user_manga domain.
type Repository interface {
	Get(ctx context.Context, data entity.GetUserMangaRequest) ([]*entity.UserManga, int, int, error)
	BatchUpdate(ctx context.Context, data []entity.UserManga) (int, error)
	IsOld(ctx context.Context, username string) (bool, int, error)
	GetOldUsernames(ctx context.Context) ([]string, int, error)
	DeleteNotInList(ctx context.Context, username string, ids []int64) (int, error)
	DeleteByMangaID(ctx context.Context, mangaID int64) (int, error)
}
