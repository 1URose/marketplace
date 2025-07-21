package routers

import (
	"context"
	"github.com/1URose/marketplace/internal/announcement/infrastructure/repository/postgresql"
	"github.com/1URose/marketplace/internal/announcement/transport/rest/ad"
	"github.com/1URose/marketplace/internal/announcement/use_cases"
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest/auth"
	"github.com/1URose/marketplace/internal/common/app"
	"github.com/1URose/marketplace/internal/common/config"
	"github.com/1URose/marketplace/internal/common/jwt"
	"github.com/1URose/marketplace/internal/common/validator"

	pgConfig "github.com/1URose/marketplace/internal/common/db/postgresql"

	"github.com/gin-gonic/gin"
	"log"
)

type AdRoute struct {
	ctx            context.Context
	engine         *gin.Engine
	pgClient       *pgConfig.Client
	jwtMgr         *jwt.Manager
	cfg            *config.GeneralConfig
	authMiddleware *auth.Middleware
}

func NewAdRoute(deps *app.Deps) *AdRoute {
	log.Println("[routers:ad] initializing AdRoute")
	return &AdRoute{
		ctx:            deps.Ctx,
		engine:         deps.Engine,
		pgClient:       deps.DB.PostgresConn,
		jwtMgr:         deps.JWTManager,
		cfg:            deps.GeneralConfig,
		authMiddleware: deps.AuthMiddleware,
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

	service := initAdService(ar.pgClient, ar.cfg.AdConfig.PageSize)

	v := validator.NewAllowedValues(ar.cfg.AdConfig)

	handler := ad.NewHandler(service, v)

	privateApiGroup := ar.engine.Group("/ad").Use(ar.authMiddleware.Require())

	{
		privateApiGroup.POST("/", handler.CreateAd)
		log.Println("[routers:ad] registered POST /ad/")
	}

	publicApiGroup := ar.engine.Group("/ads").Use(ar.authMiddleware.Optional())
	{
		publicApiGroup.GET("/", handler.GetAllAds)
		log.Println("[routers:ad] registered GET /ad/")

	}

	log.Println("[routers:ad] /ad endpoints registered successfully")
}
