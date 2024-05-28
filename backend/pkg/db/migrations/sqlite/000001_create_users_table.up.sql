
-- +migrate Up
CREATE TABLE IF NOT EXISTS USERS (
	ID INTEGER PRIMARY KEY AUTOINCREMENT, 
	email TEXT NOT NULL,
	password TEXT NOT NULL,
	firstname TEXT NOT NULL,
	lastname TEXT NOT NULL,
	birthdate TEXT NOT NULL,
	avatar TEXT,
	nickname TEXT,
	aboutme TEXT,
	isconnected BOOLEAN NOT NULL,
	ispublic BOOLEAN NOT NULL,
	isrequested BOOLEAN NOT NULL
);