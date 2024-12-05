package state

import (
	"context"

	"github.com/google/uuid"
	"github.com/rlapenok/effective_mobile_test/internal/domain/models"
	"github.com/rlapenok/effective_mobile_test/internal/domain/song_repository"
	"github.com/rlapenok/effective_mobile_test/internal/domain/swagger_client"
	"go.uber.org/zap"
)

var ServerState *state

type state struct {
	song_repo      song_repository.SongRepository
	swagger_client swagger_client.SwaggerClient
}

func New(song_repo song_repository.SongRepository, swagger_client swagger_client.SwaggerClient) {
	ServerState = &state{song_repo: song_repo, swagger_client: swagger_client}
}

func (s *state) AddSong(ctx context.Context, group string, song string) (*uuid.UUID, error) {
	zap.L().Debug("Getting Song Details...")
	song_details, err := s.swagger_client.GetSongDetail(group, song)
	if err != nil {
		return nil, err
	}
	zap.L().Debug("Song Details received")
	new_song, err := models.NewSong(group, song, song_details)
	if err != nil {
		return nil, err
	}
	zap.L().Debug("Adding a song to the database...")
	if err := s.song_repo.AddSong(ctx, new_song); err != nil {
		return nil, err
	}
	zap.L().Debug("The song has been added to the database")
	return &new_song.Id, nil
}

func (s *state) DeleteSong(ctx context.Context, id uuid.UUID) error {

	return s.song_repo.DeleteSong(ctx, id)
}

func (s *state) ChangeSong(ctx context.Context, id uuid.UUID, changes *models.Changes) error {
	return s.song_repo.ChangeSong(ctx, id, changes)
}

func (s *state) GetInfo(ctx context.Context, filters *models.Filters) ([]models.GetInfoSong, error) {
	return s.song_repo.GetInfo(ctx, filters)
}

func (s *state) GetLyrics(ctx context.Context, pagination *models.LyricsPagination) ([]models.Verse, error) {
	return s.song_repo.GetLyrics(ctx, pagination)
}
