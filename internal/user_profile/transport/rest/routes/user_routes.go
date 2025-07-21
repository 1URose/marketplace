package routes

import (
	"context"
	pgConfig "github.com/1URose/marketplace/internal/common/db/postgresql"
	pgRepo "github.com/1URose/marketplace/internal/user_profile/infrastructure/repository/postgresql"
	"github.com/1URose/marketplace/internal/user_profile/transport/rest/user"
	"github.com/1URose/marketplace/internal/user_profile/use_cases"
	"github.com/gin-gonic/gin"

	"log"
)

type UserRoute struct {
	PGClient *pgConfig.Client
	engine   *gin.Engine
	ctx      context.Context
}

func NewUserRoute(ctx context.Context, engine *gin.Engine, PGClient *pgConfig.Client) *UserRoute {
	return &UserRoute{
		ctx:      ctx,
		engine:   engine,
		PGClient: PGClient,
	}
}

func initUserService(PGClient *pgConfig.Client) *use_cases.UserService {
	log.Println("[routers:user] initializing UserService")

	repo := pgRepo.NewUserRepository(PGClient)

	service := use_cases.NewUserService(repo)

	log.Println("[routers:user] UserService initialized")

	return service
}

func (ur *UserRoute) RegisterRoutes() {

	log.Println("[routers:user] registering /user endpoints")

	api := ur.engine.Group("/user")

	service := initUserService(ur.PGClient)

	handler := user.NewUserHandler(service)

	{

		api.GET("/", handler.GetAllUsers)
		log.Println("[routers:user] registered GET /user/")

	}
}
