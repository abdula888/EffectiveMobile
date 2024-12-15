package http

import "time"

type SongJSON struct {
	ID          int       `json:"id"`
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	Text        string    `json:"text"`
	ReleaseDate time.Time `json:"releaseDate"`
	Link        string    `json:"link"`
}

type SongTextJSON struct {
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	VerseNumber int       `json:"verse_number"`
	Verse       string    `json:"verse"`
	ReleaseDate time.Time `json:"releaseDate"`
	Link        string    `json:"link"`
	FullText    string    `json:"full_text"`
}
type SongsListJSON struct {
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"releaseDate"`
	Link        string    `json:"link"`
}
