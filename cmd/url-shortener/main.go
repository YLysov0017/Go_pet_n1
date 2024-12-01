package main

import (
	"fmt"

	"github.com/YLysov0017/go_pet_n1/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)

	// TODO: init logger: slog

	// TODO: init storage: sqlite

	// TODO: init router: chi, chi render

	// TODO: run server
}
