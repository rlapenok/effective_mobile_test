package song_repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/rlapenok/effective_mobile_test/internal/domain/models"
)

type SongRepository interface {
	AddSong(ctx context.Context, song *models.Song) error
	DeleteSong(ctx context.Context, id uuid.UUID) error
	ChangeSong(ctx context.Context, id uuid.UUID, changes *models.Changes) error
	GetInfo(ctx context.Context, filters *models.Filters) ([]models.GetInfoSong, error)
	GetLyrics(ctx context.Context, pagination *models.LyricsPagination) ([]models.Verse, error)
}
