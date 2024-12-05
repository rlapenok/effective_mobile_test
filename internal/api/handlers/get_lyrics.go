package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rlapenok/effective_mobile_test/internal/api/helpers"
	"github.com/rlapenok/effective_mobile_test/internal/domain/models"
	"github.com/rlapenok/effective_mobile_test/internal/infrastructe/state"
)

// GetLyricsSongHandler godoc
//
//	@Summary		Get lyrics
//	@Description	Get lyrics with pagination
//	@Tags			Song Library
//	@Produce		json
//	@Param			id	path	string	true	"Song ID"
//	@Param			page	query	string	true	"page num"
//	@Param			limit	query	string	true	"limit num"
//
// @Success		200		{object}	helpers.GetLyricsResponse	"Response with filtered songs"
//
//	@Failure		400	{object}	helpers.ErrResponse	"Bad request"
//	@Failure		404	{object}	helpers.ErrResponse	"Not found"
//	@Failure		500	{object}	helpers.ErrResponse	"Internal error"
//	@Router			/lyrics/{id} [get]
func GetLyrics(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	page := c.Query("page")
	limit := c.Query("limit")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Desc: err.Error()})
		return
	}
	pagination, err := models.NewLyricsPagination(uuid, page, limit)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Desc: err.Error()})
		return
	}

	lyrics, err := state.ServerState.GetLyrics(ctx, pagination)
	if err != nil {
		helpers.SendError(c, err)
		return
	}
	c.JSON(http.StatusOK, helpers.GetLyricsResponse{Lyrics: lyrics})

}
