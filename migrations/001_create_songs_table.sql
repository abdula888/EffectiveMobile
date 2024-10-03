CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    group_name VARCHAR(100),
    song VARCHAR(100),
    text TEXT,
    releaseDate VARCHAR(100),
    link VARCHAR(200)
);

-- Добавляем уникальный индекс для полей group_name и song
CREATE UNIQUE INDEX unique_group_song ON songs(group_name, song);

