package cache

import (
	"context"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/hibiki/internal/domain/author/entity"
	"github.com/rl404/hibiki/internal/domain/author/repository"
)

// Cache contains functions for author cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new author cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// BatchUpdate to batch update.
func (c *Cache) BatchUpdate(ctx context.Context, data []entity.Author) (int, error) {
	return c.repo.BatchUpdate(ctx, data)
}
