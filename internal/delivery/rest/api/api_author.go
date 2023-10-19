package api

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/rl404/fairy/errors/stack"
	"github.com/rl404/hibiki/internal/errors"
	"github.com/rl404/hibiki/internal/service"
	"github.com/rl404/hibiki/internal/utils"
)

// @summary Get author list.
// @tags Author
// @produce json
// @param name query string false "name"
// @param page query integer false "page" default(1)
// @param limit query integer false "limit" default(20)
// @success 200 {object} utils.Response{data=[]service.author,meta=service.pagination}
// @failure 400 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /authors [get]
func (api *API) handleGetAuthors(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	authors, pagination, code, err := api.service.GetAuthors(r.Context(), service.GetAuthorsRequest{
		Name:  name,
		Page:  page,
		Limit: limit,
	})

	utils.ResponseWithJSON(w, code, authors, stack.Wrap(r.Context(), err), pagination)
}

// @summary Get author by id.
// @tags Author
// @produce json
// @param authorID path integer true "author id"
// @success 200 {object} utils.Response{data=service.author}
// @failure 400 {object} utils.Response
// @failure 404 {object} utils.Response
// @failure 500 {object} utils.Response
// @router /author/{authorID} [get]
func (api *API) handleGetAuthorByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "authorID"), 10, 64)
	if err != nil {
		utils.ResponseWithJSON(w, http.StatusBadRequest, nil, stack.Wrap(r.Context(), err, errors.ErrInvalidID))
		return
	}

	author, code, err := api.service.GetAuthorByID(r.Context(), id)
	utils.ResponseWithJSON(w, code, author, stack.Wrap(r.Context(), err))
}
