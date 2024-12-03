package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/YLysov0017/go_pet_n1/internal/config"
	"github.com/YLysov0017/go_pet_n1/internal/config/storage/sqlite"
	"github.com/YLysov0017/go_pet_n1/internal/http-server/middleware/mwlogger"
	"github.com/YLysov0017/go_pet_n1/internal/http-server/middleware/mwlogger/handlers/url/save"
	"github.com/YLysov0017/go_pet_n1/internal/lib/logger/handlers/slogpretty"
	"github.com/YLysov0017/go_pet_n1/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("service started", slog.String("Env", cfg.Env))
	log.Debug("debug enabled")

	var storPath string
	switch runtime.GOOS {
	case "windows":
		mydir, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
		}
		storPath = strings.ReplaceAll(strings.Replace(cfg.StoragePath, ".", "", 1), "/", "\\") // Fixin' Windows filepath
		storPath = filepath.Join(mydir, "..\\..", storPath)
	default:
		storPath = cfg.StoragePath // Unix path remain unchanged
	}

	fmt.Println(cfg)
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwlogger.New(log))
	router.Use(middleware.URLFormat)
	router.Use(middleware.Recoverer)

	storage, err := sqlite.New(storPath)

	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	router.Post("/url", save.New(log, storage, cfg.AliasLength))

	log.Info("starting server on ", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed  to start server")
	}
	log.Error("server stopped")
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger { // берем настройки из окружения
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log

}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
