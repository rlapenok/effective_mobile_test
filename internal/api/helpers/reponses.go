package helpers

import (
	"github.com/google/uuid"
	"github.com/rlapenok/effective_mobile_test/internal/domain/models"
)

type AddSongResponse struct {
	Id uuid.UUID `json:"id" example:"fe8b200c-2fe6-4ced-82cd-875751f336fb"`
}

type ErrResponse struct {
	Desc string `json:"desc" example:"error description"`
}

type GetFilteredSongs struct {
	Songs []models.GetInfoSong `json:"songs"`
}

type GetLyricsResponse struct {
	Lyrics []models.Verse `json:"lyrics"`
}
