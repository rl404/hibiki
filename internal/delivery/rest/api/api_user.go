package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/service"
	"github.com/rl404/hibiki/internal/utils"
)

// @summary Get user's manga.
// @tags User
// @produce json
// @param username path string true "username"
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.userManga,meta=service.pagination}
// @failure 202 {object} utils.Response
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /user/{username}/manga [get]
func (api *API) handleGetUserManga(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	manga, pagination, code, err := api.service.GetUserManga(r.Context(), service.GetUserMangaRequest{
		Username: username,
		Page:     page,
		Limit:    limit,
	})

	utils.ResponseWithJSON(w, code, manga, errors.Wrap(r.Context(), err), pagination)
}
