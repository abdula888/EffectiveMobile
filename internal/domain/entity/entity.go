package entity

import "time"

type SongsFilter struct {
	Group       string
	Song        string
	ReleaseDate string
	Page        string
}

type SongsList struct {
	GroupName   string
	SongName    string
	ReleaseDate time.Time
	Link        string
}
