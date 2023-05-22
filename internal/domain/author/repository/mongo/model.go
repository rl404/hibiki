package mongo

import (
	"time"

	"github.com/rl404/hibiki/internal/domain/author/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type author struct {
	ID        int64     `bson:"id"`
	FirstName string    `bson:"first_name"`
	LastName  string    `bson:"last_name"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

func (a *author) MarshalBSON() ([]byte, error) {
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}

	a.UpdatedAt = time.Now()

	type a2 author
	return bson.Marshal((*a2)(a))
}

func (a *author) toEntity() entity.Author {
	return entity.Author{
		ID:        a.ID,
		FirstName: a.FirstName,
		LastName:  a.LastName,
	}
}

func (m *Mongo) fromEntity(data entity.Author) *author {
	return &author{
		ID:        data.ID,
		FirstName: data.FirstName,
		LastName:  data.LastName,
	}
}
