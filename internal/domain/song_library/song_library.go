package song_library

import (
	"context"

	"github.com/google/uuid"
	"github.com/rlapenok/effective_mobile_test/internal/domain/models"
)

type SongLibrary interface {
	AddSong(ctx context.Context, group string, song string) error
	DeleteSong(ctx context.Context, id uuid.UUID) error
	ChangeSong(ctx context.Context, id uuid.UUID, changes *models.Changes) error
	GetInfo(ctx context.Context, filters *models.Filters) ([]models.GetInfoSong, error)
	GetLyrics(ctx context.Context, pagination *models.LyricsPagination) ([]models.Verse, error)
}
