package models

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type SongDetails struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type GetInfoSong struct {
	Id          uuid.UUID `db:"id" json:"id"`
	Group       string    `db:"group_name" json:"group"`
	Song        string    `db:"song" json:"song"`
	RaleaseDate time.Time `db:"release_date" json:"release_date"`
	Link        string    `db:"link" json:"link"`
}

type Song struct {
	Id          uuid.UUID
	Group       string
	Song        string
	RealeseDate time.Time
	Lyrics      []string
	Link        string
}

func NewSong(group string, song string, song_details *SongDetails) (*Song, error) {
	id := uuid.New()
	date, err := time.Parse("02.06.2006", song_details.ReleaseDate)
	if err != nil {
		return nil, err
	}
	lyrics := strings.Split(song_details.Text, "\n\n")
	return &Song{
		Id:          id,
		Group:       group,
		Song:        song,
		RealeseDate: date,
		Lyrics:      lyrics,
		Link:        song_details.Link,
	}, nil

}

type Changes struct {
	Group       *string
	Song        *string
	Link        *string
	ReleaseDate *time.Time
}

func NewChanges(changes map[string]string) (*Changes, error) {
	var group *string
	var song *string
	var link *string
	var releaseDate *time.Time
	if len(changes) == 0 {
		return nil, errors.New("not set query params")
	}
	for k, v := range changes {
		if len(v) > 0 {
			switch k {
			case "new_group_name":
				{

					group = &v
				}
			case "new_song_name":
				{
					song = &v
				}
			case "new_link":
				{
					link = &v
				}
			case "new_release_date":
				{
					date, err := time.Parse("02.06.2006", v)
					if err != nil {
						return nil, err
					}
					releaseDate = &date
				}
			}
		}
	}
	return &Changes{Group: group, Song: song, Link: link, ReleaseDate: releaseDate}, nil
}

type Filters struct {
	Group          *string
	Song           *string
	Link           *string
	ReleaseDate    *time.Time
	GtReleaseDate  *time.Time
	LtReleaseDate  *time.Time
	GteReleaseDate *time.Time
	LteReleaseDate *time.Time
	Offset         int
	Limit          int
}

func NewFilters(filters map[string]string, page_str string, limit_str string) (*Filters, error) {
	var group *string
	var song *string
	var link *string
	var releaseDate *time.Time
	//>
	var gtReleaseDate *time.Time
	//<
	var ltReleaseDate *time.Time
	//>=
	var gteReleaseDate *time.Time
	//<=
	var lteReleaseDate *time.Time

	page, err := strconv.Atoi(page_str)
	if err != nil {
		return nil, err
	}
	limit, err := strconv.Atoi(limit_str)
	if err != nil {
		return nil, err
	}

	offset := (page - 1) * limit

	for k, v := range filters {
		if len(v) > 0 {
			switch k {
			case "group":
				{

					group = &v
				}
			case "song":
				{
					song = &v
				}
			case "link":
				{
					link = &v
				}
			case "release_date":
				{
					date, err := time.Parse("02.06.2006", v)
					if err != nil {
						return nil, err
					}
					releaseDate = &date
				}

				//default gte
			case "release_date_start":
				{
					if strings.HasPrefix(v, "gt:") {
						date, err := time.Parse("02.06.2006", v[3:])
						if err != nil {
							return nil, err
						}
						gtReleaseDate = &date
					} else if strings.HasPrefix(v, "gte:") {
						date, err := time.Parse("02.06.2006", v[4:])
						if err != nil {
							return nil, err
						}
						gteReleaseDate = &date
					} else {
						date, err := time.Parse("02.06.2006", v)
						if err != nil {
							return nil, err
						}
						gteReleaseDate = &date
					}
				}
				//default lte
			case "release_date_end":
				{
					if strings.HasPrefix(v, "lt:") {
						date, err := time.Parse("02.06.2006", v[3:])
						if err != nil {
							return nil, err
						}
						ltReleaseDate = &date
					} else if strings.HasPrefix(v, "lte:") {
						date, err := time.Parse("02.06.2006", v[4:])
						if err != nil {
							return nil, err
						}
						lteReleaseDate = &date
					} else {
						date, err := time.Parse("02.06.2006", v)
						if err != nil {
							return nil, err
						}
						lteReleaseDate = &date
					}
				}
			}
		}
	}
	return &Filters{Group: group,
		Song:           song,
		Link:           link,
		ReleaseDate:    releaseDate,
		GtReleaseDate:  gtReleaseDate,
		LtReleaseDate:  ltReleaseDate,
		GteReleaseDate: gteReleaseDate,
		LteReleaseDate: lteReleaseDate,
		Offset:         offset,
		Limit:          limit}, nil
}

type LyricsPagination struct {
	Id     uuid.UUID
	Limit  int
	Offset int
}

func NewLyricsPagination(id uuid.UUID, page_str string, limit_str string) (*LyricsPagination, error) {
	page, err := strconv.Atoi(page_str)
	if err != nil {
		return nil, err
	}
	limit, err := strconv.Atoi(limit_str)
	if err != nil {
		return nil, err
	}
	offset := (page - 1) * limit

	return &LyricsPagination{
		Id:     id,
		Limit:  limit,
		Offset: offset,
	}, nil

}

type Verse struct {
	VerseNumber int    `db:"verse_number" json:"verse_number"`
	Text        string `db:"text" json:"text"`
}
