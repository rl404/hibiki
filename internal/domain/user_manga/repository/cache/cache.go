package cache

import (
	"context"

	"github.com/rl404/fairy/cache"
	"github.com/rl404/hibiki/internal/domain/user_manga/entity"
	"github.com/rl404/hibiki/internal/domain/user_manga/repository"
)

// Cache contains functions for user_manga cache.
type Cache struct {
	cacher cache.Cacher
	repo   repository.Repository
}

// New to create new user_manga cache.
func New(cacher cache.Cacher, repo repository.Repository) *Cache {
	return &Cache{
		cacher: cacher,
		repo:   repo,
	}
}

// Get to get user manga.
func (c *Cache) Get(ctx context.Context, data entity.GetUserMangaRequest) ([]*entity.UserManga, int, int, error) {
	return c.repo.Get(ctx, data)
}

// DeleteByMangaID to delete by manga id.
func (c *Cache) DeleteByMangaID(ctx context.Context, mangaID int64) (int, error) {
	return c.repo.DeleteByMangaID(ctx, mangaID)
}

// GetOldUsernames to get old username.
func (c *Cache) GetOldUsernames(ctx context.Context) ([]string, int, error) {
	return c.repo.GetOldUsernames(ctx)
}

// IsOld to check if old.
func (c *Cache) IsOld(ctx context.Context, username string) (bool, int, error) {
	return c.repo.IsOld(ctx, username)
}

// DeleteNotInList to delete not in list.
func (c *Cache) DeleteNotInList(ctx context.Context, username string, ids []int64) (int, error) {
	return c.repo.DeleteNotInList(ctx, username, ids)
}

// BatchUpdate to batch update.
func (c *Cache) BatchUpdate(ctx context.Context, data []entity.UserManga) (int, error) {
	return c.repo.BatchUpdate(ctx, data)
}
