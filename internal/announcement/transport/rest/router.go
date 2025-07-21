package rest

import (
	"context"
	routers "github.com/1URose/marketplace/internal/announcement/transport/rest/routes"
	"github.com/1URose/marketplace/internal/common/config/ad_limits"
	"github.com/1URose/marketplace/internal/common/db"
	"github.com/1URose/marketplace/internal/common/jwt"
	"github.com/gin-gonic/gin"
	"log"
)

func RegisterRoutes(ctx context.Context, engine *gin.Engine, connections *db.Connections, jwtMgr *jwt.Manager, cfg *ad_limits.AdConfig) {
	log.Println("[rest:announcement] registering announcement routers")

	adRoute := routers.NewAdRoute(ctx, engine, connections.PostgresConn, jwtMgr, cfg)

	adRoute.RegisterRoutes()

	log.Println("[rest:announcement] announcement routers registered successfully")
}
