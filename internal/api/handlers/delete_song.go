package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rlapenok/effective_mobile_test/internal/api/helpers"
	"github.com/rlapenok/effective_mobile_test/internal/infrastructe/state"
)

// DeleteSongHandler godoc
//
//	@Summary		Delete a song
//	@Description	Deletes a song by its ID
//	@Tags			Song Library
//	@Produce		json
//	@Param			id	path	string	true	"Song ID"
//	@Success		200
//	@Failure		400	{object}	helpers.ErrResponse	"Bad request"
//	@Failure		404	{object}	helpers.ErrResponse	"Not found"
//	@Failure		500	{object}	helpers.ErrResponse	"Internal error"
//	@Router			/delete_song/{id} [delete]
func DeleteSong(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, helpers.ErrResponse{Desc: err.Error()})
		return
	}
	if err := state.ServerState.DeleteSong(ctx, uuid); err != nil {
		helpers.SendError(c, err)
	}

}
