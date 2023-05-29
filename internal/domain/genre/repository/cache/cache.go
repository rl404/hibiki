package cache

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/hibiki/internal/domain/genre/entity"
	"github.com/rl404/hibiki/internal/domain/genre/repository"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/utils"
)

// Cache contains functions for genre cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new genre cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// BatchUpdate to batch update.
func (c *Cache) BatchUpdate(ctx context.Context, data []entity.Genre) (int, error) {
	return c.repo.BatchUpdate(ctx, data)
}

type getAllCache struct {
	Data  []entity.Genre
	Total int
}

// GetAll to get genre list.
func (c *Cache) GetAll(ctx context.Context, req entity.GetAllRequest) (_ []entity.Genre, _ int, code int, err error) {
	key := utils.GetKey("genre", utils.QueryToKey(req))

	var data getAllCache
	if c.cacher.Get(ctx, key, &data) == nil {
		return data.Data, data.Total, http.StatusOK, nil
	}

	data.Data, data.Total, code, err = c.repo.GetAll(ctx, req)
	if err != nil {
		return nil, 0, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data.Data, data.Total, code, nil
}

// GetByID to get by id.
func (c *Cache) GetByID(ctx context.Context, id int64) (data *entity.Genre, code int, err error) {
	key := utils.GetKey("genre", id)
	if c.cacher.Get(ctx, key, &data) == nil {
		return data, http.StatusOK, nil
	}

	data, code, err = c.repo.GetByID(ctx, id)
	if err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	if err := c.cacher.Set(ctx, key, data); err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return data, code, nil
}
