// Package app
// @title Marketplace API
// @version 1.0
// @description API реализующее работу с пользователями и объявлениями
// @host localhost:8080
// @BasePath /
// @schemes http
package app

import (
	"context"
	"fmt"
	"github.com/1URose/marketplace/docs"
	adApp "github.com/1URose/marketplace/internal/announcement/app"
	authApp "github.com/1URose/marketplace/internal/auth_signup/app"
	"github.com/1URose/marketplace/internal/common/config"
	"github.com/1URose/marketplace/internal/common/db"
	"github.com/1URose/marketplace/internal/common/jwt"
	userApp "github.com/1URose/marketplace/internal/user_profile/app"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"log"
)

func Run(ctx context.Context) error {
	engine := initializeGin()

	generalConfig, err := config.LoadGeneralConfigFrom(".env")

	if err != nil {
		log.Printf("failed to load env: %v", err)
		return err
	}

	log.Println("Environment variables loaded")

	connections, err := db.NewConnections(generalConfig)
	defer connections.Close()
	if err != nil {
		log.Printf("failed to establish connections: %v", err)
		return err
	}

	log.Println("Database connections established")

	jwtManager := jwt.NewManager(generalConfig.CommonConfig)

	userApp.Run(ctx, engine, connections)
	authApp.Run(ctx, engine, connections, jwtManager)
	adApp.Run(ctx, engine, connections, jwtManager, generalConfig.AdConfig)

	addr := generalConfig.CommonConfig.GinAddress
	swaggerURL := fmt.Sprintf("http://localhost%s/swagger/index.html", addr)
	log.Printf("Swagger UI available at %s", swaggerURL)

	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s", err)
		}
	}()
	log.Printf("Server is running at %s", addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("Server Shutdown: %v", err)
	}

	log.Println("Server exiting")
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

func initializeGin() *gin.Engine {
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
	return engine
}
