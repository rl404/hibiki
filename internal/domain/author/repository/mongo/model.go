package mongo

import (
	"time"

	"github.com/rl404/hibiki/internal/domain/author/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func (m *Mongo) addStage(stageKey string, stages bson.D, key string, value interface{}) bson.D {
	for i, stage := range stages {
		if stage.Key != stageKey {
			continue
		}

		matchValue, ok := stage.Value.(bson.M)
		if !ok {
			continue
		}

		if matchValue[key] == nil {
			matchValue[key] = bson.M{}
		}

		if mValue, ok := value.(bson.M); ok {
			for k, v := range mValue {
				matchValue[key].(bson.M)[k] = v
			}
		} else {
			matchValue[key] = value
		}

		stages[i].Value = matchValue
		return stages
	}

	return append(stages, bson.E{
		Key:   stageKey,
		Value: bson.M{key: value},
	})
}

func (m *Mongo) addMatch(matchStage bson.D, key string, value interface{}) bson.D {
	return m.addStage("$match", matchStage, key, value)
}

func (m *Mongo) getPipeline(stages ...bson.D) mongo.Pipeline {
	var pipelines mongo.Pipeline
	for _, stage := range stages {
		if len(stage) > 0 {
			pipelines = append(pipelines, stage)
		}
	}
	return pipelines
}
