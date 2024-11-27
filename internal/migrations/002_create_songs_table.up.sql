CREATE TABLE songs (
    id SERIAL PRIMARY KEY,
    group_id INT,
    song_name VARCHAR(100),
    text TEXT,
    releaseDate DATE,
    link VARCHAR(200),
    CONSTRAINT fk_group FOREIGN KEY (group_id) REFERENCES groups(group_id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS unique_group_song ON songs(group_id, song_name);
