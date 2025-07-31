package mongo

import (
	"time"

	"github.com/rl404/hibiki/internal/domain/magazine/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type magazine struct {
	ID        int64     `bson:"id"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (m *magazine) MarshalBSON() ([]byte, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	m.UpdatedAt = time.Now()

	type m2 magazine
	return bson.Marshal((*m2)(m))
}

func (m *Mongo) fromEntity(data entity.Magazine) *magazine {
	return &magazine{
		ID:   data.ID,
		Name: data.Name,
	}
}
