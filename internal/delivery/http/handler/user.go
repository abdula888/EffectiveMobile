package handler

type SongJSON struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group_name"`
	SongName    string `json:"song_name"`
	Text        string `json:"text"`
	ReleaseDate string `json:"releaseDate"`
	Link        string `json:"link"`
}
