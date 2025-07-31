package mongo

import (
	"context"
	_errors "errors"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/domain/author/entity"
	"github.com/rl404/hibiki/internal/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	defer cursor.Close(ctx)

	existMap := make(map[int64]bool)
	for cursor.Next(ctx) {
		var author author
		if err := cursor.Decode(&author); err != nil {
			return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
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
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return http.StatusOK, nil
}

// GetAll to get all.
func (m *Mongo) GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.Author, int, int, error) {
	newFieldStage := bson.D{{Key: "$addFields", Value: bson.M{
		"name": bson.M{
			"$cond": bson.A{
				bson.M{"$eq": bson.A{"$first_name", ""}},
				"$last_name",
				bson.M{
					"$cond": bson.A{
						bson.M{"$eq": bson.A{"$last_name", ""}},
						"$first_name",
						bson.M{
							"$concat": bson.A{"$first_name", " ", "$last_name"},
						},
					},
				},
			},
		},
	}}}
	matchStage := bson.D{}
	sortStage := bson.D{{Key: "$sort", Value: bson.M{"name": 1}}}
	skipStage := bson.D{{Key: "$skip", Value: (data.Page - 1) * data.Limit}}
	limitStage := bson.D{}
	countStage := bson.D{{Key: "$count", Value: "count"}}

	if data.Name != "" {
		matchStage = m.addMatch(matchStage, "$or", []bson.M{
			{"first_name": bson.M{"$regex": data.Name, "$options": "i"}},
			{"last_name": bson.M{"$regex": data.Name, "$options": "i"}},
		})
	}

	if data.Limit > 0 {
		limitStage = append(limitStage, bson.E{Key: "$limit", Value: data.Limit})
	}

	cursor, err := m.db.Aggregate(ctx, m.getPipeline(newFieldStage, matchStage, sortStage, skipStage, limitStage))
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	var authors []author
	if err := cursor.All(ctx, &authors); err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	res := make([]entity.Author, len(authors))
	for i, author := range authors {
		res[i] = author.toEntity()
	}

	cntCursor, err := m.db.Aggregate(ctx, m.getPipeline(newFieldStage, matchStage, countStage))
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	var total []map[string]int64
	if err := cntCursor.All(ctx, &total); err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	if len(total) == 0 {
		return res, 0, http.StatusOK, nil
	}

	return res, int(total[0]["count"]), http.StatusOK, nil
}

// GetByID to get by id.
func (m *Mongo) GetByID(ctx context.Context, id int64) (*entity.Author, int, error) {
	var author author
	if err := m.db.FindOne(ctx, bson.M{"id": id}).Decode(&author); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusNotFound, stack.Wrap(ctx, err, errors.ErrInvalidID)
		}
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return &entity.Author{
		ID:        author.ID,
		FirstName: author.FirstName,
		LastName:  author.LastName,
	}, http.StatusOK, nil
}
