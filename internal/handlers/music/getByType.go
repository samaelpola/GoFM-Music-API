package music

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/samaelpola/GoFM-Music-API/internal/models"
	"github.com/samaelpola/GoFM-Music-API/internal/repository"
	"github.com/samaelpola/GoFM-Music-API/internal/utils"
	"net/http"
	"strings"
)

type GetByType struct {
	gofmDb *repository.DB
}

func NewGetByType(musicDB *repository.DB) *GetByType {
	return &GetByType{
		gofmDb: musicDB,
	}
}

// Handle get music by type
//
// @Summary get music by Type
// @Description	get music by Type
// @ID	get-music-by-type
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param musicType path string true "music Type"
// @Success 200	{object} []models.Music	"OK"
// @Failure 401 {string} string "Unauthorized"
// @Failure 403 {string} string "Forbidden"
// @Failure 404	{string} string	"Not found"
// @Failure 500 {string} string "Internal server error"
// @Router /musics/{musicType} [get]
// @Tags Musics
func (g GetByType) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	musicType := vars["musicType"]

	if !utils.CheckTypeOfMusicExist(strings.ToUpper(musicType)) {
		http.Error(w, fmt.Sprintf("type '%s' do not exist", musicType), http.StatusNotFound)
	}

	var music models.Music
	listMusic, err := music.GetMusicByType(g.gofmDb, strings.ToUpper(musicType))
	if err != nil {
		utils.HttpErrorInternalError(
			w,
			fmt.Sprintf("unable to get music by type: %s", err),
			http.StatusInternalServerError,
		)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	response, err := json.Marshal(listMusic)
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
