package handler

import "time"

type SongJSON struct {
	ID          int       `json:"id"`
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	Text        string    `json:"text"`
	ReleaseDate time.Time `json:"releaseDate"`
	Link        string    `json:"link"`
}

type DataJSON struct {
	GroupName   string    `json:"group_name"`
	SongName    string    `json:"song_name"`
	ReleaseDate time.Time `json:"releaseDate"`
	Link        string    `json:"link"`
}
