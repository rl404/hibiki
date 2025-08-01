package mongo

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/domain/manga_stats_history/entity"
	"github.com/rl404/hibiki/internal/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Mongo contains functions for manga_stats_history mongodb.
type Mongo struct {
	db *mongo.Collection
}

// New to create new manga_stats_history database.
func New(db *mongo.Database) *Mongo {
	return &Mongo{
		db: db.Collection("manga_stats_history"),
	}
}

// Create to create new manga stats history.
func (m *Mongo) Create(ctx context.Context, data entity.MangaStatsHistory) (int, error) {
	if _, err := m.db.InsertOne(ctx, m.fromEntity(data)); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}
