package cache

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/domain/empty_id/repository"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/utils"
)

// Cache contains functions for empty_id cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new empty_id cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// Get to get empty id.
func (c *Cache) Get(ctx context.Context, id int64) (int64, int, error) {
	key := utils.GetKey("empty-id", id)
	var data int64
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	emptyID, code, err := c.repo.Get(ctx, id)
	if err != nil {
		return 0, code, stack.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, emptyID); err != nil {
		return 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}

	return emptyID, code, nil
}

// Create to create empty id.
func (c *Cache) Create(ctx context.Context, id int64) (int, error) {
	key := utils.GetKey("empty-id", id)
	if code, err := c.repo.Create(ctx, id); err != nil {
		return code, stack.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, true); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}

	return http.StatusCreated, nil
}

// Delete to delete empty id.
func (c *Cache) Delete(ctx context.Context, id int64) (int, error) {
	key := utils.GetKey("empty-id", id)
	if code, err := c.repo.Delete(ctx, id); err != nil {
		return code, stack.Wrap(ctx, err)
	}

	if err := c.cacher.Delete(ctx, key); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}

	return http.StatusOK, nil
}

// GetIDs to get all ids.
func (c *Cache) GetIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetIDs(ctx)
}
