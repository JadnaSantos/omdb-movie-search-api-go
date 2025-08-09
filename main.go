package main

import (
	"github/JadnaSantos/omdb-movie-search-api-go.git/api"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func goDotEnvVariable(key string) string {

  err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

  return os.Getenv(key)
}


func main () { 
	if err := run(); err != nil {
		slog.Error("failed to execute code", "error", err) 
		os.Exit(1) 
	}

	slog.Info("all systems offline")
}


func run () error {
	apiKey := os.Getenv("OMDB_KEY")
	handle := api.NewHandler(apiKey)

	s := http.Server{
		Addr:                         ":8080",
		Handler:                      handle,	
		ReadTimeout:                  10 * time.Second,
		WriteTimeout:                 10 * time.Second,
		IdleTimeout:                  time.Minute,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil
}