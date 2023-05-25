package api

import (
	"strconv"

	"github.com/go-chi/chi"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rl404/fairy/log"
	"github.com/rl404/fairy/monitoring/newrelic/middleware"
	"github.com/rl404/hibiki/internal/service"
	"github.com/rl404/hibiki/internal/utils"
)

// API contains all functions for api endpoints.
type API struct {
	service service.Service
}

// New to create new api endpoints.
func New(service service.Service) *API {
	return &API{
		service: service,
	}
}

// Register to register api routes.
func (api *API) Register(r chi.Router, nrApp *newrelic.Application) {
	r.Route("/", func(r chi.Router) {
		r.Use(middleware.NewHTTP(nrApp))
		r.Use(log.MiddlewareWithLog(utils.GetLogger(0), log.MiddlewareConfig{Error: true}))
		r.Use(log.MiddlewareWithLog(utils.GetLogger(1), log.MiddlewareConfig{
			RequestHeader:  true,
			RequestBody:    true,
			ResponseHeader: true,
			ResponseBody:   true,
			RawPath:        true,
			Error:          true,
		}))
		r.Use(utils.Recoverer)

		r.Get("/manga", api.handleGetManga)
		r.Get("/manga/{mangaID}", api.handleGetMangaByID)

		r.Get("/authors", api.handleGetAuthors)

		r.Get("/genres", api.handleGetGenres)

		r.Get("/magazines", api.handleGetMagazines)

		r.Get("/user/{username}/manga", api.handleGetUserManga)
	})
}

func (api *API) parseBool(str string) *bool {
	b, err := strconv.ParseBool(str)
	if err != nil {
		return nil
	}
	return &b
}
