package repository

import (
	"EffectiveMobile/internal/models"
	"errors"
	"time"

	"gorm.io/gorm"
)

func GetSongsList(db *gorm.DB, limit, offset int, groupName, songName, releaseDate string) ([]models.Song, error) {
	var songs []models.Song
	groupName, songName = "%"+groupName+"%", "%"+songName+"%"

	rows, err := db.Table("songs s").Joins("join groups g on s.group_id = g.id").
		Select("g.group_name, s.song_name, s.text, s.release_date, s.link").Limit(limit).Offset(offset).
		Where("s.song_name LIKE ? AND g.group_name LIKE ?", songName, groupName).Order("group_id").Rows()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var song models.Song
		var releaseDate time.Time

		err := rows.Scan(&song.GroupName, &song.SongName, &song.Text, &releaseDate, &song.Link)
		if err != nil {
			return nil, err
		}

		// Преобразуем дату в формат 2006-01-02
		song.ReleaseDate = releaseDate.Format("2006-01-02")

		songs = append(songs, song)
	}
	return songs, nil
}

func GetSongText(db *gorm.DB, groupName string, songName string) (*models.Song, error) {
	var song *models.Song

	db.Table("songs s").Joins("join groups g on s.group_id = g.id").
		Where("s.song_name = ? AND g.group_name = ?", songName, groupName).Last(&song)

	song.GroupName = groupName

	return song, nil
}

func UpdateSong(db *gorm.DB, song models.Song) error {
	db.Model(&song).Where("group_id = ? AND song_name = ?", song.GroupID, song.SongName).Updates(models.Song{Text: song.Text, ReleaseDate: song.ReleaseDate, Link: song.Link})

	return nil
}

func AddSong(db *gorm.DB, song models.Song) error {
	group := &models.Group{GroupName: song.GroupName}
	err := db.First(&group, "group_name = ?", group.GroupName).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = db.Create(&group).First(&group, "group_name = ?", group.GroupName).Error
		if err != nil {
			return err
		}
	}
	song.GroupID = group.ID
	err = db.Create(&song).Error
	if err != nil {
		return err
	}
	return nil
}

func DeleteSong(db *gorm.DB, songName, groupName string) error {
	var group models.Group
	err := db.First(&group, "group_name = ?", groupName).Error
	if err != nil {
		return err
	}
	err = db.Where("song_name = ? AND group_id = ?", songName, group.ID).Delete(&models.Song{}).Error
	if err != nil {
		return err
	}
	return nil
}
