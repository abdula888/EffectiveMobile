package models

type Group struct {
	ID        int `gorm:"primaryKey"`
	GroupName string
	Song      []Song `gorm:"foreignKey:GroupID;references:ID;constraint:OnDelete:CASCADE"`
}

type Song struct {
	ID          int `gorm:"primaryKey"`
	GroupID     int
	GroupName   string `gorm:"-:all"`
	SongName    string `gorm:"type:varchar(100)"`
	Text        string
	ReleaseDate string `gorm:"type:date"`
	Link        string `gorm:"type:varchar(100)"`
}

type SongJSON struct {
	ID          int    `json:"id"`
	GroupName   string `json:"group_name"`
	SongName    string `json:"song_name"`
	Text        string `json:"text"`
	ReleaseDate string `json:"releaseDate"`
	Link        string `json:"link"`
}
