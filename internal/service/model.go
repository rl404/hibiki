package service

import "github.com/rl404/hibiki/internal/domain/manga/entity"

type pagination struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type alternativeTitles struct {
	Synonyms []string `json:"synonyms"`
	English  string   `json:"english"`
	Japanese string   `json:"japanese"`
}

type date struct {
	Year  int `json:"year"`
	Month int `json:"month"`
	Day   int `json:"day"`
}

type genre struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type related struct {
	ID       int64           `json:"id"`
	Title    string          `json:"title"`
	Relation entity.Relation `json:"relation"  swaggertype:"string"`
	Picture  string          `jspn:"picture"`
}

type author struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
}

type magazine struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (c *service) mangaFromEntity(mangaDB *entity.Manga) manga {
	genres := make([]genre, len(mangaDB.Genres))
	for i, g := range mangaDB.Genres {
		genres[i] = genre{
			ID:   g.ID,
			Name: g.Name,
		}
	}

	relateds := make([]related, len(mangaDB.Related))
	for i, r := range mangaDB.Related {
		relateds[i] = related{
			ID:       r.ID,
			Title:    r.Title,
			Relation: r.Relation,
			Picture:  r.Picture,
		}
	}

	authors := make([]author, len(mangaDB.Authors))
	for i, a := range mangaDB.Authors {
		authors[i] = author{
			ID:   a.ID,
			Name: a.Name,
			Role: a.Role,
		}
	}

	serialization := make([]magazine, len(mangaDB.Serialization))
	for i, s := range mangaDB.Serialization {
		serialization[i] = magazine{
			ID:   s.ID,
			Name: s.Name,
		}
	}

	return manga{
		ID:    mangaDB.ID,
		Title: mangaDB.Title,
		AlternativeTitles: alternativeTitles{
			Synonyms: mangaDB.AlternativeTitles.Synonyms,
			Japanese: mangaDB.AlternativeTitles.Japanese,
			English:  mangaDB.AlternativeTitles.English,
		},
		Picture: mangaDB.Picture,
		StartDate: date{
			Year:  mangaDB.StartDate.Year,
			Month: mangaDB.StartDate.Month,
			Day:   mangaDB.StartDate.Day,
		},
		EndDate: date{
			Year:  mangaDB.EndDate.Year,
			Month: mangaDB.EndDate.Month,
			Day:   mangaDB.EndDate.Day,
		},
		Synopsis:      mangaDB.Synopsis,
		Background:    mangaDB.Background,
		NSFW:          mangaDB.NSFW,
		Type:          mangaDB.Type,
		Status:        mangaDB.Status,
		Chapter:       mangaDB.Chapter,
		Volume:        mangaDB.Volume,
		Mean:          mangaDB.Mean,
		Rank:          mangaDB.Rank,
		Popularity:    mangaDB.Popularity,
		Member:        mangaDB.Member,
		Voter:         mangaDB.Voter,
		Genres:        genres,
		Pictures:      mangaDB.Pictures,
		Related:       relateds,
		Authors:       authors,
		Serialization: serialization,
		UpdatedAt:     mangaDB.UpdatedAt,
	}
}
