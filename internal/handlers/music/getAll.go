package music

import (
	"encoding/json"
	"fmt"
	"github.com/samaelpola/GoFM-Music-API/internal/models"
	"github.com/samaelpola/GoFM-Music-API/internal/repository"
	"github.com/samaelpola/GoFM-Music-API/internal/utils"
	"net/http"
)

type GetAllMusic struct {
	gofmDb *repository.DB
}

func NewGetAllMusic(gofmDB *repository.DB) *GetAllMusic {
	return &GetAllMusic{
		gofmDb: gofmDB,
	}
}

// Handle get all music
//
// @Summary get all music
// @Description	get all music
// @ID	get-all-music
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200	{object} []models.Music	"OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404	{string} string	"Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /musics [get]
// @Tags Musics
func (g GetAllMusic) Handle(w http.ResponseWriter, r *http.Request) {
	var music models.Music
	listMusic, err := music.GetMusics(g.gofmDb)
	if err != nil {
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("Error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	response, err := json.Marshal(listMusic)
	if err != nil {
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("Error during parse list of music: %s", err),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(response); err != nil {
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("Error: %s", err),
			http.StatusInternalServerError,
		)
		return
	}
}
