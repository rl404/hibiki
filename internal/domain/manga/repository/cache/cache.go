package cache

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/hibiki/internal/domain/manga/entity"
	"github.com/rl404/hibiki/internal/domain/manga/repository"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/utils"
)

// Cache contains functions for manga cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new manga cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// GetByID to get manga by id.
func (c *Cache) GetByID(ctx context.Context, id int64) (data *entity.Manga, code int, err error) {
	key := utils.GetKey("manga", id)
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

// Update to update manga.
func (c *Cache) Update(ctx context.Context, data entity.Manga) (int, error) {
	if code, err := c.repo.Update(ctx, data); err != nil {
		return code, errors.Wrap(ctx, err)
	}

	key := utils.GetKey("manga", data.ID)
	if err := c.cacher.Delete(ctx, key); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalCache, err)
	}

	return http.StatusOK, nil
}

// DeleteByID to delete manga by id.
func (c *Cache) DeleteByID(ctx context.Context, id int64) (int, error) {
	return c.repo.DeleteByID(ctx, id)
}

// IsOld to check if old.
func (c *Cache) IsOld(ctx context.Context, id int64) (bool, int, error) {
	return c.repo.IsOld(ctx, id)
}

// GetMaxID to get max id.
func (c *Cache) GetMaxID(ctx context.Context) (int64, int, error) {
	return c.repo.GetMaxID(ctx)
}

// GetIDs to get ids.
func (c *Cache) GetIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetIDs(ctx)
}

// GetOldFinishedIDs to get old finished manga ids.
func (c *Cache) GetOldFinishedIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetOldFinishedIDs(ctx)
}

// GetOldReleasingIDs to get old releasing manga ids.
func (c *Cache) GetOldReleasingIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetOldReleasingIDs(ctx)
}

// GetOldNotYetIDs to get old not yet released manga ids.
func (c *Cache) GetOldNotYetIDs(ctx context.Context) ([]int64, int, error) {
	return c.repo.GetOldNotYetIDs(ctx)
}

type getAllCache struct {
	Data  []entity.Manga
	Total int
}

// GetAll to get manga list.
func (c *Cache) GetAll(ctx context.Context, req entity.GetAllRequest) (_ []entity.Manga, _ int, code int, err error) {
	key := utils.GetKey("manga", utils.QueryToKey(req))

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
