package main

import (
	"url_shortener/cmd/internal/config"
)

func main() {
	cfg := config.MustLoad()

	// TODO : init logger: slog

	// TODO: init storage: postgres

	// TODO: init router: chi

	// TODO: run server
}
