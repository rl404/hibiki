package mongo

import (
	"time"

	"github.com/rl404/hibiki/internal/domain/manga/entity"
	"go.mongodb.org/mongo-driver/bson"
)

type manga struct {
	ID                int64            `bson:"id"`
	Title             string           `bson:"title"`
	AlternativeTitles alternativeTitle `bson:"alternative_titles"`
	Picture           string           `bson:"picture"`
	StartDate         date             `bson:"start_date"`
	EndDate           date             `bson:"end_date"`
	Synopsis          string           `bson:"synopsis"`
	Background        string           `bson:"background"`
	NSFW              bool             `bson:"nsfw"`
	Type              entity.Type      `bson:"type"`
	Status            entity.Status    `bson:"status"`
	Chapter           int              `bson:"chapter"`
	Volume            int              `bson:"volume"`
	Mean              float64          `bson:"mean"`
	Rank              int              `bson:"rank"`
	Popularity        int              `bson:"popularity"`
	Member            int              `bson:"member"`
	Voter             int              `bson:"voter"`
	Genres            []genre          `bson:"genres"`
	Pictures          []string         `bson:"pictures"`
	Related           []related        `bson:"related"`
	Authors           []author         `bson:"authors"`
	Serialization     []magazine       `bson:"serialization"`
	CreatedAt         time.Time        `bson:"created_at"`
	UpdatedAt         time.Time        `bson:"updated_at"`
}

func (m *manga) MarshalBSON() ([]byte, error) {
	if m.CreatedAt.IsZero() {
		m.CreatedAt = time.Now()
	}

	m.UpdatedAt = time.Now()

	type m2 manga
	return bson.Marshal((*m2)(m))
}

type alternativeTitle struct {
	Synonyms []string `bson:"synonyms"`
	English  string   `bson:"english"`
	Japanese string   `bson:"japanese"`
}

type date struct {
	Day   int `bson:"day"`
	Month int `bson:"month"`
	Year  int `bson:"year"`
}

type genre struct {
	ID   int64  `bson:"id"`
	Name string `bson:"name"`
}

type related struct {
	ID       int64           `bson:"id"`
	Title    string          `bson:"title"`
	Relation entity.Relation `bson:"relation"`
	Picture  string          `bson:"picture"`
}

type author struct {
	ID   int64  `bson:"id"`
	Name string `bson:"name"`
	Role string `bson:"role"`
}

type magazine struct {
	ID   int64  `bson:"id"`
	Name string `bson:"name"`
}

func (m *manga) toEntity() *entity.Manga {
	genres := make([]entity.Genre, len(m.Genres))
	for i, g := range m.Genres {
		genres[i] = entity.Genre{
			ID:   g.ID,
			Name: g.Name,
		}
	}

	related := make([]entity.Related, len(m.Related))
	for i, r := range m.Related {
		related[i] = entity.Related{
			ID:       r.ID,
			Title:    r.Title,
			Relation: r.Relation,
			Picture:  r.Picture,
		}
	}

	authors := make([]entity.Author, len(m.Authors))
	for i, a := range m.Authors {
		authors[i] = entity.Author{
			ID:   a.ID,
			Name: a.Name,
			Role: a.Role,
		}
	}

	serialization := make([]entity.Magazine, len(m.Serialization))
	for i, s := range m.Serialization {
		serialization[i] = entity.Magazine{
			ID:   s.ID,
			Name: s.Name,
		}
	}

	return &entity.Manga{
		ID:    m.ID,
		Title: m.Title,
		AlternativeTitles: entity.AlternativeTitle{
			Synonyms: m.AlternativeTitles.Synonyms,
			English:  m.AlternativeTitles.English,
			Japanese: m.AlternativeTitles.Japanese,
		},
		Picture: m.Picture,
		StartDate: entity.Date{
			Year:  m.StartDate.Year,
			Month: m.StartDate.Month,
			Day:   m.StartDate.Day,
		},
		EndDate: entity.Date{
			Year:  m.EndDate.Year,
			Month: m.EndDate.Month,
			Day:   m.EndDate.Day,
		},
		Synopsis:      m.Synopsis,
		Background:    m.Background,
		NSFW:          m.NSFW,
		Type:          m.Type,
		Status:        m.Status,
		Chapter:       m.Chapter,
		Volume:        m.Volume,
		Mean:          m.Mean,
		Rank:          m.Rank,
		Popularity:    m.Popularity,
		Member:        m.Member,
		Voter:         m.Voter,
		Genres:        genres,
		Pictures:      m.Pictures,
		Related:       related,
		Authors:       authors,
		Serialization: serialization,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (m *Mongo) mangaFromEntity(ma entity.Manga) *manga {
	genres := make([]genre, len(ma.Genres))
	for i, g := range ma.Genres {
		genres[i] = genre{
			ID:   g.ID,
			Name: g.Name,
		}
	}

	relateds := make([]related, len(ma.Related))
	for i, r := range ma.Related {
		relateds[i] = related{
			ID:       r.ID,
			Title:    r.Title,
			Relation: r.Relation,
		}
	}

	authors := make([]author, len(ma.Authors))
	for i, a := range ma.Authors {
		authors[i] = author{
			ID:   a.ID,
			Name: a.Name,
			Role: a.Role,
		}
	}

	serialization := make([]magazine, len(ma.Serialization))
	for i, s := range ma.Serialization {
		serialization[i] = magazine{
			ID:   s.ID,
			Name: s.Name,
		}
	}

	return &manga{
		ID:    ma.ID,
		Title: ma.Title,
		AlternativeTitles: alternativeTitle{
			Synonyms: ma.AlternativeTitles.Synonyms,
			Japanese: ma.AlternativeTitles.Japanese,
			English:  ma.AlternativeTitles.English,
		},
		Picture: ma.Picture,
		StartDate: date{
			Year:  ma.StartDate.Year,
			Month: ma.StartDate.Month,
			Day:   ma.StartDate.Day,
		},
		EndDate: date{
			Year:  ma.EndDate.Year,
			Month: ma.EndDate.Month,
			Day:   ma.EndDate.Day,
		},
		Synopsis:      ma.Synopsis,
		Background:    ma.Background,
		NSFW:          ma.NSFW,
		Type:          ma.Type,
		Status:        ma.Status,
		Chapter:       ma.Chapter,
		Volume:        ma.Volume,
		Mean:          ma.Mean,
		Rank:          ma.Rank,
		Popularity:    ma.Popularity,
		Member:        ma.Member,
		Voter:         ma.Voter,
		Genres:        genres,
		Pictures:      ma.Pictures,
		Related:       relateds,
		Authors:       authors,
		Serialization: serialization,
		UpdatedAt:     ma.UpdatedAt,
	}
}
