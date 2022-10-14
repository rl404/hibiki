package mongo

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/hibiki/internal/domain/genre/entity"
	"github.com/rl404/hibiki/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo contains functions for genre mongodb.
type Mongo struct {
	db *mongo.Collection
}

// New to create new genre database.
func New(mdb *mongo.Database) *Mongo {
	return &Mongo{
		db: mdb.Collection("genre"),
	}
}

// BatchUpdate to batch update.
func (m *Mongo) BatchUpdate(ctx context.Context, data []entity.Genre) (int, error) {
	ids := make([]int64, len(data))
	for i, g := range data {
		ids[i] = g.ID
	}

	cursor, err := m.db.Find(ctx, bson.M{"id": bson.M{"$in": ids}})
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	defer cursor.Close(ctx)

	existMap := make(map[int64]bool)
	for cursor.Next(ctx) {
		var genre genre
		if err := cursor.Decode(&genre); err != nil {
			return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
		}

		existMap[genre.ID] = true
	}

	models := make([]mongo.WriteModel, len(data))
	for i, d := range data {
		if _, ok := existMap[d.ID]; ok {
			models[i] = mongo.NewUpdateOneModel().
				SetFilter(bson.M{"id": d.ID}).
				SetUpdate(bson.M{"$set": bson.M{"name": d.Name, "updated_at": time.Now()}})
		} else {
			models[i] = mongo.NewInsertOneModel().SetDocument(m.fromEntity(d))
		}
	}

	if _, err := m.db.BulkWrite(ctx, models); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return http.StatusOK, nil
}
