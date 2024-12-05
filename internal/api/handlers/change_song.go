package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rlapenok/effective_mobile_test/internal/api/helpers"
	"github.com/rlapenok/effective_mobile_test/internal/domain/models"
	"github.com/rlapenok/effective_mobile_test/internal/infrastructe/state"
)

// ChangeSongHandler godoc
//
//	@Summary		Change song data
//	@Description	Changing song data by its ID
//	@Tags			Song Library
//	@Produce		json
//	@Param			id							path	string	true	"Song ID"
//
//	@Param			changes[new_group_name]		query	string	false	"The new group name"
//	@Param			changes[new_song_name]		query	string	false	"The new song name"
//	@Param			changes[new_link]			query	string	false	"The new song link"
//	@Param			changes[new_release_date]	query	string	false	"The new release date"
//
//	@Success		200
//	@Failure		400	{object}	helpers.ErrResponse	"Bad request"
//	@Failure		404	{object}	helpers.ErrResponse	"Not found"
//	@Failure		500	{object}	helpers.ErrResponse	"Internal error"
//	@Router			/change_song/{id} [patch]
func ChangeSong(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Desc: err.Error()})
		return
	}
	changesQuery := c.QueryMap("changes")
	changes, err := models.NewChanges(changesQuery)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Desc: err.Error()})
		return
	}
	if err := state.ServerState.ChangeSong(ctx, uuid, changes); err != nil {
		helpers.SendError(c, err)
		return
	}

}
