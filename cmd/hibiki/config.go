package main

import (
	"context"
	"io/ioutil"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/newrelic/go-agent/v3/integrations/nrmongo"
	"github.com/rl404/hibiki/internal/utils"
	"github.com/rl404/hibiki/pkg/cache"
	"github.com/rl404/hibiki/pkg/log"
	"github.com/rl404/hibiki/pkg/pubsub"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	App      appConfig      `envconfig:"APP"`
	HTTP     httpConfig     `envconfig:"HTTP"`
	GRPC     grpcConfig     `envconfig:"GRPC"`
	Cache    cacheConfig    `envconfig:"CACHE"`
	DB       dbConfig       `envconfig:"DB"`
	PubSub   pubsubConfig   `envconfig:"PUBSUB"`
	Mal      malConfig      `envconfig:"MAL"`
	Cron     cronConfig     `envconfig:"CRON"`
	Log      logConfig      `envconfig:"LOG"`
	Newrelic newrelicConfig `envconfig:"NEWRELIC"`
}

type appConfig struct {
	Env string `envconfig:"ENV" validate:"required,oneof=dev prod" mod:"default=dev,no_space,lcase"`
}

type httpConfig struct {
	Port            string        `envconfig:"PORT" validate:"required" mod:"default=45002,no_space"`
	ReadTimeout     time.Duration `envconfig:"READ_TIMEOUT" validate:"required,gt=0" mod:"default=5s"`
	WriteTimeout    time.Duration `envconfig:"WRITE_TIMEOUT" validate:"required,gt=0" mod:"default=5s"`
	GracefulTimeout time.Duration `envconfig:"GRACEFUL_TIMEOUT" validate:"required,gt=0" mod:"default=10s"`
}

type grpcConfig struct {
	Port    string        `envconfig:"PORT" validate:"required" mod:"default=46002,no_space"`
	Timeout time.Duration `envconfig:"TIMEOUT" validate:"required,gt=0" mod:"default=10s"`
}

type cacheConfig struct {
	Dialect  string        `envconfig:"DIALECT" validate:"required,oneof=nocache redis inmemory memcache" mod:"default=inmemory,no_space,lcase"`
	Address  string        `envconfig:"ADDRESS"`
	Password string        `envconfig:"PASSWORD"`
	Time     time.Duration `envconfig:"TIME" default:"24h" validate:"required,gt=0"`
}

type dbConfig struct {
	Dialect         string        `envconfig:"DIALECT" validate:"required,oneof=mongodb" mod:"default=mongodb,no_space,lcase"`
	Address         string        `envconfig:"ADDRESS" validate:"required" mod:"default=mongodb://localhost:27017,no_space"`
	Name            string        `envconfig:"NAME" validate:"required" mod:"default=hibiki"`
	User            string        `envconfig:"USER" validate:"required" mod:"default=root"`
	Password        string        `envconfig:"PASSWORD"`
	MaxConnOpen     int           `envconfig:"MAX_CONN_OPEN" validate:"required,gt=0" mod:"default=10"`
	MaxConnIdle     int           `envconfig:"MAX_CONN_IDLE" validate:"required,gt=0" mod:"default=10"`
	MaxConnLifetime time.Duration `envconfig:"MAX_CONN_LIFETIME" validate:"required,gt=0" mod:"default=1m"`
}

type pubsubConfig struct {
	Dialect  string `envconfig:"DIALECT" validate:"required,oneof=nsq rabbitmq redis google" mod:"default=rabbitmq,no_space,lcase"`
	Address  string `envconfig:"ADDRESS" validate:"required"`
	Password string `envconfig:"PASSWORD"`
}

type malConfig struct {
	ClientID     string `envconfig:"CLIENT_ID" validate:"required" mod:"no_space"`
	ClientSecret string `envconfig:"CLIENT_SECRET" mod:"no_space"`
}

type cronConfig struct {
	UpdateLimit  int `envconfig:"UPDATE_LIMIT" validate:"required,gte=0" mod:"default=10"`
	FillLimit    int `envconfig:"FILL_LIMIT" validate:"required,gte=0" mod:"default=30"`
	ReleasingAge int `envconfig:"RELEASING_AGE" validate:"required,gt=0" mod:"default=7"`  // days
	FinishedAge  int `envconfig:"FINISHED_AGE" validate:"required,gt=0" mod:"default=30"`  // days
	NotYetAge    int `envconfig:"NOT_YET_AGE" validate:"required,gt=0" mod:"default=7"`    // days
	UserMangaAge int `envconfig:"USER_MANGA_AGE" validate:"required,gt=0" mod:"default=7"` // days
}

type logConfig struct {
	Level log.LogLevel `envconfig:"LEVEL" default:"-1"`
	JSON  bool         `envconfig:"JSON" default:"false"`
	Color bool         `envconfig:"COLOR" default:"true"`
}

type newrelicConfig struct {
	Name       string `envconfig:"NAME" default:"hibiki"`
	LicenseKey string `envconfig:"LICENSE_KEY"`
}

const envPath = "../../.env"
const envPrefix = "HIBIKI"
const pubsubTopic = "hibiki-pubsub"

var cacheType = map[string]cache.CacheType{
	"nocache":  cache.NOP,
	"redis":    cache.Redis,
	"inmemory": cache.InMemory,
}

var pubsubType = map[string]pubsub.PubsubType{
	"rabbitmq": pubsub.RabbitMQ,
	"redis":    pubsub.Redis,
	"google":   pubsub.Google,
}

func getConfig() (*config, error) {
	var cfg config

	// Load .env file.
	_ = godotenv.Load(envPath)

	// Convert env to struct.
	if err := envconfig.Process(envPrefix, &cfg); err != nil {
		return nil, err
	}

	// Override PORT env.
	if port := os.Getenv("PORT"); port != "" {
		cfg.HTTP.Port = port
	}

	// Handle google pubsub credential.
	if cfg.PubSub.Dialect == "google" {
		credFilename, err := generateGoogleServiceAccountJSON("gcp-service-account.json", cfg.PubSub.Password)
		if err != nil {
			return nil, err
		}
		cfg.PubSub.Password = credFilename
	}

	// Validate.
	if err := utils.Validate(&cfg); err != nil {
		return nil, err
	}

	// Init global log.
	if err := utils.InitLog(cfg.Log.Level, cfg.Log.JSON, cfg.Log.Color); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func newDB(cfg dbConfig) (*mongo.Database, error) {
	nrMongo := nrmongo.NewCommandMonitor(nil)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start connection.
	client, err := mongo.Connect(ctx, options.
		Client().
		ApplyURI(cfg.Address).
		SetAuth(options.Credential{
			Username: cfg.User,
			Password: cfg.Password,
		}).
		SetMonitor(nrMongo))
	if err != nil {
		return nil, err
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	// Ping test.
	if err = client.Ping(ctx2, nil); err != nil {
		return nil, err
	}

	return client.Database(cfg.Name), nil

}

func generateGoogleServiceAccountJSON(filename, value string) (string, error) {
	if err := ioutil.WriteFile(filename, []byte(value), 0644); err != nil {
		return "", err
	}
	return filename, nil
}
