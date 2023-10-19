# Hibiki

[![Go Report Card](https://goreportcard.com/badge/github.com/rl404/hibiki)](https://goreportcard.com/report/github.com/rl404/hibiki)
![License: MIT](https://img.shields.io/github/license/rl404/hibiki)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/rl404/hibiki)
[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/rl404/hibiki)](https://hub.docker.com/r/rl404/hibiki)
[![publish & deploy](https://github.com/rl404/hibiki/actions/workflows/publish-deploy.yml/badge.svg)](https://github.com/rl404/hibiki/actions/workflows/publish-deploy.yml)

Hibiki is [MyAnimeList](https://myanimelist.net/) manga database dump and REST API.

Powered by my [nagato](https://github.com/rl404/nagato) library and [MyAnimeList API](https://myanimelist.net/apiconfig/references/api/v2) as reference.

## Features

- Save manga details
  - Manga data
  - Manga genres
  - Manga pictures
  - Manga relation (with other manga)
  - Manga authors
  - Manga serialization (magazines)
- Save manga stats history
- Save user manga list
- Handle empty manga id
- Auto update manga & user data (cron)
- Interchangeable cache
  - no cache
  - inmemory
  - [Redis](https://redis.io/)
- Interchangeable pubsub
  - [RabbitMQ](https://www.rabbitmq.com/)
  - [Redis](https://redis.io/)
  - [Google PubSub](https://cloud.google.com/pubsub)
- [Swagger](https://github.com/swaggo/swag)
- [Docker](https://www.docker.com/)
- [Newrelic](https://newrelic.com/) monitoring
  - HTTP
  - Cron
  - Database
  - Cache
  - Pubsub
  - External API

_More will be coming soon..._

## Requirement

- [Go](https://go.dev/)
- [MyAnimeList](https://myanimelist.net/) [client id](https://myanimelist.net/apiconfig)
- [MongoDB](https://www.mongodb.com/)
- PubSub ([RabbitMQ](https://www.rabbitmq.com/)/[Redis](https://redis.io/)/[Google PubSub](https://cloud.google.com/pubsub))
- (optional) Cache ([Redis](https://redis.io/))
- (optional) [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)
- (optional) [Newrelic](https://newrelic.com/) license key

## Installation

### Without [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)

1. Clone the repository.

```sh
git clone github.com/rl404/hibiki
```

2. Rename `.env.sample` to `.env` and modify the values according to your setup.
3. Run. You need at least 2 consoles/terminals.

```sh
# Run the API.
make

# Run the consumer.
make consumer
```

6. [localhost:45002](http://localhost:45002) is ready (port may varies depend on your `.env`).

#### Other commands

```sh
# Update old manga data.
make cron-update

# Fill missing manga data.
make cron-fill
```

### With [Docker](https://www.docker.com/) & [Docker Compose](https://docs.docker.com/compose/)

1. Clone the repository.

```sh
git clone github.com/rl404/hibiki
```

2. Rename `.env.sample` to `.env` and modify the values according to your setup.
3. Get docker image.

```sh
# Pull existing image.
docker pull rl404/hibiki

# Or build your own.
make docker-build
```

4. Run the container. You need at least 2 consoles/terminals.

```sh
# Run the API.
make docker-api

# Run the consumer.
make docker-consumer
```

5. [localhost:45002](http://localhost:45002) is ready (port may varies depend on your `.env`).

#### Other commands

```sh
# Update old manga data.
make docker-cron-update

# Fill missing manga data.
make docker-cron-fill

# Stop running containers.
make docker-stop
```

## Environment Variables

| Env                            |           Default           | Description                                                                                                |
| ------------------------------ | :-------------------------: | ---------------------------------------------------------------------------------------------------------- |
| `HIBIKI_APP_ENV`               |            `dev`            | Environment type (`dev`/`prod`).                                                                           |
| `HIBIKI_HTTP_PORT`             |           `45002`           | HTTP server port.                                                                                          |
| `HIBIKI_HTTP_READ_TIMEOUT`     |            `5s`             | HTTP read timeout.                                                                                         |
| `HIBIKI_HTTP_WRITE_TIMEOUT`    |            `5s`             | HTTP write timeout.                                                                                        |
| `HIBIKI_HTTP_GRACEFUL_TIMEOUT` |            `10s`            | HTTP gracefull timeout.                                                                                    |
| `HIBIKI_GRPC_PORT`             |           `46002`           | GRPC server port.                                                                                          |
| `HIBIKI_GRPC_TIMEOUT`          |            `10s`            | GRPC timeout.                                                                                              |
| `HIBIKI_CACHE_DIALECT`         |         `inmemory`          | Cache type (`nocache`/`redis`/`inmemory`)                                                                  |
| `HIBIKI_CACHE_ADDRESS`         |                             | Cache address.                                                                                             |
| `HIBIKI_CACHE_PASSWORD`        |                             | Cache password.                                                                                            |
| `HIBIKI_CACHE_TIME`            |            `24h`            | Cache time.                                                                                                |
| `HIBIKI_DB_DIALECT`            |          `mongodb`          | Database type.                                                                                             |
| `HIBIKI_DB_ADDRESS`            | `mongodb://localhost:27017` | Database address with port.                                                                                |
| `HIBIKI_DB_NAME`               |          `hibiki`           | Database name.                                                                                             |
| `HIBIKI_DB_USER`               |                             | Database username.                                                                                         |
| `HIBIKI_DB_PASSWORD`           |                             | Database password.                                                                                         |
| `HIBIKI_DB_MAX_CONN_OPEN`      |            `10`             | Max open database connection.                                                                              |
| `HIBIKI_DB_MAX_CONN_IDLE`      |            `10`             | Max idle database connection.                                                                              |
| `HIBIKI_DB_MAX_CONN_LIFETIME`  |            `1m`             | Max database connection lifetime.                                                                          |
| `HIBIKI_PUBSUB_DIALECT`        |         `rabbitmq`          | Pubsub type (`rabbitmq`/`redis`/`google`)                                                                  |
| `HIBIKI_PUBSUB_ADDRESS`        |                             | Pubsub address (if you are using `google`, this will be your google project id).                           |
| `HIBIKI_PUBSUB_PASSWORD`       |                             | Pubsub password (if you are using `google`, this will be the content of your google service account json). |
| `HIBIKI_MAL_CLIENT_ID`         |                             | MyAnimeList client id.                                                                                     |
| `HIBIKI_CRON_UPDATE_LIMIT`     |            `10`             | Manga count limit when updating old data.                                                                  |
| `HIBIKI_CRON_FILL_LIMIT`       |            `30`             | Manga count limit when filling missing manga data.                                                         |
| `HIBIKI_CRON_RELEASING_AGE`    |             `7`             | Age of old releasing/airing manga data (in days).                                                          |
| `HIBIKI_CRON_FINISHED_AGE`     |            `30`             | Age of old finished manga data (in days).                                                                  |
| `HIBIKI_CRON_NOT_YET_AGE`      |             `7`             | Age of old not yet released/aired manga (in days).                                                         |
| `HIBIKI_CRON_USER_ANIME_AGE`   |             `7`             | Age of old user manga list (in days).                                                                      |
| `HIBIKI_NEWRELIC_NAME`         |          `hibiki`           | Newrelic application name.                                                                                 |
| `HIBIKI_NEWRELIC_LICENSE_KEY`  |                             | Newrelic license key.                                                                                      |

## Trivia

[Hibiki](<https://en.wikipedia.org/wiki/Japanese_destroyer_Hibiki_(1932)>)'s name is taken from japanese destroyer with her sisters (Inazuma, Akatsuki, Ikazuchi). Also, [exists](https://en.kancollewiki.net/Hibiki) in Kantai Collection games and manga.

## Disclaimer

Hibiki is meant for educational purpose and personal usage only. Please use it responsibly according to MyAnimeList [API License and Developer Agreement](https://myanimelist.net/static/apiagreement.html).

All data belong to their respective copyrights owners, hibiki does not have any affiliation with content providers.

## License

MIT License

Copyright (c) 2022 Axel
