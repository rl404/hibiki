package service

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/hibiki/internal/domain/manga/entity"
	publisherEntity "github.com/rl404/hibiki/internal/domain/publisher/entity"
	"github.com/rl404/hibiki/internal/errors"
)

// Manga is manga model.
type Manga struct {
	ID                int64             `json:"id"`
	Title             string            `json:"title"`
	AlternativeTitles alternativeTitles `json:"alternative_titles"`
	Picture           string            `json:"picture"`
	StartDate         date              `json:"start_date"`
	EndDate           date              `json:"end_date"`
	Synopsis          string            `json:"synopsis"`
	Background        string            `json:"background"`
	NSFW              bool              `json:"nsfw"`
	Type              entity.Type       `json:"type" swaggertype:"string"`
	Status            entity.Status     `json:"status" swaggertype:"string"`
	Chapter           int               `json:"chapter"`
	Volume            int               `json:"volume"`
	Mean              float64           `json:"mean"`
	Rank              int               `json:"rank"`
	Popularity        int               `json:"popularity"`
	Member            int               `json:"member"`
	Voter             int               `json:"voter"`
	Genres            []genre           `json:"genres"`
	Pictures          []string          `json:"pictures"`
	Related           []related         `json:"related"`
	Authors           []author          `json:"authors"`
	Serialization     []magazine        `json:"serialization"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// GetMangaByID to get manga by id.
func (s *service) GetMangaByID(ctx context.Context, id int64) (*Manga, int, error) {
	if code, err := s.validateID(ctx, id); err != nil {
		return nil, code, errors.Wrap(ctx, err)
	}

	// Get manga from db.
	mangaDB, code, err := s.manga.GetByID(ctx, id)
	if err != nil {
		if code == http.StatusNotFound {
			// Queue to parse.
			if err := s.publisher.PublishParseManga(ctx, publisherEntity.ParseMangaRequest{ID: id}); err != nil {
				return nil, http.StatusInternalServerError, errors.Wrap(ctx, errors.ErrInternalServer, err)
			}
			return nil, http.StatusAccepted, nil
		}
		return nil, code, errors.Wrap(ctx, err)
	}

	manga := s.mangaFromEntity(mangaDB)

	return &manga, http.StatusOK, nil
}

func (s *service) validateID(ctx context.Context, id int64) (int, error) {
	if id <= 0 {
		return http.StatusBadRequest, errors.Wrap(ctx, errors.ErrInvalidMangaID)
	}

	if _, code, err := s.emptyID.Get(ctx, id); err != nil {
		if code == http.StatusNotFound {
			return http.StatusOK, nil
		}
		return code, errors.Wrap(ctx, err)
	}

	return http.StatusNotFound, errors.Wrap(ctx, errors.ErrMangaNotFound)
}
