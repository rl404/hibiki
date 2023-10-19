package service

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	authorEntity "github.com/rl404/hibiki/internal/domain/author/entity"
	genreEntity "github.com/rl404/hibiki/internal/domain/genre/entity"
	magazineEntity "github.com/rl404/hibiki/internal/domain/magazine/entity"
	mangaEntity "github.com/rl404/hibiki/internal/domain/manga/entity"
	historyEntity "github.com/rl404/hibiki/internal/domain/manga_stats_history/entity"
)

func (s *service) updateManga(ctx context.Context, id int64) (int, error) {
	// Call mal api.
	manga, code, err := s.nagato.GetMangaByID(ctx, int(id))
	if err != nil {
		if code == http.StatusNotFound {
			// Insert empty id.
			if code, err := s.emptyID.Create(ctx, id); err != nil {
				return code, stack.Wrap(ctx, err)
			}

			// Delete existing data.
			if code, err := s.manga.DeleteByID(ctx, id); err != nil {
				return code, stack.Wrap(ctx, err)
			}

			if code, err := s.userManga.DeleteByMangaID(ctx, id); err != nil {
				return code, stack.Wrap(ctx, err)
			}
		}
		return code, stack.Wrap(ctx, err)
	}

	// Update genre data.
	if len(manga.Genres) > 0 {
		genres := make([]genreEntity.Genre, len(manga.Genres))
		for i, g := range manga.Genres {
			genres[i] = genreEntity.Genre{
				ID:   int64(g.ID),
				Name: g.Name,
			}
		}

		if code, err := s.genre.BatchUpdate(ctx, genres); err != nil {
			return code, stack.Wrap(ctx, err)
		}
	}

	// Update author data.
	if len(manga.Authors) > 0 {
		authors := make([]authorEntity.Author, len(manga.Authors))
		for i, a := range manga.Authors {
			authors[i] = authorEntity.Author{
				ID:        int64(a.Person.ID),
				FirstName: a.Person.FirstName,
				LastName:  a.Person.LastName,
			}
		}

		if code, err := s.author.BatchUpdate(ctx, authors); err != nil {
			return code, stack.Wrap(ctx, err)
		}
	}

	// Update serialization data.
	if len(manga.Serialization) > 0 {
		magazines := make([]magazineEntity.Magazine, len(manga.Serialization))
		for i, m := range manga.Serialization {
			magazines[i] = magazineEntity.Magazine{
				ID:   int64(m.Magazine.ID),
				Name: m.Magazine.Name,
			}
		}

		if code, err := s.magazine.BatchUpdate(ctx, magazines); err != nil {
			return code, stack.Wrap(ctx, err)
		}
	}

	// Update manga data.
	mangaE, err := mangaEntity.MangaFromNagato(ctx, manga)
	if err != nil {
		return http.StatusInternalServerError, stack.Wrap(ctx, err)
	}

	if code, err := s.manga.Update(ctx, *mangaE); err != nil {
		return code, stack.Wrap(ctx, err)
	}

	// Insert manga stats history.
	if code, err := s.mangaStatsHistory.Create(ctx, historyEntity.MangaStatsHistory{
		MangaID:    mangaE.ID,
		Mean:       mangaE.Mean,
		Rank:       mangaE.Rank,
		Popularity: mangaE.Popularity,
		Member:     mangaE.Member,
		Voter:      mangaE.Voter,
		Favorite:   mangaE.Favorite,
	}); err != nil {
		return code, stack.Wrap(ctx, err)
	}

	// Queue related manga.
	for _, r := range manga.RelatedManga {
		if err := s.publisher.PublishParseManga(ctx, int64(r.Manga.ID)); err != nil {
			return http.StatusInternalServerError, stack.Wrap(ctx, err)
		}
	}

	return http.StatusOK, nil
}
