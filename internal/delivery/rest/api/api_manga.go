package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rl404/hibiki/internal/domain/manga/entity"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/service"
	"github.com/rl404/hibiki/internal/utils"
)

// @summary Get manga list.
// @tags Manga
// @produce json
// @param mode query string false "mode" enums(all, simple) default(all)
// @param title query string false "title"
// @param sort query string false "sort" enums(title,-title) default(title)
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.Manga}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /manga [get]
func (api *API) handleGetMangaByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "mangaID"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidMangaID, err))
		return
	}

	manga, code, err := api.service.GetMangaByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, manga, errors.Wrap(r.Context(), err))
}

// @summary Get manga by id.
// @tags Manga
// @produce json
// @param mangaID path integer true "manga id"
// @success 200 {object} utils.Response{data=service.Manga}
// @failure 202 {object} utils.Response
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /manga/{mangaID} [get]
func (api *API) handleGetManga(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	title := r.URL.Query().Get("title")
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	vtubers, pagination, code, err := api.service.GetManga(r.Context(), service.GetMangaRequest{
		Mode:  entity.SearchMode(mode),
		Title: title,
		Sort:  sort,
		Page:  page,
		Limit: limit,
	})

	utils.ResponseWithJSON(w, code, vtubers, errors.Wrap(r.Context(), err), pagination)
}
