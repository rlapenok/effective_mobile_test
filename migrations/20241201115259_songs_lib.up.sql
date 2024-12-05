-- Add up migration script here
CREATE TABLE IF NOT EXISTS songs(
    id UUID PRIMARY KEY,
    group_name TEXT NOT NULL,
    song TEXT NOT NULL UNIQUE,
    link TEXT NOT NULL UNIQUE,
    release_date DATE NOT NULL
);
CREATE TABLE IF NOT EXISTS lyrics(
    id UUID NOT NULL,
    verse_number INT NOT NULL,
    text TEXT NOT NULL,
    FOREIGN KEY (id) REFERENCES songs(id) ON DELETE CASCADE
)