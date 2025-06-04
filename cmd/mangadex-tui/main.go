package main

import (
	"log"

	"mangadex-tui/internal/ui"
)

func main() {
	app := ui.NewApp()
	if err := app.Run(); err != nil {
		log.Fatalf("Application error: %v", err)
	}
}
