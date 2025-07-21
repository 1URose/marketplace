package main

import (
	"context"
	"github.com/1URose/marketplace/internal/app"
	"github.com/1URose/marketplace/internal/common/logger"

	"log"
)

func main() {
	ctx := context.Background()
	logger.Init()

	log.Println("[cmd] Starting marketplace-service...")

	if err := app.Run(ctx); err != nil {
		log.Fatalf("[cmd] Application terminated with error: %v", err)
	}
}
