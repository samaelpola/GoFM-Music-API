package music

import (
	"encoding/json"
	"errors"
	"fmt"
	midllewares "github.com/samaelpola/GoFM-Music-API/internal/handlers/middlewares"
	"github.com/samaelpola/GoFM-Music-API/internal/models"
	"github.com/samaelpola/GoFM-Music-API/internal/repository"
	"github.com/samaelpola/GoFM-Music-API/internal/service"
	"github.com/samaelpola/GoFM-Music-API/internal/utils"
	"net/http"
)

type Delete struct {
	gofmDb *repository.DB
}

func NewDelete(musicDB *repository.DB) *Delete {
	return &Delete{
		gofmDb: musicDB,
	}
}

// Handle delete music
//
// @Summary delete music
// @Description	delete music
// @ID	delete-music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param musicID path int true "music ID"
// @Success 200	{object} utils.HttpResponse	"OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404	{string} string "Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /musics/{musicID} [delete]
// @Tags Musics
func (d Delete) Handle(w http.ResponseWriter, r *http.Request) {
	music := r.Context().Value(midllewares.MusicKey).(models.Music)
	if music.Picture != "" {
		errDeleteFile := service.DeleteFile(music.Picture)
		if errDeleteFile != nil && !errors.Is(service.FileNotFound, errDeleteFile) {
			utils.HttpErrorInternalError(
				w,
				fmt.Sprintf("unable to delete file to s3: %s", errDeleteFile),
				http.StatusInternalServerError,
			)
			return
		}
	}

	if music.Track != "" {
		errDeleteFile := service.DeleteFile(music.Track)
		if errDeleteFile != nil && !errors.Is(service.FileNotFound, errDeleteFile) {
			utils.HttpErrorInternalError(
				w,
				fmt.Sprintf("unable to delete file to s3: %s", errDeleteFile),
				http.StatusInternalServerError,
			)
			return
		}
	}

	if err := music.Delete(d.gofmDb); err != nil {
		fmt.Println("Error to delete music in db: ", err)
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("Error to delete music in db: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(map[string]string{"success": "music delete"})
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
			fmt.Sprintf("Error during write result: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
}
