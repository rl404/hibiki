package cache

import (
	"context"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/hibiki/internal/domain/magazine/entity"
	"github.com/rl404/hibiki/internal/domain/magazine/repository"
)

// Cache contains functions for magazine cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new magazine cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// BatchUpdate to batch update.
func (c *Cache) BatchUpdate(ctx context.Context, data []entity.Magazine) (int, error) {
	return c.repo.BatchUpdate(ctx, data)
}
