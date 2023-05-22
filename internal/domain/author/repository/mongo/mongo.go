package mongo

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/hibiki/internal/domain/author/entity"
	"github.com/rl404/hibiki/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo contains functions for author mongodb.
type Mongo struct {
	db *mongo.Collection
}

// New to create new author database.
func New(mdb *mongo.Database) *Mongo {
	return &Mongo{
		db: mdb.Collection("author"),
	}
}

// BatchUpdate to batch update.
func (m *Mongo) BatchUpdate(ctx context.Context, data []entity.Author) (int, error) {
	ids := make([]int64, len(data))
	for i, a := range data {
		ids[i] = a.ID
	}

	cursor, err := m.db.Find(ctx, bson.M{"id": bson.M{"$in": ids}})
	if err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	defer cursor.Close(ctx)

	existMap := make(map[int64]bool)
	for cursor.Next(ctx) {
		var author author
		if err := cursor.Decode(&author); err != nil {
			return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
		}

		existMap[author.ID] = true
	}

	models := make([]mongo.WriteModel, len(data))
	for i, d := range data {
		if _, ok := existMap[d.ID]; ok {
			models[i] = mongo.NewUpdateOneModel().
				SetFilter(bson.M{"id": d.ID}).
				SetUpdate(bson.M{"$set": bson.M{
					"first_name": d.FirstName,
					"last_name":  d.LastName,
					"updated_at": time.Now(),
				}})
		} else {
			models[i] = mongo.NewInsertOneModel().SetDocument(m.fromEntity(d))
		}
	}

	if _, err := m.db.BulkWrite(ctx, models); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return http.StatusOK, nil
}

// GetAll to get all.
func (m *Mongo) GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.Author, int, int, error) {
	filter := bson.M{}
	opt := options.Find().SetSort(bson.M{"first_name": 1}).SetSkip(int64((data.Page - 1) * data.Limit)).SetLimit(int64(data.Limit))

	if data.Name != "" {
		filter = bson.M{"$or": []bson.M{
			{"first_name": bson.M{"$regex": data.Name, "$options": "i"}},
			{"last_name": bson.M{"$regex": data.Name, "$options": "i"}},
		}}
	}

	if data.Limit < 0 {
		opt.SetLimit(0)
	}

	c, err := m.db.Find(ctx, filter, opt)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	defer c.Close(ctx)

	var authors []entity.Author
	for c.Next(ctx) {
		var author author
		if err := c.Decode(&author); err != nil {
			return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}
		authors = append(authors, author.toEntity())
	}

	total, err := m.db.CountDocuments(ctx, filter, options.Count())
	if err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return authors, int(total), http.StatusOK, nil
}
