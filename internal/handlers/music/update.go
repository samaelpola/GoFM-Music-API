package music

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samaelpola/GoFM-Music-API/internal/models"
	"github.com/samaelpola/GoFM-Music-API/internal/repository"
	"github.com/samaelpola/GoFM-Music-API/internal/utils"
	"net/http"
	"strings"
)

type Update struct {
	gofmDb *repository.DB
}

func NewUpdate(musicDB *repository.DB) *Update {
	return &Update{
		gofmDb: musicDB,
	}
}

// Handle update music
//
// @Summary update music
// @Description update a music
// @ID update-music
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Name of the artist"
// @Param title formData string true "Title of the music"
// @Param type formData string true "Type of music (e.g., GO-ROCK, GO-POP, GO-RAP, GO-SLOW)"
// @Param picture formData file true "Image file for the music"
// @Param audio formData file true "Audio file for the music"
// @Success 200 {object} utils.HttpResponse "OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 409 {string} string "Conflict"
// @Failure 415 {string} string "Unsupported Media Type"
// @Failure 500 {string} string "Internal server error"
// @Router /musics/{musicID} [put]
// @Tags Musics
func (u Update) Handle(w http.ResponseWriter, r *http.Request) {
	music := r.Context().Value("music").(models.Music)
	musicData := utils.ParseFormMusic(r)

	_, audioHeader, err := r.FormFile("audio")
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if audioHeader.Header.Get("Content-Type") != "audio/mp3" &&
		audioHeader.Header.Get("Content-Type") != "audio/mpeg" {
		http.Error(w, "Invalid audio file type mp3 required", http.StatusUnsupportedMediaType)
		return
	}

	_, pictureHeader, err := r.FormFile("picture")
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if strings.Contains(pictureHeader.Header.Get("Content-Type"), "image/") != true {
		http.Error(w, "Invalid image file type", http.StatusUnsupportedMediaType)
		return
	}

	if err := music.Update(u.gofmDb, musicData); err != nil {
		if err != nil {
			if errors.Is(err, models.ErrMusicAlreadyExist) {
				http.Error(
					w,
					fmt.Sprintf("music with name '%s' and title '%s' already exists", musicData.Name, musicData.Title),
					http.StatusConflict,
				)
				return
			}

			utils.HttpErrorInternalError(
				w,
				fmt.Sprintf("Error: %s", err),
				http.StatusInternalServerError,
			)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(map[string]string{"success": "music update"})
	if err != nil {
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("Error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	if _, err := w.Write(response); err != nil {
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("Error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
}
