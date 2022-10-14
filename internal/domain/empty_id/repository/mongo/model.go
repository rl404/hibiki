package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type emptyID struct {
	MangaID   int64     `bson:"manga_id"`
	CreatedAt time.Time `bson:"created_at"`
}

func (e *emptyID) MarshalBSON() ([]byte, error) {
	if e.CreatedAt.IsZero() {
		e.CreatedAt = time.Now()
	}

	type e2 emptyID
	return bson.Marshal((*e2)(e))
}
