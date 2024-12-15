package entity

import "time"

type Song struct {
	ID          int
	GroupName   string
	SongName    string
	Text        string
	ReleaseDate time.Time
	Link        string
}

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

type SongText struct {
	GroupName   string
	SongName    string
	VerseNumber int
	Verse       string
	ReleaseDate time.Time
	Link        string
	FullText    string
}
