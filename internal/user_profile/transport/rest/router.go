package rest

import (
	"context"
	"github.com/1URose/marketplace/internal/common/db"
	"github.com/1URose/marketplace/internal/user_profile/transport/rest/routes"
	"github.com/gin-gonic/gin"
	"log"
)

func RegisterRoutes(ctx context.Context, engine *gin.Engine, connections *db.Connections) {
	log.Println("[rest:user_profile] registering user_profile routers")

	userRoute := routes.NewUserRoute(ctx, engine, connections.UserPostgresConn)

	userRoute.RegisterRoutes()

	log.Println("[rest:user_profile] user_profile routers registered")
}
