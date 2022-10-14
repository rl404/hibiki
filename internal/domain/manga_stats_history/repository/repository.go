package repository

import (
	"context"

	"github.com/rl404/hibiki/internal/domain/manga_stats_history/entity"
)

// Repository contains functions for manga_stats_history domain.
type Repository interface {
	Create(ctx context.Context, data entity.MangaStatsHistory) (int, error)
}
