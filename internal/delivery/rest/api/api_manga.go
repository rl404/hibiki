package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/utils"
)

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
func (api *API) handleGetMangaByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "mangaID"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidMangaID, err))
		return
	}

	manga, code, err := api.service.GetMangaByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, manga, errors.Wrap(r.Context(), err))
}
