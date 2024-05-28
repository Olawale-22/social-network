
-- +migrate Up
CREATE TABLE IF NOT EXISTS POSTS (
	ID INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	content TEXT NOT NULL,
    image TEXT,
    privacy TEXT CHECK (privacy IN ('public', 'private', 'mentions')) NOT NULL,
	created_at DATETIME NOT NULL,
	FOREIGN KEY (user_id) REFERENCES USERS(ID)
);