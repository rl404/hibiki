package api

import (
	"net/http"
	"strconv"

	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/service"
	"github.com/rl404/hibiki/internal/utils"
)

// @summary Get genre list.
// @tags Genre
// @produce json
// @param name query string false "name"
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.genre,meta=service.pagination}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /genres [get]
func (api *API) handleGetGenres(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	genres, pagination, code, err := api.service.GetGenres(r.Context(), service.GetGenresRequest{
		Name:  name,
		Page:  page,
		Limit: limit,
	})

	utils.ResponseWithJSON(w, code, genres, errors.Wrap(r.Context(), err), pagination)
}
