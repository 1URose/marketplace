package app

import (
	"context"
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest/auth"
	"github.com/1URose/marketplace/internal/common/config"
	"github.com/1URose/marketplace/internal/common/db"
	"github.com/1URose/marketplace/internal/common/jwt"
	"github.com/gin-gonic/gin"
)

type Deps struct {
	Ctx            context.Context
	Engine         *gin.Engine
	DB             *db.Connections
	GeneralConfig  *config.GeneralConfig
	JWTManager     *jwt.Manager
	AuthMiddleware *auth.Middleware
}

func NewDeps(ctx context.Context, engine *gin.Engine, connections *db.Connections, generalCfg *config.GeneralConfig) *Deps {
	jwtMgr := jwt.NewManager(generalCfg.CommonConfig)

	authMiddleware := auth.NewMiddleware(generalCfg.CommonConfig.BearerPrefix, jwtMgr)

	return &Deps{
		Ctx:            ctx,
		Engine:         engine,
		DB:             connections,
		GeneralConfig:  generalCfg,
		JWTManager:     jwtMgr,
		AuthMiddleware: authMiddleware,
	}
}
