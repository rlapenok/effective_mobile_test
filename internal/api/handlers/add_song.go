package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rlapenok/effective_mobile_test/internal/api/helpers"
	"github.com/rlapenok/effective_mobile_test/internal/infrastructe/state"
)

// AddSongHandler godoc
//
// @Summary		Add a new song
// @Description	Adds a new song to the library
// @Tags			Song Library
// @Accept			json
// @Produce		json
// @Param			request	body		helpers.AddSongRequest	true	"Request body containing song details"
// @Success		200		{object}	helpers.AddSongResponse	"Response with song ID"
// @Failure		400		{object}	helpers.ErrResponse		"Bad request"
// @Failure		404		{object}	helpers.ErrResponse		"Not found"
// @Failure		500		{object}	helpers.ErrResponse		"Internal error"
// @Router			/add_song [post]
func AddSong(c *gin.Context) {
	ctx := c.Request.Context()
	requset, err := helpers.ExtractorJsonWithValidation[helpers.AddSongRequest](c)
	if err != nil {
		helpers.SendError(c, err)
		return
	}
	song_id, err := state.ServerState.AddSong(ctx, requset.Group, requset.Song)
	if err != nil {
		helpers.SendError(c, err)
		return
	}
	c.JSON(http.StatusOK, helpers.AddSongResponse{Id: *song_id})
}
