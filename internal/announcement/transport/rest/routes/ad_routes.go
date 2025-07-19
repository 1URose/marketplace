package routers

import (
	"context"
	"github.com/1URose/marketplace/internal/announcement/infrastructure/repository/postgresql"
	"github.com/1URose/marketplace/internal/announcement/transport/rest/ad"
	"github.com/1URose/marketplace/internal/announcement/use_cases"
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest/auth"

	pgConfig "github.com/1URose/marketplace/internal/common/db/postgresql"

	"github.com/gin-gonic/gin"
	"log"
)

type AdRoute struct {
	PGClient *pgConfig.Client
	engine   *gin.Engine
	ctx      context.Context
}

func NewAdRoute(ctx context.Context, engine *gin.Engine, PGClient *pgConfig.Client) *AdRoute {
	return &AdRoute{
		ctx:      ctx,
		engine:   engine,
		PGClient: PGClient,
	}
}

func initAdService(PGClient *pgConfig.Client) *use_cases.AdService {
	log.Println("[routers:ad] initializing AdService")

	repo := postgresql.NewAdRepository(PGClient)

	service := use_cases.NewAdService(repo)

	log.Println("[routers:ad] AdService initialized")

	return service

}

func (ar *AdRoute) RegisterRoutes() {
	log.Println("[routers:ad] registering /ad endpoints")

	service := initAdService(ar.PGClient)

	handler := ad.NewHandler(service)

	privateApiGroup := ar.engine.Group("/ad").Use(auth.RequireAuthMiddleware())

	{
		privateApiGroup.POST("/", handler.CreateAd)
		log.Println("[routers:ad] registered POST /ad/")
	}

	publicApiGroup := ar.engine.Group("/ad").Use(auth.OptionalAuthMiddleware())
	{
		publicApiGroup.GET("/", handler.GetAllAds)
		log.Println("[routers:ad] registered GET /ad/")

		publicApiGroup.GET("/:id", handler.GetAdByID)
		log.Println("[routers:ad] registered GET /ad/:id")
	}
}
