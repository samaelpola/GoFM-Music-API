package music

import (
	"errors"
	"fmt"
	midllewares "github.com/samaelpola/GoFM-Music-API/internal/handlers/middlewares"
	"github.com/samaelpola/GoFM-Music-API/internal/models"
	"github.com/samaelpola/GoFM-Music-API/internal/service"
	"github.com/samaelpola/GoFM-Music-API/internal/utils"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

// GetAudioFileHandle get audio file of a music.
//
// @Summary Get audio file of music
// @Description Retrieve the audio file of a music.
// @ID get-audio-file-of-music
// @Accept json
// @Produce octet-stream
// @Security BearerAuth
// @Param musicID path int true "Music ID"
// @Success 200 {file} file "OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404	{string} string	"Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /musics/{musicID}/audio [get]
// @Tags Musics
func GetAudioFileHandle(w http.ResponseWriter, r *http.Request) {
	music := r.Context().Value(midllewares.MusicKey).(models.Music)
	if music.Track == "" {
		utils.HttpError(
			w,
			fmt.Sprintf("Audio file of music by id %d not found", music.ID),
			http.StatusNotFound,
		)
		return
	}

	file, err := service.DownloadFile(music.Track)
	if err != nil {
		if errors.Is(err, service.FileNotFound) {
			utils.HttpError(
				w,
				fmt.Sprintf("audio file for music with id %d not found", music.ID),
				http.StatusNotFound,
			)
			return
		}

		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("unable to download file to s3: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	info, err := file.Stat()
	if err != nil {
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("error during get info of file: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(
			`attachment; 
				filename="%s_%s%s"`, music.Name, music.Title, filepath.Ext(music.Track),
		),
	)
	w.Header().Set("Content-Length", strconv.Itoa(int(info.Size)))
	w.WriteHeader(http.StatusOK)

	_, err = io.Copy(w, file)
	if err != nil {
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("Error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
}
