package usecase

import (
	"EffectiveMobile/internal/config"
	"EffectiveMobile/internal/domain/entity"
	"EffectiveMobile/internal/infrastructure/postgres/model"
	"EffectiveMobile/pkg/api/audd"
	"EffectiveMobile/pkg/api/lastfm"
	"EffectiveMobile/pkg/log"
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type repository interface {
	GetSongs(limit, offset int, filter entity.SongsFilter) ([]model.Song, error)
	GetSongText(groupName, songName string) (model.Song, error)
	AddSong(song entity.Song) error
	UpdateSong(song entity.Song) error
	DeleteSong(groupName, songName string) error
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

	// Получаем песни с учётом лимита и смещения
	songs, err := u.repository.GetSongs(songsPerPage, offset, filter)
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

func (u *Usecase) GetSongText(groupName, songName, verse string) (entity.SongText, error) {
	song, err := u.repository.GetSongText(groupName, songName)
	if err != nil {
		log.Logger.Warnf("Song not found: Group=%s, Song=%s", groupName, songName)
		return entity.SongText{}, err
	}
	log.Logger.Infof("Song retrieved: %s - %s", groupName, songName)

	// Разделяем текст песни на куплеты
	verses := strings.Split(song.Text, "\n\n")
	log.Logger.Debugf("Total verses: %d", len(verses))

	// Получаем параметр текущего куплета из URL
	if verse == "" {
		verse = "1"
	}
	log.Logger.Debugf("Verse parameter: %s", verse)

	verse = strings.Trim(verse, "/")
	verseNumber, err := strconv.Atoi(verse)
	if err != nil || verseNumber < 1 || verseNumber > len(verses) {
		log.Logger.Warnf("Invalid verse number: %s", verse)
		return entity.SongText{}, err
	}
	log.Logger.Infof("Displaying verse %d for song %s - %s", verseNumber, groupName, songName)

	// Текущий куплет
	currentVerse := verses[verseNumber-1]

	return entity.SongText{GroupName: song.GroupName, SongName: song.SongName, VerseNumber: verseNumber,
		Verse: currentVerse, ReleaseDate: song.ReleaseDate, Link: song.Link, FullText: song.Text}, nil
}

func (u *Usecase) AddSong(groupName, songName string) error {

	song := entity.Song{GroupName: groupName, SongName: songName}
	audDData, err := audd.GetAudDData(song.GroupName, song.SongName, u.conf.AuddAPI.AuddAPIKey, u.conf.AuddAPI.AuddAPIURL)
	if err != nil {
		log.Logger.Error("Error fetching song data from AudD: ", err)
		return err
	}
	log.Logger.Debugf("AudD data retrieved")

	lastFmData, err := lastfm.GetLastFmData(song.GroupName, song.SongName, u.conf.LastFMAPI.LastFMAPIKey, u.conf.LastFMAPI.LastFMAPIURL)
	if err != nil {
		log.Logger.Error("Error fetching song data from Last.fm: ", err)
		return err
	}
	log.Logger.Debugf("Last.fm data retrieved")

	if lastFmData.Track.Wiki.Published != "" {
		song.ReleaseDate, err = time.Parse("02 Jan 2006, 15:04", lastFmData.Track.Wiki.Published)
		if err != nil {
			log.Logger.Error("Error parsing release date: ", err)
			return err
		}
		log.Logger.Debugf("Parsed release date: %s", song.ReleaseDate)
	} else {
		log.Logger.Info("No release date found in Last.fm API response")
		song.ReleaseDate = time.Time{}
	}

	song.Text = audDData.Result[0].Lyrics
	log.Logger.Debugf("Fetched song lyrics")

	var media []audd.Media
	err = json.Unmarshal([]byte(audDData.Result[0].Media), &media)
	if err != nil {
		log.Logger.Error("Error parsing media field: ", err)
		return err
	}
	log.Logger.Debugf("Parsed media field: %+v", media)

	for _, m := range media {
		if m.Provider == "youtube" {
			song.Link = m.URL
			break
		}
	}
	if song.Link == "" {
		log.Logger.Warn("YouTube link not found")
	} else {
		log.Logger.Debugf("YouTube link found: %s", song.Link)
	}

	err = u.repository.AddSong(song)
	if err != nil {
		log.Logger.Warnf("Error adding song: %s", err)
		return err
	}

	log.Logger.Infof("Song inserted into database: %s - %s", song.GroupName, song.SongName)
	return nil
}

func (u *Usecase) UpdateSong(song entity.Song) error {
	err := u.repository.UpdateSong(song)
	if err != nil {
		log.Logger.Error("Error updating song in database:", err)
		return err
	}
	log.Logger.Infof("Song updated successfully: %+v", song)

	return nil
}

func (u *Usecase) DeleteSong(groupName, songName string) error {
	err := u.repository.DeleteSong(groupName, songName)
	if err != nil {
		log.Logger.Warnf("Error delete song: %s", err)
		return err
	}
	log.Logger.Infof("Song deleted successfully: Group=%s, Song=%s", groupName, songName)

	return nil
}
