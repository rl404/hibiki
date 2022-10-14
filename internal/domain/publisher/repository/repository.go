package repository

import (
	"context"

	"github.com/rl404/hibiki/internal/domain/publisher/entity"
)

// Repository contains functions for publisher domain.
type Repository interface {
	PublishParseManga(ctx context.Context, data entity.ParseMangaRequest) error
	PublishParseUserManga(ctx context.Context, data entity.ParseUserMangaRequest) error
}
