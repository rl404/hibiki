package mongo

import (
	"time"

	"github.com/rl404/hibiki/internal/domain/genre/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type genre struct {
	ID        int64     `bson:"id"`
	Name      string    `bson:"name"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (g *genre) MarshalBSON() ([]byte, error) {
	if g.CreatedAt.IsZero() {
		g.CreatedAt = time.Now()
	}

	g.UpdatedAt = time.Now()

	type g2 genre
	return bson.Marshal((*g2)(g))
}

func (m *Mongo) fromEntity(data entity.Genre) *genre {
	return &genre{
		ID:   data.ID,
		Name: data.Name,
	}
}
