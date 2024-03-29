package service

import (
	"context"
	"net/http"
	"time"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/domain/manga/entity"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/utils"
)

type manga struct {
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
	Favorite          int               `json:"favorite"`
	Genres            []genre           `json:"genres"`
	Pictures          []string          `json:"pictures"`
	Related           []related         `json:"related"`
	Authors           []mangaAuthor     `json:"authors"`
	Serialization     []magazine        `json:"serialization"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

// GetMangaByID to get manga by id.
func (s *service) GetMangaByID(ctx context.Context, id int64) (*manga, int, error) {
	if code, err := s.validateID(ctx, id); err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	// Get manga from db.
	mangaDB, code, err := s.manga.GetByID(ctx, id)
	if err != nil {
		if code == http.StatusNotFound {
			// Queue to parse.
			if err := s.publisher.PublishParseManga(ctx, id); err != nil {
				return nil, http.StatusInternalServerError, stack.Wrap(ctx, err)
			}
			return nil, http.StatusAccepted, nil
		}
		return nil, code, stack.Wrap(ctx, err)
	}

	manga := s.mangaFromEntity(mangaDB)

	return &manga, http.StatusOK, nil
}

func (s *service) validateID(ctx context.Context, id int64) (int, error) {
	if id <= 0 {
		return http.StatusBadRequest, stack.Wrap(ctx, errors.ErrInvalidID)
	}

	if _, code, err := s.emptyID.Get(ctx, id); err != nil {
		if code == http.StatusNotFound {
			return http.StatusOK, nil
		}
		return code, stack.Wrap(ctx, err)
	}

	return http.StatusNotFound, stack.Wrap(ctx, errors.ErrMangaNotFound)
}

// GetMangaRequest is get manga request model.
type GetMangaRequest struct {
	Mode       entity.SearchMode `validate:"oneof=ALL SIMPLE" mod:"default=SIMPLE,trim,ucase"`
	Title      string            `validate:"omitempty,gte=3" mod:"trim,lcase"`
	Type       entity.Type       `validate:"omitempty,oneof=MANGA NOVEL ONE_SHOT DOUJINSHI MANHWA MANHUA OEL LIGHT_NOVEL" mod:"trim,ucase"`
	Status     entity.Status     `validate:"omitempty,oneof=FINISHED RELEASING NOT_YET HIATUS DISCONTINUED" mod:"trim,ucase"`
	StartDate  string            `validate:"omitempty,datetime=2006-01-02" mod:"trim"`
	EndDate    string            `validate:"omitempty,datetime=2006-01-02" mod:"trim"`
	AuthorID   int64             `validate:"omitempty,gt=0"`
	MagazineID int64             `validate:"omitempty,gt=0"`
	GenreID    int64             `validate:"omitempty,gt=0"`
	NSFW       *bool             ``
	Sort       string            `validate:"omitempty,oneof=title -title mean -mean rank -rank popularity -popularity member -member favorite -favorite start_date -start_date" mod:"default=popularity,trim,lcase"`
	Page       int               `validate:"required,gte=1" mod:"default=1"`
	Limit      int               `validate:"required,gte=-1" mod:"default=20"`
}

// GetManga to get manga list.
func (s *service) GetManga(ctx context.Context, data GetMangaRequest) ([]manga, *pagination, int, error) {
	if err := utils.Validate(&data); err != nil {
		return nil, nil, http.StatusBadRequest, stack.Wrap(ctx, err)
	}

	mangas, total, code, err := s.manga.GetAll(ctx, entity.GetAllRequest{
		Mode:       data.Mode,
		Title:      data.Title,
		Type:       data.Type,
		Status:     data.Status,
		StartDate:  utils.ParseToTimePtr("2006-01-02", data.StartDate),
		EndDate:    utils.ParseToTimePtr("2006-01-02", data.EndDate),
		AuthorID:   data.AuthorID,
		MagazineID: data.MagazineID,
		GenreID:    data.GenreID,
		NSFW:       data.NSFW,
		Sort:       data.Sort,
		Page:       data.Page,
		Limit:      data.Limit,
	})
	if err != nil {
		return nil, nil, code, stack.Wrap(ctx, err)
	}

	res := make([]manga, len(mangas))
	for i, m := range mangas {
		res[i] = s.mangaFromEntity(&m)
	}

	return res, &pagination{
		Page:  data.Page,
		Limit: data.Limit,
		Total: total,
	}, http.StatusOK, nil
}
