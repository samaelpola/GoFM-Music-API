package models

import (
	"errors"
	"github.com/samaelpola/GoFM-Music-API/internal/repository"
	"log"
)

type Music struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	Picture string `json:"picture"`
	Track   string `json:"track"`
}

const (
	GORAP  = "GO-RAP"
	GOPOP  = "GO-POP"
	GOROCK = "GO-ROCK"
	GOSLOW = "GO-SLOW"
)

var ErrMusicAlreadyExist = errors.New("music already Exist")

func (m Music) checkMusicAlreadyExist(goFmDb *repository.DB, name, title string) (bool, error) {
	res, err := goFmDb.GetDotSql().Query(goFmDb.GetSqlDb(), "check-music-already-exist", name, title)
	if err != nil {
		return false, err
	}
	defer res.Close()

	var count int
	if res.Next() {
		if err = res.Scan(&count); err != nil {
			return false, err
		}
	}

	return count > 0, nil
}

func (m Music) validateMusic(goFmDb *repository.DB, music Music) error {
	existEmail, err := m.checkMusicAlreadyExist(goFmDb, music.Name, music.Title)

	if err != nil {
		return err
	}

	if existEmail {
		return ErrMusicAlreadyExist
	}

	return nil
}

func (m Music) Create(goFmDb *repository.DB) (Music, error) {
	if err := m.validateMusic(goFmDb, m); err != nil {
		return Music{}, err
	}

	res, err := goFmDb.GetDotSql().Exec(
		goFmDb.GetSqlDb(),
		"create-music",
		m.Name,
		m.Title,
		m.Type,
		m.Picture,
		m.Track,
	)

	if err != nil {
		return Music{}, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return Music{}, err
	}

	return m.GetMusic(goFmDb, id)
}

func (m Music) GetMusic(goFmDb *repository.DB, musicId int64) (Music, error) {
	res, err := goFmDb.GetDotSql().QueryRow(goFmDb.GetSqlDb(), "find-music-by-id", musicId)

	if err != nil {
		return Music{}, err
	}

	if err = res.Scan(&m.ID, &m.Name, &m.Title, &m.Type, &m.Picture, &m.Track); err != nil {
		return Music{}, err
	}

	return m, nil
}

func (m Music) GetMusicByType(goFmDb *repository.DB, typeOfMusic string) ([]Music, error) {
	res, err := goFmDb.GetDotSql().Query(goFmDb.GetSqlDb(), "find-music-by-type", typeOfMusic)
	var listMusic []Music

	if err != nil {
		return []Music{}, err
	}

	for res.Next() {
		err = res.Scan(&m.ID, &m.Name, &m.Title, &m.Type, &m.Picture, &m.Track)
		if err != nil {
			log.Printf("Failed to retrieve row because: %s", err)
			continue
		}

		listMusic = append(listMusic, m)
	}

	return listMusic, nil
}

func (m Music) GetMusics(goFmDb *repository.DB) ([]Music, error) {
	res, err := goFmDb.GetDotSql().Query(goFmDb.GetSqlDb(), "find-all-music")
	var listMusic []Music

	if err != nil {
		return []Music{}, err
	}

	for res.Next() {
		err = res.Scan(&m.ID, &m.Name, &m.Title, &m.Type, &m.Picture, &m.Track)
		if err != nil {
			log.Printf("Failed to retrieve row because: %s", err)
			continue
		}

		listMusic = append(listMusic, m)
	}

	return listMusic, nil
}

func (m Music) Update(goFmDb *repository.DB, newData Music) error {
	if err := m.validateMusic(goFmDb, newData); err != nil {
		return err
	}

	_, err := goFmDb.GetDotSql().Exec(
		goFmDb.GetSqlDb(),
		"update-music",
		newData.Name,
		newData.Title,
		newData.Type,
		newData.Picture,
		newData.Track,
		m.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m Music) Delete(goFmDb *repository.DB) error {
	_, err := goFmDb.GetDotSql().Exec(goFmDb.GetSqlDb(), "delete-music", m.ID)

	if err != nil {
		return err
	}

	return nil
}

func (m Music) UpdateTrackAndPicture(goFmDb *repository.DB, picture, track string) error {
	_, err := goFmDb.GetDotSql().Exec(
		goFmDb.GetSqlDb(),
		"update-picture-track",
		picture,
		track,
		m.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
