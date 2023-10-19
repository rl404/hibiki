package service

import (
	"context"
	"net/http"
	"strings"

	"github.com/rl404/fairy/errors/stack"
	nagatoEntity "github.com/rl404/hibiki/internal/domain/nagato/entity"
	"github.com/rl404/hibiki/internal/domain/user_manga/entity"
)

func (s *service) updateUserManga(ctx context.Context, username string) (int, error) {
	username = strings.ToLower(username)

	var ids []int64
	var mangaList []entity.UserManga
	limit, offset := 500, 0
	for {
		// Call mal api.
		manga, code, err := s.nagato.GetUserManga(ctx, nagatoEntity.GetUserMangaRequest{
			Username: username,
			Limit:    limit + 1,
			Offset:   offset,
		})
		if err != nil {
			return code, stack.Wrap(ctx, err)
		}

		for _, m := range manga {
			ids = append(ids, int64(m.Manga.ID))
			mangaList = append(mangaList, entity.UserMangaFromNagato(username, m))

			// Queue related manga.
			if err := s.publisher.PublishParseManga(ctx, int64(m.Manga.ID)); err != nil {
				return http.StatusInternalServerError, stack.Wrap(ctx, err)
			}
		}

		if len(manga) <= limit || len(manga) == 0 {
			break
		}

		offset += limit
	}

	// Update.
	if code, err := s.userManga.BatchUpdate(ctx, mangaList); err != nil {
		return code, stack.Wrap(ctx, err)
	}

	// Delete manga not in list.
	if code, err := s.userManga.DeleteNotInList(ctx, username, ids); err != nil {
		return code, stack.Wrap(ctx, err)
	}

	return http.StatusOK, nil
}
