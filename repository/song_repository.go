package repository

import (
	"EffectiveMobile/models"
	"database/sql"
	"strconv"
)

func GetSongsWithPagination(db *sql.DB, limit, offset int, group, song, releaseDate string) ([]models.Song, error) {
	var songs []models.Song

	query := "SELECT group_name, song, text, releaseDate, link FROM songs WHERE 1=1"

	var args []interface{}
	argIndex := 1 // Индекс для параметров

	if group != "" {
		query += " AND group_name LIKE $" + strconv.Itoa(argIndex)
		args = append(args, "%"+group+"%")
		argIndex++
	}

	if song != "" {
		query += " AND song LIKE $" + strconv.Itoa(argIndex)
		args = append(args, "%"+song+"%")
		argIndex++
	}

	if releaseDate != "" {
		query += " AND releaseDate = $" + strconv.Itoa(argIndex)
		args = append(args, releaseDate)
		argIndex++
	}

	// Добавляем ORDER BY и передаем параметры для LIMIT и OFFSET
	query += " ORDER BY group_name LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
	args = append(args, limit, offset)

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var song models.Song
		err := rows.Scan(&song.GroupName, &song.Song, &song.Text, &song.ReleaseDate, &song.Link)
		if err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	return songs, nil
}

func GetSongByName(db *sql.DB, groupName string, songName string) (*models.Song, error) {
	var song models.Song
	query := "SELECT group_name, song, text, releaseDate, link FROM songs WHERE group_name = $1 AND song = $2"
	err := db.QueryRow(query, groupName, songName).Scan(&song.GroupName, &song.Song, &song.Text, &song.ReleaseDate, &song.Link)
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func UpdateSongByName(db *sql.DB, song models.Song) error {
	// Сохраняем данные в БД
	_, err := db.Exec(
		"UPDATE songs SET text = $1, releaseDate = $2, link = $3 WHERE group_name = $4 AND song = $5",
		song.Text,
		song.ReleaseDate,
		song.Link,
		song.GroupName,
		song.Song,
	)

	if err != nil {
		return err
	}
	return nil
}
