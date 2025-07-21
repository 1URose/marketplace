package app

import (
	"context"
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest"
	"github.com/1URose/marketplace/internal/common/db"
	"github.com/1URose/marketplace/internal/common/jwt"
	"github.com/gin-gonic/gin"
	"log"
)

func Run(ctx context.Context, engine *gin.Engine, connections *db.Connections, jwtMgr *jwt.Manager) {
	log.Println("[auth_signup] registering auth routers")

	rest.RegisterRoutes(ctx, engine, connections, jwtMgr)

	log.Println("[auth_signup] auth routers registered successfully")
}
