package app

import (
	"context"
	"github.com/1URose/marketplace/internal/common/db"
	"github.com/1URose/marketplace/internal/user_profile/transport/rest"

	"github.com/gin-gonic/gin"
	"log"
)

func Run(ctx context.Context, engine *gin.Engine, connections *db.Connections) {
	log.Println("[user_profile] registering routers")

	rest.RegisterRoutes(ctx, engine, connections)

	log.Println("[user_profile] routers registered successfully")
}
