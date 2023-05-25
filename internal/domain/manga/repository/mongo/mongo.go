package mongo

import (
	"context"
	_errors "errors"
	"net/http"
	"time"

	"github.com/rl404/hibiki/internal/domain/manga/entity"
	"github.com/rl404/hibiki/internal/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Mongo contains functions for manga mongodb.
type Mongo struct {
	db           *mongo.Collection
	finishedAge  time.Duration
	releasingAge time.Duration
	notYetAge    time.Duration
}

// New to create new manga database.
func New(db *mongo.Database, finishedAge, releasingAge, notYetAge int) *Mongo {
	return &Mongo{
		db:           db.Collection("manga"),
		finishedAge:  time.Duration(finishedAge) * 24 * time.Hour,
		releasingAge: time.Duration(releasingAge) * 24 * time.Hour,
		notYetAge:    time.Duration(notYetAge) * 24 * time.Hour,
	}
}

// GetByID to get manga by id.
func (m *Mongo) GetByID(ctx context.Context, id int64) (*entity.Manga, int, error) {
	var manga manga
	if err := m.db.FindOne(ctx, bson.M{"id": id}).Decode(&manga); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return nil, http.StatusNotFound, errors.Wrap(ctx, errors.ErrMangaNotFound, err)
		}
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return manga.toEntity(), http.StatusOK, nil
}

// Update to update manga.
func (m *Mongo) Update(ctx context.Context, data entity.Manga) (int, error) {
	var manga manga
	if err := m.db.FindOne(ctx, bson.M{"id": data.ID}).Decode(&manga); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			if _, err := m.db.InsertOne(ctx, m.mangaFromEntity(data)); err != nil {
				return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
			}
			return http.StatusOK, nil
		}
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	mm := m.mangaFromEntity(data)
	mm.CreatedAt = manga.CreatedAt

	if _, err := m.db.UpdateOne(ctx, bson.M{"id": data.ID}, bson.M{"$set": mm}); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return http.StatusOK, nil
}

// DeleteByID to delete manga by id.
func (m *Mongo) DeleteByID(ctx context.Context, id int64) (int, error) {
	if _, err := m.db.DeleteOne(ctx, bson.M{"id": id}); err != nil {
		return http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return http.StatusOK, nil
}

// IsOld to check if old.
func (m *Mongo) IsOld(ctx context.Context, id int64) (bool, int, error) {
	filter := bson.M{
		"id": id,
		"$or": bson.A{
			bson.M{"status": bson.M{"$in": []entity.Status{entity.StatusFinished, entity.StatusHiatus, entity.StatusDiscontinued}}, "updated_at": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().Add(-m.finishedAge))}},
			bson.M{"status": entity.StatusReleasing, "updated_at": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().Add(-m.releasingAge))}},
			bson.M{"status": entity.StatusNotYet, "updated_at": bson.M{"$gte": primitive.NewDateTimeFromTime(time.Now().Add(-m.notYetAge))}},
		},
	}

	if err := m.db.FindOne(ctx, filter).Decode(&manga{}); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return true, http.StatusNotFound, nil
		}
		return true, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	return false, http.StatusOK, nil
}

// GetMaxID to get max id.
func (m *Mongo) GetMaxID(ctx context.Context) (int64, int, error) {
	var manga manga
	if err := m.db.FindOne(ctx, bson.M{}, options.FindOne().SetSort(bson.M{"id": -1}).SetProjection(bson.M{"id": 1})).Decode(&manga); err != nil {
		if _errors.Is(err, mongo.ErrNoDocuments) {
			return 1, http.StatusOK, nil
		}
		return 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	return manga.ID, http.StatusOK, nil
}

// GetIDs to get ids.
func (m *Mongo) GetIDs(ctx context.Context) ([]int64, int, error) {
	cursor, err := m.db.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"id": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var ids []int64
	for cursor.Next(ctx) {
		var manga manga
		if err := cursor.Decode(&manga); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}

		ids = append(ids, manga.ID)
	}

	return ids, http.StatusOK, nil
}

func (m *Mongo) getOldIDs(ctx context.Context, statuses []entity.Status, age time.Duration) ([]int64, int, error) {
	cursor, err := m.db.Find(ctx, bson.M{
		"status":     bson.M{"$in": statuses},
		"updated_at": bson.M{"$lte": primitive.NewDateTimeFromTime(time.Now().Add(-age))},
	}, options.Find().SetProjection(bson.M{"id": 1}))
	if err != nil {
		return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}
	defer cursor.Close(ctx)

	var ids []int64
	for cursor.Next(ctx) {
		var manga manga
		if err := cursor.Decode(&manga); err != nil {
			return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
		}

		ids = append(ids, manga.ID)
	}

	return ids, http.StatusOK, nil
}

// GetOldFinishedIDs to get old finished manga id.
func (m *Mongo) GetOldFinishedIDs(ctx context.Context) ([]int64, int, error) {
	return m.getOldIDs(ctx, []entity.Status{entity.StatusFinished, entity.StatusHiatus, entity.StatusDiscontinued}, m.finishedAge)
}

// GetOldReleasingIDs to get old releasing manga id.
func (m *Mongo) GetOldReleasingIDs(ctx context.Context) ([]int64, int, error) {
	return m.getOldIDs(ctx, []entity.Status{entity.StatusReleasing}, m.releasingAge)
}

// GetOldNotYetIDs to get not yet released manga id.
func (m *Mongo) GetOldNotYetIDs(ctx context.Context) ([]int64, int, error) {
	return m.getOldIDs(ctx, []entity.Status{entity.StatusFinished}, m.notYetAge)
}

// GetAll to get all data.
func (m *Mongo) GetAll(ctx context.Context, data entity.GetAllRequest) ([]entity.Manga, int, int, error) {
	newFieldStage := bson.D{{Key: "$addFields", Value: bson.M{
		"start_date_2": bson.M{"$dateFromParts": bson.M{
			"year":  bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$start_date.year", 0}}, 1, "$start_date.year"}},
			"month": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$start_date.month", 0}}, 1, "$start_date.month"}},
			"day":   bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$start_date.day", 0}}, 1, "$start_date.day"}},
		}},
		"end_date_2": bson.M{"$dateFromParts": bson.M{
			"year":  bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$end_date.year", 0}}, 1, "$end_date.year"}},
			"month": bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$end_date.month", 0}}, 1, "$end_date.month"}},
			"day":   bson.M{"$cond": bson.A{bson.M{"$eq": bson.A{"$end_date.day", 0}}, 1, "$end_date.day"}},
		}}}}}
	matchStage := bson.D{}
	projectStage := bson.D{}
	sortStage := bson.D{{Key: "$sort", Value: m.convertSort(data.Sort)}}
	skipStage := bson.D{{Key: "$skip", Value: (data.Page - 1) * data.Limit}}
	limitStage := bson.D{}
	countStage := bson.D{{Key: "$count", Value: "count"}}

	if data.Mode == entity.SearchModeSimple {
		projectStage = bson.D{{Key: "$project", Value: bson.M{
			"id":      1,
			"title":   1,
			"picture": 1,
			"nsfw":    1,
			"mean":    1,
		}}}
	}

	if data.Title != "" {
		matchStage = m.addMatch(matchStage, "$or", []bson.M{
			{"title": bson.M{"$regex": data.Title, "$options": "i"}},
			{"alternative_titles.synonyms": bson.M{"$regex": data.Title, "$options": "i"}},
			{"alternative_titles.english": bson.M{"$regex": data.Title, "$options": "i"}},
			{"alternative_titles.japanese": bson.M{"$regex": data.Title, "$options": "i"}},
		})
	}

	if data.Type != "" {
		matchStage = m.addMatch(matchStage, "type", data.Type)
	}

	if data.StartDate != nil {
		matchStage = m.addMatch(matchStage, "start_date_2", bson.M{"$gte": data.StartDate})
	}

	if data.EndDate != nil {
		matchStage = m.addMatch(matchStage, "start_date_2", bson.M{"$lte": data.EndDate})
	}

	if data.NSFW != nil {
		matchStage = m.addMatch(matchStage, "nsfw", data.NSFW)
	}

	if data.Limit > 0 {
		limitStage = append(limitStage, bson.E{Key: "$limit", Value: data.Limit})
	}

	cursor, err := m.db.Aggregate(ctx, m.getPipeline(newFieldStage, matchStage, sortStage, projectStage, skipStage, limitStage))
	if err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var mangas []manga
	if err := cursor.All(ctx, &mangas); err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	res := make([]entity.Manga, len(mangas))
	for i, manga := range mangas {
		res[i] = *manga.toEntity()
	}

	cntCursor, err := m.db.Aggregate(ctx, m.getPipeline(newFieldStage, matchStage, countStage))
	if err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	var total []map[string]int64
	if err := cntCursor.All(ctx, &total); err != nil {
		return nil, 0, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalDB, err)
	}

	if len(total) == 0 {
		return res, 0, http.StatusOK, nil
	}

	return res, int(total[0]["count"]), http.StatusOK, nil
}
