package mongo

import (
	"context"
	_errors "errors"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/domain/genre/entity"
	"github.com/rl404/hibiki/internal/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	defer cursor.Close(ctx)

	existMap := make(map[int64]bool)
	for cursor.Next(ctx) {
		var genre genre
		if err := cursor.Decode(&genre); err != nil {
			return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
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
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return http.StatusOK, nil
}

// GetAll to get all.
func (m *Mongo) GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.Genre, int, int, error) {
	filter := bson.M{}
	opt := options.Find().SetSort(bson.M{"name": 1}).SetSkip(int64((data.Page - 1) * data.Limit)).SetLimit(int64(data.Limit))

	if data.Name != "" {
		filter = bson.M{"name": bson.M{"$regex": data.Name, "$options": "i"}}
	}

	if data.Limit < 0 {
		opt.SetLimit(0)
	}

	c, err := m.db.Find(ctx, filter, opt)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	defer c.Close(ctx)

	var genres []entity.Genre
	for c.Next(ctx) {
		var genre genre
		if err := c.Decode(&genre); err != nil {
			return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
		}
		genres = append(genres, entity.Genre{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}

	total, err := m.db.CountDocuments(ctx, filter, options.Count())
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return genres, int(total), http.StatusOK, nil
}

// GetByID to get by id.
func (m *Mongo) GetByID(ctx context.Context, id int64) (*entity.Genre, int, error) {
	var genre genre
	if err := m.db.FindOne(ctx, bson.M{"id": id}).Decode(&genre); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusNotFound, stack.Wrap(ctx, err, errors.ErrInvalidID)
		}
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return &entity.Genre{
		ID:   genre.ID,
		Name: genre.Name,
	}, http.StatusOK, nil
}
