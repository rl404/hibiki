package client

import (
	"context"
	"net/http"

	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/nagato"
)

// GetMangaByID to get manga by id.
func (c *Client) GetMangaByID(ctx context.Context, id int) (*nagato.Manga, int, error) {
	manga, code, err := c.client.GetMangaDetailsWithContext(ctx, id,
		nagato.MangaFieldAlternativeTitles,
		nagato.MangaFieldStartDate,
		nagato.MangaFieldEndDate,
		nagato.MangaFieldSynopsis,
		nagato.MangaFieldMean,
		nagato.MangaFieldRank,
		nagato.MangaFieldPopularity,
		nagato.MangaFieldNumListUsers,
		nagato.MangaFieldNumScoringUsers,
		nagato.MangaFieldNSFW,
		nagato.MangaFieldGenres,
		nagato.MangaFieldMediaType,
		nagato.MangaFieldStatus,
		nagato.MangaFieldNumVolumes,
		nagato.MangaFieldNumChapters,
		nagato.MangaFieldAuthors,
		nagato.MangaFieldPictures,
		nagato.MangaFieldBackground,
		nagato.MangaFieldSerialization,
		nagato.MangaFieldNumFavorites,
		nagato.MangaFieldRelatedManga(),
	)
	if err != nil {
		return nil, code, stack.Wrap(ctx, err)
	}

	return manga, http.StatusOK, nil
}
