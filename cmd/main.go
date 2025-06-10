package main

import (
	"fmt"
	"log"
	"os"

	"github.com/sangnt1552314/mangadex-tui/internal/ui"
)

func main() {
	if err := os.MkdirAll("storage/logs", 0755); err != nil {
		panic(fmt.Errorf("failed to create logs directory: %w", err))
	}

	// Setup logging
	logFile, err := os.OpenFile("storage/logs/develop.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	app := ui.NewApp()

	if err := app.Run(); err != nil {
		panic(fmt.Errorf("failed to run application: %w", err))
	}
}
