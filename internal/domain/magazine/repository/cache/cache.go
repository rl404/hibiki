package cache

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/domain/magazine/entity"
	"github.com/rl404/hibiki/internal/domain/magazine/repository"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/utils"
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

type getAllCache struct {
	Data  []entity.Magazine
	Total int
}

// GetAll to get magazine list.
func (c *Cache) GetAll(ctx context.Context, req entity.GetAllRequest) (_ []entity.Magazine, _ int, code int, err error) {
	key := utils.GetKey("magazine", utils.QueryToKey(req))

	var data getAllCache
	if c.cacher.Get(ctx, key, &data) == nil {
		return data.Data, data.Total, http.StatusOK, nil
	}

	data.Data, data.Total, code, err = c.repo.GetAll(ctx, req)
	if err != nil {
		return nil, 0, code, stack.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}

	return data.Data, data.Total, code, nil
}

// GetByID to get by id.
func (c *Cache) GetByID(ctx context.Context, id int64) (data *entity.Magazine, code int, err error) {
	key := utils.GetKey("magazine", id)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetByID(ctx, id)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalCache)
	}

	return data, code, nil
}
