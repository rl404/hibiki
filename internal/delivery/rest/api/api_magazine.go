package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/service"
	"github.com/rl404/hibiki/internal/utils"
)

// @summary Get magazine list.
// @tags Magazine
// @produce json
// @param name query string false "name"
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.magazine,meta=service.pagination}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /magazines [get]
func (api *API) handleGetMagazines(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	magazines, pagination, code, err := api.service.GetMagazines(r.Context(), service.GetMagazinesRequest{
		Name:  name,
		Page:  page,
		Limit: limit,
	})

	utils.ResponseWithJSON(w, code, magazines, errors.Wrap(r.Context(), err), pagination)
}

// @summary Get magazine by id.
// @tags Magazine
// @produce json
// @param magazineID path integer true "magazine id"
// @success 200 {object} utils.Response{data=service.magazine}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /magazine/{magazineID} [get]
func (api *API) handleGetMagazineByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "magazineID"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, errors.Wrap(r.Context(), errors.ErrInvalidID, err))
		return
	}

	magazine, code, err := api.service.GetMagazineByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, magazine, errors.Wrap(r.Context(), err))
}
