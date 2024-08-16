package main

import (
	"URL-Shortener/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// TODO: init logger: slog

	// TODO: init storage: SQLite

	// TODO: init router: chi, render

	// TODO: run server
}
