package repository

import (
	"context"
)

// Repository contains functions for publisher domain.
type Repository interface {
	PublishParseManga(ctx context.Context, id int64) error
	PublishParseUserManga(ctx context.Context, username string) error
}
