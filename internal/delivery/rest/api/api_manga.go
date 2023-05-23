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
// @param mode query string false "mode" enums(ALL,SIMPLE) default(SIMPLE)
// @param title query string false "title"
// @param type query string false "type" enums(MANGA,NOVEL,ONE_SHOT,DOUJINSHI,MANHWA,MANHUA,OEL,LIGHT_NOVEL)
// @param start_date query string false "start date (yyyy-mm-dd)"
// @param end_date query string false "end date (yyyy-mm-dd)"
// @param sort query string false "sort" enums(title,-title,mean,-mean,rank,-rank,popularity,-popularity,member,-member,favorite,-favorite,start_date,-start_date) default(popularity)
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.manga,meta=service.pagination}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /manga [get]
func (api *API) handleGetManga(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	title := r.URL.Query().Get("title")
	_type := r.URL.Query().Get("type")
	startDate := r.URL.Query().Get("start_date")
	endDate := r.URL.Query().Get("end_date")
	sort := r.URL.Query().Get("sort")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	manga, pagination, code, err := api.service.GetManga(r.Context(), service.GetMangaRequest{
		Mode:      entity.SearchMode(mode),
		Type:      entity.Type(_type),
		Title:     title,
		StartDate: startDate,
		EndDate:   endDate,
		Sort:      sort,
		Page:      page,
		Limit:     limit,
	})

	utils.ResponseWithJSON(w, code, manga, errors.Wrap(r.Context(), err), pagination)
}

// @summary Get manga by id.
// @tags Manga
// @produce json
// @param mangaID path integer true "manga id"
// @success 200 {object} utils.Response{data=service.manga}
// @failure 202 {object} utils.Response
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /manga/{mangaID} [get]
func (api *API) handleGetMangaByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "mangaID"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidMangaID, err))
		return
	}

	manga, code, err := api.service.GetMangaByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, manga, errors.Wrap(r.Context(), err))
}
