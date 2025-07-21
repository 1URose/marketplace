package app

import (
	"context"
	"github.com/1URose/marketplace/internal/announcement/transport/rest"
	"github.com/1URose/marketplace/internal/common/config/ad_limits"
	"github.com/1URose/marketplace/internal/common/db"
	"github.com/1URose/marketplace/internal/common/jwt"
	"github.com/gin-gonic/gin"
	"log"
)

func Run(ctx context.Context, engine *gin.Engine, connections *db.Connections, jwtMgr *jwt.Manager, cfg *ad_limits.AdConfig) {
	log.Println("[announcement] registering routers")

	rest.RegisterRoutes(ctx, engine, connections, jwtMgr, cfg)
	log.Println("[announcement] routers registered successfully")
}
