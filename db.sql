CREATE TABLE channels (
    id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE posts (
    channel_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    rating INTEGER NOT NULL,
    description TEXT NOT NULL,
    tags TEXT NOT NULL,
    video_id TEXT PRIMARY KEY,
    video_path TEXT NOT NULL,
    thumbnail_id TEXT NOT NULL,
    thumbnail_path TEXT NOT NULL,
    g_thumbnail_id TEXT,
    g_thumbnail_path TEXT,
    FOREIGN KEY (channel_id) REFERENCES channels(id)
)