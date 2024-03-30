package music

import (
	"errors"
	"fmt"
	"github.com/samaelpola/GoFM-Music-API/internal/models"
	"github.com/samaelpola/GoFM-Music-API/internal/service"
	"github.com/samaelpola/GoFM-Music-API/internal/utils"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
)

// GetImageHandle get the image of a music.
//
// @Summary Get image of music
// @Description Retrieve the image of a music by ID.
// @ID get-image-of-music
// @Accept json
// @Produce octet-stream
// @Security BearerAuth
// @Param musicID path int true "Music ID"
// @Success 200 {file} file "OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404	{string} string	"Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /musics/{musicID}/image [get]
// @Tags Musics
func GetImageHandle(w http.ResponseWriter, r *http.Request) {
	music := r.Context().Value("music").(models.Music)
	if music.Picture == "" {
		http.Error(
			w,
			fmt.Sprintf("Image for id %d not found", music.ID),
			http.StatusNotFound,
		)
		return
	}

	file, err := service.DownloadFile(music.Picture)
	if err != nil {
		if errors.Is(err, service.FileNotFound) {
			utils.HttpError(
				w,
				fmt.Sprintf("Image for id %d not found", music.ID),
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

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Length", strconv.Itoa(int(info.Size)))
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(
			`attachment; 
				filename="%s_%s%s"`, music.Name, music.Title, filepath.Ext(music.Picture),
		),
	)

	_, err = io.Copy(w, file)
	if err != nil {
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("error during rendering music: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
}
