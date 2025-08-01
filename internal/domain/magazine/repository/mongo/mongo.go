package mongo

import (
	"context"
	_errors "errors"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/domain/magazine/entity"
	"github.com/rl404/hibiki/internal/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Mongo contains functions for magazine mongodb.
type Mongo struct {
	db *mongo.Collection
}

// New to create new magazine database.
func New(mdb *mongo.Database) *Mongo {
	return &Mongo{
		db: mdb.Collection("magazine"),
	}
}

// BatchUpdate to batch update.
func (m *Mongo) BatchUpdate(ctx context.Context, data []entity.Magazine) (int, error) {
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
		var magazine magazine
		if err := cursor.Decode(&magazine); err != nil {
			return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
		}

		existMap[magazine.ID] = true
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
func (m *Mongo) GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.Magazine, int, int, error) {
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

	var magazines []entity.Magazine
	for c.Next(ctx) {
		var magazine magazine
		if err := c.Decode(&magazine); err != nil {
			return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
		}
		magazines = append(magazines, entity.Magazine{
			ID:   magazine.ID,
			Name: magazine.Name,
		})
	}

	total, err := m.db.CountDocuments(ctx, filter, options.Count())
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return magazines, int(total), http.StatusOK, nil
}

// GetByID to get by id.
func (m *Mongo) GetByID(ctx context.Context, id int64) (*entity.Magazine, int, error) {
	var magazine magazine
	if err := m.db.FindOne(ctx, bson.M{"id": id}).Decode(&magazine); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusNotFound, stack.Wrap(ctx, err, errors.ErrInvalidID)
		}
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return &entity.Magazine{
		ID:   magazine.ID,
		Name: magazine.Name,
	}, http.StatusOK, nil
}
