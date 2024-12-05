package swagger_client

import "github.com/rlapenok/effective_mobile_test/internal/domain/models"

type SwaggerClient interface {
	GetSongDetail(group string, song string) (*models.SongDetails, error)
}
