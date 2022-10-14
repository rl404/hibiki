package service

import (
	"context"

	authorRepository "github.com/rl404/hibiki/internal/domain/author/repository"
	emptyIDRepository "github.com/rl404/hibiki/internal/domain/empty_id/repository"
	genreRepository "github.com/rl404/hibiki/internal/domain/genre/repository"
	magazineRepository "github.com/rl404/hibiki/internal/domain/magazine/repository"
	mangaRepository "github.com/rl404/hibiki/internal/domain/manga/repository"
	mangaStatsHistoryRepository "github.com/rl404/hibiki/internal/domain/manga_stats_history/repository"
	nagatoRepository "github.com/rl404/hibiki/internal/domain/nagato/repository"
	"github.com/rl404/hibiki/internal/domain/publisher/entity"
	publisherRepository "github.com/rl404/hibiki/internal/domain/publisher/repository"
	userMangaRepository "github.com/rl404/hibiki/internal/domain/user_manga/repository"
)

// Service contains functions for service.
type Service interface {
	GetMangaByID(ctx context.Context, id int64) (*Manga, int, error)

	GetUserManga(ctx context.Context, data GetUserMangaRequest) ([]UserManga, *Pagination, int, error)

	ConsumeMessage(ctx context.Context, msg entity.Message) error

	QueueOldReleasingManga(ctx context.Context, limit int) (int, int, error)
	QueueOldFinishedManga(ctx context.Context, limit int) (int, int, error)
	QueueOldNotYetManga(ctx context.Context, limit int) (int, int, error)
	QueueMissingManga(ctx context.Context, limit int) (int, int, error)
	QueueOldUserManga(ctx context.Context, limit int) (int, int, error)
}

type service struct {
	manga             mangaRepository.Repository
	genre             genreRepository.Repository
	author            authorRepository.Repository
	magazine          magazineRepository.Repository
	userManga         userMangaRepository.Repository
	mangaStatsHistory mangaStatsHistoryRepository.Repository
	emptyID           emptyIDRepository.Repository
	publisher         publisherRepository.Repository
	nagato            nagatoRepository.Repository
}

// New to create new service.
func New(
	manga mangaRepository.Repository,
	genre genreRepository.Repository,
	author authorRepository.Repository,
	magazine magazineRepository.Repository,
	userManga userMangaRepository.Repository,
	mangaStatsHistory mangaStatsHistoryRepository.Repository,
	emptyID emptyIDRepository.Repository,
	publisher publisherRepository.Repository,
	nagato nagatoRepository.Repository,
) Service {
	return &service{
		manga:             manga,
		genre:             genre,
		author:            author,
		magazine:          magazine,
		userManga:         userManga,
		mangaStatsHistory: mangaStatsHistory,
		emptyID:           emptyID,
		publisher:         publisher,
		nagato:            nagato,
	}
}
