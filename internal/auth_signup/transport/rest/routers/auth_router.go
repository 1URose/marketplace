package routers

import (
	"context"
	"github.com/1URose/marketplace/internal/auth_signup/infrastructure/repository/redis"
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest/auth"
	"github.com/1URose/marketplace/internal/auth_signup/use_cases"
	"github.com/1URose/marketplace/internal/common/db"
	"github.com/1URose/marketplace/internal/user_profile/infrastructure/repository/postgresql"
	"github.com/gin-gonic/gin"
	"log"
)

type AuthRouter struct {
	engine      *gin.Engine
	ctx         context.Context
	Connections *db.Connections
}

func NewAuthRouter(ctx context.Context, engine *gin.Engine, connections *db.Connections) *AuthRouter {
	log.Println("[routers:auth] initializing AuthRouter")

	return &AuthRouter{
		ctx:         ctx,
		engine:      engine,
		Connections: connections,
	}
}

func initRedisServer(connections *db.Connections) *use_cases.AuthService {

	log.Println("[routers:auth] initializing AuthService with Redis and Postgres repositories")

	redisR := redis.NewRedisRepository(connections.RedisConn)

	userR := postgresql.NewUserRepository(connections.UserPostgresConn)

	svc := use_cases.NewAccountService(redisR, userR)

	log.Println("[routers:auth] AuthService initialized")

	return svc
}

func (ar *AuthRouter) RegisterRoutes() {

	log.Println("[routers:auth] registering /auth endpoints")

	apiGroup := ar.engine.Group("/auth")

	service := initRedisServer(ar.Connections)
	handler := auth.NewAuthHandler(service)

	{

		apiGroup.POST("/signup", handler.SignUp)

		apiGroup.POST("/login/", handler.Login)

		apiGroup.POST("/refresh", handler.Refresh)

	}

	log.Println("[routers:auth] all auth routers registered")

}
