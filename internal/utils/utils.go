package utils

import (
	"github.com/samaelpola/GoFM-Music-API/internal/models"
	"log"
	"net/http"
	"slices"
	"strings"
)

func HttpError(w http.ResponseWriter, err string, status int) {
	log.Println(err)
	http.Error(w, err, status)
}

func HttpErrorInternalError(w http.ResponseWriter, err string, status int) {
	log.Println(err)
	http.Error(w, "internal server error", status)
}

func CheckTypeOfMusicExist(musicType string) bool {
	return slices.Contains([]string{models.GOPOP, models.GORAP, models.GOROCK, models.GOSLOW}, musicType)
}

func ParseFormMusic(r *http.Request) models.Music {
	var musicData models.Music
	musicData.Name = r.FormValue("name")
	musicData.Title = r.FormValue("title")
	musicData.Type = strings.ToUpper(r.FormValue("type"))
	musicData.Track = ""
	musicData.Picture = ""

	return musicData
}

type HttpResponse struct {
	Message string `json:"message"`
}
