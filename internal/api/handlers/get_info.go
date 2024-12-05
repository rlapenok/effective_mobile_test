package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rlapenok/effective_mobile_test/internal/api/helpers"
	"github.com/rlapenok/effective_mobile_test/internal/domain/models"
	"github.com/rlapenok/effective_mobile_test/internal/infrastructe/state"
)

// GetInfoHandler godoc
//
//	@Summary		Get songs info
//	@Description	Get songs info from library
//	@Tags			Song Library
//	@Produce		json
//	@Param			filters[group]		query	string	false	"group filter"
//	@Param			filters[song]		query	string	false	"song filter"
//	@Param			filters[release_date_start]	query	string	false	"release start date filter (support gt,gte)"
//	@Param			filters[release_date_end]	query	string	false	"release end date filter (support lt,lte)"
//	@Param			filters[release_date]	query	string	false	"release date filter"
//	@Param			page	query	string	true	"page num"
//	@Param			limit	query	string	true	"limit num"
//
// @Success		200		{object}	helpers.GetFilteredSongs	"Response with filtered songs"
//
//	@Failure		400	{object}	helpers.ErrResponse	"Bad request"
//	@Failure		404	{object}	helpers.ErrResponse	"Not found"
//	@Failure		500	{object}	helpers.ErrResponse	"Internal error"
//	@Router			/get_info [get]
func GetInfo(c *gin.Context) {
	ctx := c.Request.Context()
	filtersQuery := c.QueryMap("filters")
	page := c.Query("page")
	limit := c.Query("limit")
	filters, err := models.NewFilters(filtersQuery, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Desc: err.Error()})
		return
	}
	songs, err := state.ServerState.GetInfo(ctx, filters)
	if err != nil {
		helpers.SendError(c, err)
		return
	}
	c.JSON(http.StatusOK, helpers.GetFilteredSongs{Songs: songs})
}
