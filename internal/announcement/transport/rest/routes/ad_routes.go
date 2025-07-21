package routers

import (
	"context"
	"github.com/1URose/marketplace/internal/announcement/infrastructure/repository/postgresql"
	"github.com/1URose/marketplace/internal/announcement/transport/rest/ad"
	"github.com/1URose/marketplace/internal/announcement/use_cases"
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest/auth"
	"github.com/1URose/marketplace/internal/common/config/ad_limits"
	"github.com/1URose/marketplace/internal/common/jwt"
	"github.com/1URose/marketplace/internal/common/validator"

	pgConfig "github.com/1URose/marketplace/internal/common/db/postgresql"

	"github.com/gin-gonic/gin"
	"log"
)

type AdRoute struct {
	ctx      context.Context
	engine   *gin.Engine
	pgClient *pgConfig.Client
	jwtMgr   *jwt.Manager
	cfg      *ad_limits.AdConfig
}

func NewAdRoute(ctx context.Context, engine *gin.Engine, PGClient *pgConfig.Client, jwtMgr *jwt.Manager, cfg *ad_limits.AdConfig) *AdRoute {
	log.Println("[routers:ad] initializing AdRoute")
	return &AdRoute{
		ctx:      ctx,
		engine:   engine,
		pgClient: PGClient,
		jwtMgr:   jwtMgr,
		cfg:      cfg,
	}
}

func initAdService(PGClient *pgConfig.Client, pageSize int) *use_cases.AdService {
	log.Println("[routers:ad] initializing AdService")

	repo := postgresql.NewAdRepository(PGClient)

	service := use_cases.NewAdService(repo, pageSize)

	log.Println("[routers:ad] AdService initialized")

	return service

}

func (ar *AdRoute) RegisterRoutes() {
	log.Println("[routers:ad] registering /ad endpoints")

	service := initAdService(ar.pgClient, ar.cfg.PageSize)

	v := validator.NewAllowedValues(ar.cfg)

	handler := ad.NewHandler(service, v)

	privateApiGroup := ar.engine.Group("/ad").Use(auth.RequireAuthMiddleware(ar.jwtMgr))

	{
		privateApiGroup.POST("/", handler.CreateAd)
		log.Println("[routers:ad] registered POST /ad/")
	}

	publicApiGroup := ar.engine.Group("/ads").Use(auth.OptionalAuthMiddleware(ar.jwtMgr))
	{
		publicApiGroup.GET("/", handler.GetAllAds)
		log.Println("[routers:ad] registered GET /ad/")

	}

	log.Println("[routers:ad] /ad endpoints registered successfully")
}
