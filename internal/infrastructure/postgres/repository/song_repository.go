package repository

import (
	"EffectiveMobile/internal/domain/entity"
	"EffectiveMobile/internal/infrastructure/postgres/model"
	"errors"

	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB
}

func New(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetSongs(limit, offset int, filter entity.SongsFilter) ([]model.Song, error) {
	var songs []model.Song
	groupName, songName := "%"+filter.Group+"%", "%"+filter.Song+"%"

	rows, err := r.db.Table("songs s").Joins("join groups g on s.group_id = g.id").
		Select("g.group_name, s.song_name, s.text, s.release_date, s.link").Limit(limit).Offset(offset).
		Where("s.song_name LIKE ? AND g.group_name LIKE ?", songName, groupName).Order("group_id").Rows()

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var song model.Song

		err := rows.Scan(&song.GroupName, &song.SongName, &song.Text, &song.ReleaseDate, &song.Link)
		if err != nil {
			return nil, err
		}

		songs = append(songs, song)
	}
	return songs, nil
}

func (r *Repo) GetSongText(groupName string, songName string) (model.Song, error) {
	var song model.Song

	err := r.db.Table("songs s").Joins("join groups g on s.group_id = g.id").
		Where("s.song_name = ? AND g.group_name = ?", songName, groupName).Last(&song).Error

	if err != nil {
		return model.Song{}, err
	}
	song.GroupName = groupName

	return song, nil
}

func (r *Repo) AddSong(song entity.Song) error {
	group := &model.Group{GroupName: song.GroupName}
	err := r.db.First(&group, "group_name = ?", group.GroupName).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = r.db.Create(&group).First(&group, "group_name = ?", group.GroupName).Error
		if err != nil {
			return err
		}
	}
	songDB := model.Song{GroupID: group.ID, SongName: song.SongName, Text: song.Text,
		ReleaseDate: song.ReleaseDate, Link: song.Link}
	err = r.db.Create(&songDB).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repo) UpdateSong(song entity.Song) error {
	songDB := model.Song{GroupName: song.GroupName, SongName: song.SongName,
		Text: song.Text, ReleaseDate: song.ReleaseDate, Link: song.Link}

	r.db.Model(&songDB).Where("group_id = ? AND song_name = ?", songDB.GroupID, songDB.SongName).Updates(model.Song{Text: songDB.Text, ReleaseDate: songDB.ReleaseDate, Link: songDB.Link})

	return nil
}

func (r *Repo) DeleteSong(groupName, songName string) error {
	var group model.Group
	err := r.db.First(&group, "group_name = ?", groupName).Error
	if err != nil {
		return err
	}
	err = r.db.Where("song_name = ? AND group_id = ?", songName, group.ID).Delete(&model.Song{}).Error
	if err != nil {
		return err
	}
	return nil
}
