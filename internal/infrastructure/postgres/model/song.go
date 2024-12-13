package model

import "time"

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
	ReleaseDate time.Time `gorm:"type:date"`
	Link        string    `gorm:"type:varchar(100)"`
}
