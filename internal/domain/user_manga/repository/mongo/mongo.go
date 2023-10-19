package mongo

import (
	"context"
	_errors "errors"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/domain/user_manga/entity"
	"github.com/rl404/hibiki/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo contains functions for user_manga mongodb.
type Mongo struct {
	db  *mongo.Collection
	age time.Duration
}

// New to create new user_manga database.
func New(mdb *mongo.Database, age int) *Mongo {
	return &Mongo{
		db:  mdb.Collection("user_manga"),
		age: time.Duration(age) * 24 * time.Hour,
	}
}

// Get to get user manga.
func (m *Mongo) Get(ctx context.Context, data entity.GetUserMangaRequest) ([]*entity.UserManga, int, int, error) {
	filter := bson.M{"username": data.Username}

	count, err := m.db.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	if count == 0 {
		return nil, 0, http.StatusOK, nil
	}

	c, err := m.db.Find(ctx, filter, options.Find().SetSort(bson.M{"title": 1}).SetLimit(int64(data.Limit)).SetSkip(int64((data.Page-1)*data.Limit)))
	if err != nil {
		return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	defer c.Close(ctx)

	var res []*entity.UserManga
	for c.Next(ctx) {
		var um userManga
		if err := c.Decode(&um); err != nil {
			return nil, 0, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
		}

		res = append(res, um.toEntity())
	}

	return res, int(count), http.StatusOK, nil
}

// DeleteByMangaID to delete by manga id.
func (m *Mongo) DeleteByMangaID(ctx context.Context, mangaID int64) (int, error) {
	if _, err := m.db.DeleteMany(ctx, bson.M{"manga_id": mangaID}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}

// IsOld to check if old.
func (m *Mongo) IsOld(ctx context.Context, username string) (bool, int, error) {
	filter := bson.M{
		"username":   username,
		"updated_at": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().Add(-m.age))},
	}

	if err := m.db.FindOne(ctx, filter).Decode(&userManga{}); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return true, http.StatusNotFound, nil
		}
		return true, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return false, http.StatusOK, nil
}

// GetOldUsernames to get old usernames.
func (m *Mongo) GetOldUsernames(ctx context.Context) ([]string, int, error) {
	res, err := m.db.Distinct(ctx, "username", bson.M{"updated_at": bson.M{"$lte": primitive.NewDateTimeFromTime(time.Now().Add(-m.age))}})
	if err != nil {
		return nil, http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	usernames := make([]string, len(res))
	for i, r := range res {
		usernames[i] = r.(string)
	}

	return usernames, http.StatusOK, nil
}

// DeleteNotInList to delete not in list.
func (m *Mongo) DeleteNotInList(ctx context.Context, username string, ids []int64) (int, error) {
	if _, err := m.db.DeleteMany(ctx, bson.M{"username": username, "manga_id": bson.M{"$nin": ids}}); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	return http.StatusOK, nil
}

// BatchUpdate to batch update.
func (m *Mongo) BatchUpdate(ctx context.Context, data []entity.UserManga) (int, error) {
	username, ids := data[0].Username, make([]int64, len(data))
	for i, um := range data {
		ids[i] = um.MangaID
	}

	cursor, err := m.db.Find(ctx, bson.M{"username": username, "manga_id": bson.M{"$in": ids}})
	if err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}
	defer cursor.Close(ctx)

	existMap := make(map[int64]*userManga)
	for cursor.Next(ctx) {
		var userManga userManga
		if err := cursor.Decode(&userManga); err != nil {
			return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalServer)
		}

		existMap[userManga.MangaID] = &userManga
	}

	models := make([]mongo.WriteModel, len(data))
	for i, d := range data {
		if ex, ok := existMap[d.MangaID]; ok {
			mm := m.fromEntity(d)
			mm.CreatedAt = ex.CreatedAt

			models[i] = mongo.NewUpdateOneModel().
				SetFilter(bson.M{"username": username, "manga_id": d.MangaID}).
				SetUpdate(bson.M{"$set": mm})
		} else {
			models[i] = mongo.NewInsertOneModel().SetDocument(m.fromEntity(d))
		}
	}

	if _, err := m.db.BulkWrite(ctx, models); err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err, errors.ErrInternalDB)
	}

	return http.StatusOK, nil
}
