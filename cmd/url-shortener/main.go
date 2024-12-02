package main

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/YLysov0017/go_pet_n1/internal/config"
	"github.com/YLysov0017/go_pet_n1/internal/config/storage/sqlite"
	"github.com/YLysov0017/go_pet_n1/internal/lib/logger/sl"
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

	storage, err := sqlite.New(storPath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	// TODO: init router: chi, chi render

	// TODO: run server
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
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log

}
