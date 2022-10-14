package repository

import (
	"context"

	"github.com/rl404/hibiki/internal/domain/nagato/entity"
	"github.com/rl404/nagato"
)

// Repository contains functions for nagato domain.
type Repository interface {
	GetMangaByID(ctx context.Context, id int) (*nagato.Manga, int, error)
	GetUserManga(ctx context.Context, data entity.GetUserMangaRequest) ([]nagato.UserManga, int, error)
}
