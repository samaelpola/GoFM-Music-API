package midllewares

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/samaelpola/GoFM-Music-API/internal/models"
	"github.com/samaelpola/GoFM-Music-API/internal/repository"
	"net/http"
	"strconv"
)

func CheckMusicExist(gofmDb *repository.DB) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			id, err := strconv.Atoi(vars["musicID"])
			if err != nil {
				http.Error(w, fmt.Sprintf("Invalid music ID: '{%d}'", id), http.StatusBadRequest)
				return
			}

			var music models.Music
			currentMusic, err := music.GetMusic(gofmDb, int64(id))
			if err != nil {
				http.Error(w, fmt.Sprintf("Music for id '%d' not found", id), http.StatusNotFound)
				return
			}

			ctx := context.WithValue(r.Context(), "music", currentMusic)
			h.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
