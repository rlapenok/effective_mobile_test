package helpers

// AddSongRequest represents the request body for adding a song
type AddSongRequest struct {
	Group string `json:"group" validate:"required"`
	Song  string `json:"song" validate:"required"`
}
