package mongo

import (
	"context"
	_errors "errors"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo contains functions for empty_id mongodb.
type Mongo struct {
	db *mongo.Collection
}

// New to create new empty_id database.
func New(mdb *mongo.Database) *Mongo {
	return &Mongo{
		db: mdb.Collection("empty_id"),
	}
}

// Get to get empty id.
func (m *Mongo) Get(ctx context.Context, id int64) (int64, int, error) {
	var emptyID emptyID
	if err := m.db.FindOne(ctx, bson.M{"manga_id": id}).Decode(&emptyID); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return 0, http.StatusNotFound, stack.Wrap(ctx, err, errors.ErrMangaNotFound)
		}
		return 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return emptyID.MangaID, http.StatusOK, nil
}

// Create to create empty id.
func (m *Mongo) Create(ctx context.Context, id int64) (int, error) {
	if _, err := m.db.InsertOne(ctx, &emptyID{
		MangaID: id,
	}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusCreated, nil
}

// Delete to empty id.
func (m *Mongo) Delete(ctx context.Context, id int64) (int, error) {
	if _, err := m.db.DeleteOne(ctx, bson.M{"manga_id": id}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}

// GetIDs to all empty id.
func (m *Mongo) GetIDs(ctx context.Context) ([]int64, int, error) {
	var ids []int64
	c, err := m.db.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"manga_id": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	defer c.Close(ctx)

	for c.Next(ctx) {
		var emptyID emptyID
		if err := c.Decode(&emptyID); err != nil {
			return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
		}
		ids = append(ids, emptyID.MangaID)
	}

	return ids, http.StatusOK, nil
}
