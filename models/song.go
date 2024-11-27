package models

type Song struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group"`
	Song        string `json:"song"`
	Text        string `json:"text"`
	ReleaseDate string `json:"releaseDate"`
	Link        string `json:"link"`
}
