package music

import (
	"encoding/json"
	"fmt"
	midllewares "github.com/samaelpola/GoFM-Music-API/internal/handlers/middlewares"
	"github.com/samaelpola/GoFM-Music-API/internal/models"
	"github.com/samaelpola/GoFM-Music-API/internal/utils"
	"net/http"
)

// GetMusicHandle get music by id
//
// @Summary get music by ID
// @Description	get music by ID
// @ID	get-music-by-id
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param musicID path int true "music ID"
// @Success 200	{object} models.Music	"OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404	{string} string	"Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /musics/{musicID} [get]
// @Tags Musics
func GetMusicHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	music := r.Context().Value(midllewares.MusicKey).(models.Music)

	response, err := json.Marshal(music)
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
