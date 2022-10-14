package mongo

import (
	"time"

	"github.com/rl404/hibiki/internal/domain/manga_stats_history/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type mangaStatsHistory struct {
	MangaID    int64     `bson:"manga_id"`
	Mean       float64   `bson:"mean"`
	Rank       int       `bson:"rank"`
	Popularity int       `bson:"popularity"`
	Member     int       `bson:"member"`
	Voter      int       `bson:"voter"`
	CreatedAt  time.Time `bson:"created_at"`
}

func (m *mangaStatsHistory) MarshalBSON() ([]byte, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	type m2 mangaStatsHistory
	return bson.Marshal((*m2)(m))
}

func (m *Mongo) fromEntity(data entity.MangaStatsHistory) *mangaStatsHistory {
	return &mangaStatsHistory{
		MangaID:    data.MangaID,
		Mean:       data.Mean,
		Rank:       data.Rank,
		Popularity: data.Popularity,
		Member:     data.Member,
		Voter:      data.Voter,
	}
}
