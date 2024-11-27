package repository

import (
	"EffectiveMobile/internal/models"
	"database/sql"
	"strconv"
)

func GetSongsWithPagination(db *sql.DB, limit, offset int, group, song, releaseDate string) ([]models.Song, error) {
	var songs []models.Song

	query := `
        SELECT 
            g.group_name, 
            s.song_name, 
            s.text, 
            s.releaseDate, 
            s.link
        FROM 
            songs s
        JOIN 
            groups g 
        ON 
            s.group_id = g.group_id
        WHERE 1=1
    `

	var args []interface{}
	argIndex := 1

	// Фильтр по названию группы
	if group != "" {
		query += " AND g.group_name LIKE $" + strconv.Itoa(argIndex)
		args = append(args, "%"+group+"%")
		argIndex++
	}

	// Фильтр по названию песни
	if song != "" {
		query += " AND s.song_name LIKE $" + strconv.Itoa(argIndex)
		args = append(args, "%"+song+"%")
		argIndex++
	}

	// Фильтр по дате релиза
	if releaseDate != "" {
		query += " AND s.releaseDate = $" + strconv.Itoa(argIndex)
		args = append(args, releaseDate)
		argIndex++
	}

	// Сортировка и лимит
	query += " ORDER BY g.group_name LIMIT $" + strconv.Itoa(argIndex) + " OFFSET $" + strconv.Itoa(argIndex+1)
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
	query := `
    SELECT 
        g.group_name, 
        s.song_name, 
        s.text, 
        s.releaseDate, 
        s.link
    FROM 
        songs s
    JOIN 
        groups g 
    ON 
        s.group_id = g.group_id
    WHERE 
        g.group_name = $1 AND s.song_name = $2;
`
	err := db.QueryRow(query, groupName, songName).Scan(&song.GroupName, &song.Song, &song.Text, &song.ReleaseDate, &song.Link)
	if err != nil {
		return nil, err
	}
	return &song, nil
}

func UpdateSongByName(db *sql.DB, song models.Song) error {
	// Сохраняем данные в БД
	_, err := db.Exec(
		`UPDATE songs SET text = $1, releaseDate = $2, link = $3 
		WHERE group_id = (SELECT group_id FROM groups WHERE group_name = $4) AND song_name = $5`,
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
