package main

import (
	"context"
	"github.com/1URose/marketplace/internal/app"
	"github.com/1URose/marketplace/internal/common/logger"

	"log"
)

// TODO: добавить тесты
// TODO: добавить валидацию особенно на обновление
// TODO: добавить единую обработку ошибок

func main() {
	ctx := context.Background()
	logger.Init()

	log.Println("Starting user-service")

	if err := app.Run(ctx); err != nil {
		log.Fatalf("Application terminated with error: %v", err)
	}
}
