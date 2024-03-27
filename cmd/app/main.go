package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	midllewares "github.com/samaelpola/GoFM-Music-API/internal/handlers/middlewares"
	"github.com/samaelpola/GoFM-Music-API/internal/handlers/music"
	"github.com/samaelpola/GoFM-Music-API/internal/repository"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
	"time"

	_ "github.com/samaelpola/GoFM-Music-API/docs"
)

// @title Swagger GO-FM Music  API
// @version 1.0
// @description Api go fm.
// @contact.name   Samael POLA
// @contact.email  me@lan.lan
// @SecurityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Enter the token with the `Bearer ` prefix, e.g. "Bearer lnez564".
// @host localhost:8083
func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	router := mux.NewRouter()

	myDb, err := repository.Initialize("sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer myDb.GetSqlDb().Close()

	createMusic := music.NewCreate(myDb)
	updateMusic := music.NewUpdate(myDb)
	deleteMusic := music.NewDelete(myDb)
	getMusics := music.NewGetAllMusic(myDb)
	getMusicByType := music.NewGetByType(myDb)

	musicRouterGlobal := router.PathPrefix("/musics").Subrouter()
	musicRouterGlobal.Use(midllewares.Authenticate)
	musicRouterGlobal.HandleFunc("", createMusic.Handle).Methods(http.MethodPost)
	musicRouterGlobal.HandleFunc("", getMusics.Handle).Methods(http.MethodGet)
	musicRouterGlobal.HandleFunc("/{musicType:[A-Za-z]+-[A-Za-z]+$}", getMusicByType.Handle)

	musicRouter := musicRouterGlobal.PathPrefix("/{musicID:[0-9]+$}").Subrouter()
	musicRouter.Use(midllewares.CheckMusicExist(myDb))
	musicRouter.HandleFunc("", music.GetMusicHandle).Methods(http.MethodGet)
	musicRouter.HandleFunc("", updateMusic.Handle).Methods(http.MethodPut)
	musicRouter.HandleFunc("", deleteMusic.Handle).Methods(http.MethodDelete)

	audioFileRouter := musicRouterGlobal.PathPrefix("/{musicID}").Subrouter()
	audioFileRouter.Use(midllewares.CheckMusicExist(myDb))
	audioFileRouter.HandleFunc("/audio", music.GetAudioFileHandle)

	musicImageRouter := router.PathPrefix("/musics/{musicID}").Subrouter()
	musicImageRouter.Use(midllewares.CheckMusicExist(myDb))
	musicImageRouter.HandleFunc("/image", music.GetImageHandle)

	router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":8083",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Println("listening on 8083 ....")
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
