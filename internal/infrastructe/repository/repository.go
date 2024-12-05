package repository

import (
	"context"
	"database/sql"
	"errors"
	"os"

	"github.com/Masterminds/squirrel"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/rlapenok/effective_mobile_test/internal/domain/models"
	"github.com/rlapenok/effective_mobile_test/internal/domain/song_repository"
	"go.uber.org/zap"
)

type Repository struct {
	db *sqlx.DB
	sq squirrel.StatementBuilderType
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db, sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)}
}

func (repo *Repository) RunMigration(path string) {
	path = "file:" + path
	driver, err := postgres.WithInstance(repo.db.DB, &postgres.Config{})
	if err != nil {
		zap.L().Error("Error while crete driver for run migrations", zap.Error(err))
		os.Exit(1)
	}
	migrator, err := migrate.NewWithDatabaseInstance(path, "postgres", driver)
	if err != nil {
		zap.L().Error("Error while crete instance for run migration", zap.Error(err))
		os.Exit(1)
	}
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		zap.L().Error("Error while run migrations", zap.Error(err))
		os.Exit(1)
	}
	zap.L().Info("Migration applied")

}

func (repo *Repository) AddSong(ctx context.Context, song *models.Song) error {

	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		zap.L().Error("Error while begin transaction", zap.Error(err))
		return err
	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = song_repository.SongRepositoryError{Err: errors.Join(err, rollbackErr)}
				zap.L().Error("Error while rollback tranaction", zap.Error(err))
			}
		}
	}()
	query1, args1, err := repo.sq.Insert("songs").
		Columns("id", "group_name", "song", "link", "release_date").
		Values(song.Id, song.Group, song.Song, song.Link, song.RealeseDate).
		ToSql()
	if err != nil {
		zap.L().Error("Error while create query", zap.Error(err))
		return song_repository.SongRepositoryError{Err: err}
	}
	_, err = tx.ExecContext(ctx, query1, args1...)
	if err != nil {
		zap.L().Error("Error while execute query", zap.Error(err))
		return song_repository.SongRepositoryError{Err: err}
	}
	queryBuilder := repo.sq.Insert("lyrics").Columns("id", "verse_number", "text")
	for i, lyric := range song.Lyrics {
		queryBuilder = queryBuilder.Values(song.Id, i+1, lyric)
	}
	query2, args2, err := queryBuilder.ToSql()
	if err != nil {
		zap.L().Error("Error while create query", zap.Error(err))
		return song_repository.SongRepositoryError{Err: err}
	}
	_, err = tx.ExecContext(ctx, query2, args2...)
	if err != nil {
		zap.L().Error("Error while execute query", zap.Error(err))
		return song_repository.SongRepositoryError{Err: err}
	}

	if err := tx.Commit(); err != nil {
		zap.L().Error("Error while commit tranaction", zap.Error(err))
		return song_repository.SongRepositoryError{Err: err}
	}
	return nil
}

func (repo *Repository) DeleteSong(ctx context.Context, id uuid.UUID) error {
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		zap.L().Error("Error while begin transaction", zap.Error(err))

	}
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = song_repository.SongRepositoryError{Err: errors.Join(err, rollbackErr)}
				zap.L().Error("Error while rollback transaction", zap.Error(err))

			}
		}
	}()
	query, args, err := repo.sq.Delete("songs").
		Where(squirrel.Eq{"id": id}).
		ToSql()
	if err != nil {
		zap.L().Error("Error while create query", zap.Error(err))

		return song_repository.SongRepositoryError{Err: err}
	}
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		zap.L().Error("Error while execute query", zap.Error(err))
		return song_repository.SongRepositoryError{Err: err}
	}

	if err != nil {
		zap.L().Error("Error while get affected rows", zap.Error(err))
		return song_repository.SongRepositoryError{Err: err}

	}
	if err := tx.Commit(); err != nil {
		zap.L().Error("Error while commit transaction", zap.Error(err))
		return song_repository.SongRepositoryError{Err: err}
	}
	numRows, err := res.RowsAffected()
	if numRows == 0 {
		return errors.New("not found")
	}
	return nil
}

func (repo *Repository) ChangeSong(ctx context.Context, id uuid.UUID, changes *models.Changes) error {
	queryBuilder := repo.sq.Update("songs").Where(squirrel.Eq{"id": id})

	if changes.Group != nil {
		queryBuilder = queryBuilder.Set("group_name", *changes.Group)
	}
	if changes.Song != nil {
		queryBuilder = queryBuilder.Set("song", *changes.Song)

	}
	if changes.Link != nil {
		queryBuilder = queryBuilder.Set("link", *changes.Link)

	}
	if changes.ReleaseDate != nil {
		queryBuilder = queryBuilder.Set("release_date", *changes.ReleaseDate)

	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		zap.L().Error("Error while create query", zap.Error(err))
		return err
	}
	tx, err := repo.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	defer func() {
		if err != nil {
			rollbackErr := tx.Rollback()
			if rollbackErr != nil {
				err = song_repository.SongRepositoryError{Err: errors.Join(err, rollbackErr)}
				zap.L().Error("Error rollback transaction", zap.Error(err))
			}
		}
	}()
	res, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		zap.L().Error("Error execute query", zap.Error(err))
		return song_repository.SongRepositoryError{Err: err}
	}
	if err := tx.Commit(); err != nil {
		zap.L().Error("Error while commit tranaction", zap.Error(err))
		return song_repository.SongRepositoryError{Err: err}
	}
	numRows, err := res.RowsAffected()
	if numRows == 0 {
		return errors.New("not found")
	}

	return nil
}

func (repo *Repository) GetInfo(ctx context.Context, filters *models.Filters) ([]models.GetInfoSong, error) {
	queryBuilder := repo.sq.Select("*").From("songs")

	if filters.Group != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"group_name": *filters.Group})
	}
	if filters.Song != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"song": *filters.Song})

	}
	if filters.Link != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"link": *filters.Link})

	}
	if filters.ReleaseDate != nil {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"release_date": *filters.ReleaseDate})

	}
	if filters.GtReleaseDate != nil {
		queryBuilder = queryBuilder.Where(squirrel.Gt{"release_date": *filters.GtReleaseDate})
	}
	if filters.GteReleaseDate != nil {
		queryBuilder = queryBuilder.Where(squirrel.GtOrEq{"release_date": *filters.GteReleaseDate})
	}
	if filters.LtReleaseDate != nil {
		queryBuilder = queryBuilder.Where(squirrel.Lt{"release_date": *filters.LtReleaseDate})
	}
	if filters.LteReleaseDate != nil {
		queryBuilder = queryBuilder.Where(squirrel.LtOrEq{"release_date": *filters.LteReleaseDate})
	}

	query, args, err := queryBuilder.Limit(uint64(filters.Limit)).Offset(uint64(filters.Offset)).ToSql()
	if err != nil {
		zap.L().Error("Error while create query", zap.Error(err))
		return nil, err
	}
	rows, err := repo.db.QueryxContext(ctx, query, args...)
	defer func() {
		if rows != nil {
			err = rows.Close()
			if err != nil {
				zap.L().Error("Error while closing rows", zap.Error(err))
			}
		}
	}()
	if err != nil {
		zap.L().Error("Error execute query", zap.Error(err))
		return nil, song_repository.SongRepositoryError{Err: err}
	}
	results := make([]models.GetInfoSong, 0)
	for rows.Next() {
		var song models.GetInfoSong
		if err := rows.StructScan(&song); err != nil {
			return nil, err
		}
		results = append(results, song)
	}

	return results, nil
}

func (repo *Repository) GetLyrics(ctx context.Context, pagination *models.LyricsPagination) ([]models.Verse, error) {
	query, args, err := repo.sq.Select("verse_number,text").
		From("lyrics").
		Where(squirrel.Eq{"id": pagination.Id}).
		Limit(uint64(pagination.Limit)).
		Offset(uint64(pagination.Offset)).
		ToSql()
	if err != nil {
		zap.L().Error("Error while create query", zap.Error(err))
		return nil, err
	}
	rows, err := repo.db.QueryxContext(ctx, query, args...)
	defer func() {
		if rows != nil {
			err = rows.Close()
			if err != nil {
				zap.L().Error("Error while closing rows", zap.Error(err))
			}
		}
	}()
	if err != nil {
		zap.L().Error("Error execute query", zap.Error(err))
		return nil, song_repository.SongRepositoryError{Err: err}
	}
	results := make([]models.Verse, 0)
	for rows.Next() {
		var verse models.Verse
		if err := rows.StructScan(&verse); err != nil {
			return nil, err
		}
		results = append(results, verse)
	}

	return results, nil
}
