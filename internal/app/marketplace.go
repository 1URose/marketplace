// Package app
// @title Marketplace API
// @version 1.0
// @description API реализующее работу с пользователями и объявлениями
// @host localhost:8080
// @BasePath /
// @schemes http
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Тип "Bearer" следует указать в качестве значения этого заголовка
package app

import (
	"context"
	"fmt"
	"github.com/1URose/marketplace/docs"
	authApp "github.com/1URose/marketplace/internal/auth_signup/app"
	"github.com/1URose/marketplace/internal/common/db"
	"github.com/1URose/marketplace/internal/common/settings"
	userApp "github.com/1URose/marketplace/internal/user_profile/app"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"

	"log"
)

func Run(ctx context.Context) error {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(loggingMiddleware)

	engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		MaxAge:           1 * time.Hour,
	}))

	docs.SwaggerInfo.BasePath = ""
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	if err := settings.LoadEnv(".env"); err != nil {

		log.Printf("failed to load env: %v", err)

		return err
	}

	log.Println("Environment variables loaded")

	connections, err := db.NewConnections()

	if err != nil {

		log.Printf("failed to establish connections: %v", err)

		return err
	}

	defer connections.Close()

	log.Println("Database connections established")

	userApp.Run(ctx, engine, connections)
	authApp.Run(ctx, engine, connections)

	addr := ":8080"
	swaggerURL := fmt.Sprintf("http://localhost%s/swagger/index.html", addr)
	log.Printf("Swagger UI available at %s", swaggerURL)

	if err = engine.Run(addr); err != nil {

		log.Printf("Error starting the application: %v", err)

		return err
	}

	return nil
}

func loggingMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()

	latency := time.Since(start)
	status := c.Writer.Status()
	method := c.Request.Method
	path := c.Request.URL.Path
	ip := c.ClientIP()

	if len(c.Errors) > 0 {

		log.Printf("ERROR: status=%d method=%s path=%s ip=%s latency=%v errors=%s",
			status, method, path, ip, latency, c.Errors.String())
	} else {

		log.Printf("INFO: status=%d method=%s path=%s ip=%s latency=%v",
			status, method, path, ip, latency)
	}
}
