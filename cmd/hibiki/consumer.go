package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/cache"
	_nr "github.com/rl404/fairy/log/newrelic"
	nrCache "github.com/rl404/fairy/monitoring/newrelic/cache"
	nrPS "github.com/rl404/fairy/monitoring/newrelic/pubsub"
	"github.com/rl404/fairy/pubsub"
	_consumer "github.com/rl404/hibiki/internal/delivery/consumer"
	authorRepository "github.com/rl404/hibiki/internal/domain/author/repository"
	authorCache "github.com/rl404/hibiki/internal/domain/author/repository/cache"
	authorMongo "github.com/rl404/hibiki/internal/domain/author/repository/mongo"
	emptyIDRepository "github.com/rl404/hibiki/internal/domain/empty_id/repository"
	emptyIDCache "github.com/rl404/hibiki/internal/domain/empty_id/repository/cache"
	emptyIDMongo "github.com/rl404/hibiki/internal/domain/empty_id/repository/mongo"
	genreRepository "github.com/rl404/hibiki/internal/domain/genre/repository"
	genreCache "github.com/rl404/hibiki/internal/domain/genre/repository/cache"
	genreMongo "github.com/rl404/hibiki/internal/domain/genre/repository/mongo"
	magazineRepository "github.com/rl404/hibiki/internal/domain/magazine/repository"
	magazineCache "github.com/rl404/hibiki/internal/domain/magazine/repository/cache"
	magazineMongo "github.com/rl404/hibiki/internal/domain/magazine/repository/mongo"
	mangaRepository "github.com/rl404/hibiki/internal/domain/manga/repository"
	mangaCache "github.com/rl404/hibiki/internal/domain/manga/repository/cache"
	mangaMongo "github.com/rl404/hibiki/internal/domain/manga/repository/mongo"
	mangaStatsHistoryRepository "github.com/rl404/hibiki/internal/domain/manga_stats_history/repository"
	mangaStatsHistoryMongo "github.com/rl404/hibiki/internal/domain/manga_stats_history/repository/mongo"
	nagatoRepository "github.com/rl404/hibiki/internal/domain/nagato/repository"
	nagatoClient "github.com/rl404/hibiki/internal/domain/nagato/repository/client"
	publisherRepository "github.com/rl404/hibiki/internal/domain/publisher/repository"
	publisherPubsub "github.com/rl404/hibiki/internal/domain/publisher/repository/pubsub"
	userMangaRepository "github.com/rl404/hibiki/internal/domain/user_manga/repository"
	userMangaCache "github.com/rl404/hibiki/internal/domain/user_manga/repository/cache"
	userMangaMongo "github.com/rl404/hibiki/internal/domain/user_manga/repository/mongo"
	"github.com/rl404/hibiki/internal/service"
	"github.com/rl404/hibiki/internal/utils"
)

func consumer() error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}
	utils.Info("config initialized")

	// Init newrelic.
	nrApp, err := newrelic.NewApplication(
		newrelic.ConfigAppName(cfg.Newrelic.Name),
		newrelic.ConfigLicense(cfg.Newrelic.LicenseKey),
		newrelic.ConfigDistributedTracerEnabled(true),
		newrelic.ConfigAppLogForwardingEnabled(true),
	)
	if err != nil {
		utils.Error(err.Error())
	} else {
		defer nrApp.Shutdown(10 * time.Second)
		utils.AddLog(_nr.NewFromNewrelicApp(nrApp, _nr.ErrorLevel))
		utils.Info("newrelic initialized")
	}

	// Init cache.
	c, err := cache.New(cacheType[cfg.Cache.Dialect], cfg.Cache.Address, cfg.Cache.Password, cfg.Cache.Time)
	if err != nil {
		return err
	}
	c = nrCache.New(cfg.Cache.Dialect, cfg.Cache.Address, c)
	utils.Info("cache initialized")
	defer c.Close()

	// Init db.
	db, err := newDB(cfg.DB)
	if err != nil {
		return err
	}
	utils.Info("database initialized")
	defer db.Client().Disconnect(context.Background())

	// Init pubsub.
	ps, err := pubsub.New(pubsubType[cfg.PubSub.Dialect], cfg.PubSub.Address, cfg.PubSub.Password)
	if err != nil {
		return err
	}
	ps = nrPS.New(cfg.PubSub.Dialect, ps)
	utils.Info("pubsub initialized")
	defer ps.Close()

	// Init manga.
	var manga mangaRepository.Repository
	manga = mangaMongo.New(db, cfg.Cron.FinishedAge, cfg.Cron.ReleasingAge, cfg.Cron.NotYetAge)
	manga = mangaCache.New(c, manga)
	utils.Info("repository manga initialized")

	// Init genre.
	var genre genreRepository.Repository
	genre = genreMongo.New(db)
	genre = genreCache.New(c, genre)
	utils.Info("repository genre initialized")

	// Init author.
	var author authorRepository.Repository
	author = authorMongo.New(db)
	author = authorCache.New(c, author)
	utils.Info("repository author initialized")

	// Init magazine.
	var magazine magazineRepository.Repository
	magazine = magazineMongo.New(db)
	magazine = magazineCache.New(c, magazine)
	utils.Info("repository magazine initialized")

	// Init user manga.
	var userManga userMangaRepository.Repository
	userManga = userMangaMongo.New(db, cfg.Cron.UserMangaAge)
	userManga = userMangaCache.New(c, userManga)
	utils.Info("repository user manga initialized")

	// Init manga stats history.
	var mangaStatsHistory mangaStatsHistoryRepository.Repository = mangaStatsHistoryMongo.New(db)
	utils.Info("repository manga stats history initialized")

	// Init empty id.
	var emptyID emptyIDRepository.Repository
	emptyID = emptyIDMongo.New(db)
	emptyID = emptyIDCache.New(c, emptyID)
	utils.Info("repository manga initialized")

	// Init publisher.
	var publisher publisherRepository.Repository = publisherPubsub.New(ps, pubsubTopic)
	utils.Info("repository publisher initialized")

	// Init nagato.
	var nagato nagatoRepository.Repository = nagatoClient.New(cfg.Mal.ClientID, cfg.Mal.ClientSecret)
	utils.Info("repository nagato initialized")

	// Init service.
	service := service.New(manga, genre, author, magazine, userManga, mangaStatsHistory, emptyID, publisher, nagato)
	utils.Info("service initialized")

	// Init consumer.
	consumer, err := _consumer.New(service, ps, pubsubTopic)
	if err != nil {
		return err
	}
	utils.Info("consumer initialized")
	defer consumer.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Start subscribe.
	if err := consumer.Subscribe(nrApp); err != nil {
		return err
	}

	utils.Info("consumer ready")
	<-sigChan

	return nil
}
