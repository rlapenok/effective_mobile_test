package song_repository

type SongRepositoryError struct {
	Err error
}

func (e SongRepositoryError) Error() string {
	return e.Err.Error()
}
