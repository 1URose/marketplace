package rest

import (
	"context"
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest/routers"
	"github.com/1URose/marketplace/internal/common/db"
	"github.com/gin-gonic/gin"

	"log"
)

func RegisterRoutes(ctx context.Context, engine *gin.Engine, connections *db.Connections) {
	log.Println("[rest:auth_signup] registering auth routers")

	authRouter := routers.NewAuthRouter(ctx, engine, connections)
	authRouter.RegisterRoutes()

	log.Println("[rest:auth_signup] auth routers registered successfully")
}
