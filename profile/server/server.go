package server

import (
	"context"
	"database/sql"
	"date-bot-go/profile"
	"date-bot-go/profile/repository"
	http2 "date-bot-go/profile/server/http"
	"date-bot-go/profile/service"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

func initDb() *sql.DB {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	return db
}

type App struct {
	profileService profile.Service
	server         *http.Server
}

func NewApp(server *http.Server) *App {
	db := initDb()
	profileRepo := repository.NewPostgresProfileRepository(db)
	profileService := service.NewProfileService(profileRepo)
	//todo grpc
	return &App{
		profileService: profileService,
		server:         server,
	}
}

func (a *App) Run(port string) error {
	profileHandler := http2.NewHandler(a.profileService)
	router := mux.NewRouter()
	router.HandleFunc("/api/profiles/", profileHandler.GetAll).Methods("GET")
	router.HandleFunc("/api/profiles/{id}/", profileHandler.GetById).Methods("GET")
	router.HandleFunc("/api/profiles/", profileHandler.Create).Methods("POST")
	router.HandleFunc("/api/profiles/{id}/", profileHandler.Update).Methods("PUT")
	router.HandleFunc("/api/profiles/{id}/", profileHandler.Delete).Methods("DELETE")
	//router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	a.server = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := a.server.ListenAndServe(); err != nil {
			log.Fatalf("Server error: %s", err.Error())
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()
	return a.server.Shutdown(ctx)

}
