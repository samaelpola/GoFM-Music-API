-- name: create-database
CREATE DATABASE IF NOT EXISTS musics;

-- name: create-musics-table
CREATE TABLE IF NOT EXISTS musics (
    music_id INTEGER PRIMARY KEY ,
    name VARCHAR(255) NOT NULL ,
    title VARCHAR(255) NOT NULL ,
    type VARCHAR(255) NOT NULL ,
    picture VARCHAR(255) NULL ,
    track VARCHAR(255) NULL
);

-- name: create-music
INSERT INTO musics (name, title, type, picture, track) VALUES(?, ?, ?, ?, ?);

-- name: update-picture-track
UPDATE musics SET picture = ?, track = ? WHERE music_id = ?;

-- name: find-music-by-id
SELECT * FROM musics WHERE music_id = ?;

-- name: find-music-by-type
SELECT * FROM musics WHERE type = ?;

-- name: check-music-already-exist
SELECT count(*) FROM musics WHERE name = ? and title = ?;

-- name: delete-music
DELETE FROM musics WHERE music_id = ?;

-- name: update-music
UPDATE musics SET name = ?, title = ?, type = ?, picture = ?, track = ? WHERE music_id = ?;

-- name: find-all-music
SELECT * FROM musics;

--name: drop-musics-table
DROP TABLE musics;
