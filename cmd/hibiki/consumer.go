package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/newrelic/go-agent/v3/newrelic"
	_nr "github.com/rl404/fairy/log/newrelic"
	nrPS "github.com/rl404/fairy/monitoring/newrelic/pubsub"
	_consumer "github.com/rl404/hibiki/internal/delivery/consumer"
	authorRepository "github.com/rl404/hibiki/internal/domain/author/repository"
	authorMongo "github.com/rl404/hibiki/internal/domain/author/repository/mongo"
	emptyIDRepository "github.com/rl404/hibiki/internal/domain/empty_id/repository"
	emptyIDMongo "github.com/rl404/hibiki/internal/domain/empty_id/repository/mongo"
	genreRepository "github.com/rl404/hibiki/internal/domain/genre/repository"
	genreMongo "github.com/rl404/hibiki/internal/domain/genre/repository/mongo"
	magazineRepository "github.com/rl404/hibiki/internal/domain/magazine/repository"
	magazineMongo "github.com/rl404/hibiki/internal/domain/magazine/repository/mongo"
	mangaRepository "github.com/rl404/hibiki/internal/domain/manga/repository"
	mangaMongo "github.com/rl404/hibiki/internal/domain/manga/repository/mongo"
	mangaStatsHistoryRepository "github.com/rl404/hibiki/internal/domain/manga_stats_history/repository"
	mangaStatsHistoryMongo "github.com/rl404/hibiki/internal/domain/manga_stats_history/repository/mongo"
	nagatoRepository "github.com/rl404/hibiki/internal/domain/nagato/repository"
	nagatoClient "github.com/rl404/hibiki/internal/domain/nagato/repository/client"
	publisherRepository "github.com/rl404/hibiki/internal/domain/publisher/repository"
	publisherPubsub "github.com/rl404/hibiki/internal/domain/publisher/repository/pubsub"
	userMangaRepository "github.com/rl404/hibiki/internal/domain/user_manga/repository"
	userMangaMongo "github.com/rl404/hibiki/internal/domain/user_manga/repository/mongo"
	"github.com/rl404/hibiki/internal/service"
	"github.com/rl404/hibiki/internal/utils"
	"github.com/rl404/hibiki/pkg/pubsub"
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
	ps = nrPS.New(cfg.PubSub.Dialect, ps, nrApp)
	utils.Info("pubsub initialized")
	defer ps.Close()

	// Init manga.
	var manga mangaRepository.Repository = mangaMongo.New(db, cfg.Cron.FinishedAge, cfg.Cron.ReleasingAge, cfg.Cron.NotYetAge)
	utils.Info("repository manga initialized")

	// Init genre.
	var genre genreRepository.Repository = genreMongo.New(db)
	utils.Info("repository genre initialized")

	// Init author.
	var author authorRepository.Repository = authorMongo.New(db)
	utils.Info("repository author initialized")

	// Init magazine.
	var magazine magazineRepository.Repository = magazineMongo.New(db)
	utils.Info("repository magazine initialized")

	// Init user manga.
	var userManga userMangaRepository.Repository = userMangaMongo.New(db, cfg.Cron.UserMangaAge)
	utils.Info("repository user manga initialized")

	// Init manga stats history.
	var mangaStatsHistory mangaStatsHistoryRepository.Repository = mangaStatsHistoryMongo.New(db)
	utils.Info("repository manga stats history initialized")

	// Init empty id.
	var emptyID emptyIDRepository.Repository = emptyIDMongo.New(db)
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
	consumer := _consumer.New(service, ps, pubsubTopic)
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
