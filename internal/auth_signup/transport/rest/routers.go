package rest

import (
	"context"
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest/routers"
	"github.com/1URose/marketplace/internal/common/db"
	"github.com/1URose/marketplace/internal/common/jwt"
	"github.com/gin-gonic/gin"

	"log"
)

func RegisterRoutes(ctx context.Context, engine *gin.Engine, connections *db.Connections, jwtMgr *jwt.Manager) {
	log.Println("[rest:auth_signup] registering auth routers")

	authRouter := routers.NewAuthRouter(ctx, engine, connections, jwtMgr)
	authRouter.RegisterRoutes()

	log.Println("[rest:auth_signup] auth routers registered successfully")
}
