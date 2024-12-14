package usecase

import (
	"EffectiveMobile/internal/config"
	"EffectiveMobile/internal/domain/entity"
	"EffectiveMobile/internal/infrastructure/postgres/model"
	"EffectiveMobile/pkg/db/conn"
	"EffectiveMobile/pkg/log"
	"strconv"

	"gorm.io/gorm"
)

type repository interface {
	GetSongs(db *gorm.DB, limit, offset int, filter entity.SongsFilter) ([]model.Song, error)
	GetSongText()
	AddSong()
	UpdateSong()
	DeleteSong()
}

type Usecase struct {
	repository
	conf config.APIConfig
}

func New(repo repository, conf config.APIConfig) *Usecase {
	return &Usecase{
		repo,
		conf,
	}
}

func (u *Usecase) GetSongs(filter entity.SongsFilter) ([]entity.SongsList, error) {
	if filter.Page == "" {
		filter.Page = "1"
	}

	log.Logger.Debugf("Filter parameters: group=%s, song=%s, releaseDate=%s, page=%s",
		filter.Group, filter.Song, filter.ReleaseDate, filter.Page)

	pageNumber, err := strconv.Atoi(filter.Page)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	// Указываем количество песен на одной странице
	songsPerPage := 20
	offset := (pageNumber - 1) * songsPerPage

	db, err := conn.InitDB("postgres://test_user:password@localhost:5432/test_db?sslmode=disable")
	if err != nil {
		log.Logger.Error("Failed to connect to the database:", err)
		return nil, err
	}
	log.Logger.Debug("Successfully connected to the database")

	// Получаем песни с учётом лимита и смещения
	songs, err := u.repository.GetSongs(db, songsPerPage, offset, filter)
	if err != nil {
		log.Logger.Error("Error fetching songs from database:", err)
		return nil, err
	}
	log.Logger.Infof("Fetched %d songs for page %d", len(songs), pageNumber)

	var songsList []entity.SongsList
	for _, song := range songs {
		songsList = append(songsList, entity.SongsList{GroupName: song.GroupName,
			SongName: song.SongName, ReleaseDate: song.ReleaseDate, Link: song.Link})
	}
	return songsList, nil
}

func (u *Usecase) GetSongText() {

}

func (u *Usecase) AddSongHandler() {

}

func (u *Usecase) UpdateSong() {

}

func (u *Usecase) DeleteSongHandler() {

}
