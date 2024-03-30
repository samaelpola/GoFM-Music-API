package music

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/samaelpola/GoFM-Music-API/internal/models"
	"github.com/samaelpola/GoFM-Music-API/internal/repository"
	"github.com/samaelpola/GoFM-Music-API/internal/service"
	"github.com/samaelpola/GoFM-Music-API/internal/utils"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

type Create struct {
	gofmDb *repository.DB
}

func NewCreate(musicDB *repository.DB) *Create {
	return &Create{
		gofmDb: musicDB,
	}
}

// Handle create music
//
// @Summary Create a new music
// @Description Create a new music with the provided details
// @ID create-music
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param name formData string true "Name of the artist"
// @Param title formData string true "Title of the music"
// @Param type formData string true "Type of music (e.g., GO-ROCK, GO-POP, GO-RAP, GO-SLOW)"
// @Param picture formData file true "Image file for the music"
// @Param audio formData file true "Audio file for the music"
// @Success 201 {object} utils.HttpResponse "Created"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 409 {string} string "Conflict"
// @Failure 415 {string} string "Unsupported Media Type"
// @Failure 500 {string} string "Internal server error"
// @Router /musics [post]
// @Tags Musics
func (c Create) Handle(w http.ResponseWriter, r *http.Request) {
	newMusic := utils.ParseFormMusic(r)

	if !utils.CheckTypeOfMusicExist(newMusic.Type) {
		http.Error(w, fmt.Sprintf("type '%s' do not exist", newMusic.Type), http.StatusNotFound)
	}

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

	if !strings.Contains(pictureHeader.Header.Get("Content-Type"), "image/") {
		http.Error(w, "Invalid image file type", http.StatusUnsupportedMediaType)
		return
	}

	musicCreated, err := newMusic.Create(c.gofmDb)
	if err != nil {
		if errors.Is(err, models.ErrMusicAlreadyExist) {
			http.Error(
				w,
				fmt.Sprintf("music with name '%s' and title '%s' already exists", newMusic.Name, newMusic.Title),
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

	pictureKey := fmt.Sprintf(
		"/%s/%d/image/%d%s",
		musicCreated.Type,
		musicCreated.ID,
		musicCreated.ID,
		filepath.Ext(pictureHeader.Filename),
	)

	audioKey := fmt.Sprintf(
		"/%s/%d/audio/%d.mp3",
		musicCreated.Type,
		musicCreated.ID,
		musicCreated.ID,
	)

	var wg sync.WaitGroup
	var errUploadFile error
	wg.Add(2)
	go func() {
		defer wg.Done()
		errUploadFile = service.UploadFile(pictureHeader, pictureKey)
	}()
	go func() {
		defer wg.Done()
		errUploadFile = service.UploadFile(audioHeader, audioKey)
	}()
	wg.Wait()

	if errUploadFile != nil {
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("unable to upload file to s3: %s", errUploadFile),
			http.StatusInternalServerError,
		)
		return
	}

	if err = musicCreated.UpdateTrackAndPicture(
		c.gofmDb,
		pictureKey,
		audioKey,
	); err != nil {
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("Error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(map[string]string{"success": strconv.Itoa(musicCreated.ID)})
	if err != nil {
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("Error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	if _, err = w.Write(response); err != nil {
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("Error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
}
