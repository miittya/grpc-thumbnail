CREATE TABLE IF NOT EXISTS thumbnails
(
    id        INTEGER PRIMARY KEY,
    video_url TEXT NOT NULL UNIQUE,
    thumbnail BLOB NOT NULL
);