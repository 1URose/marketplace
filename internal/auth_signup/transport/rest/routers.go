package rest

import (
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest/routers"
	"github.com/1URose/marketplace/internal/common/app"
	"log"
)

func RegisterRoutes(deps *app.Deps) {
	log.Println("[rest:auth_signup] registering auth routers")

	authRouter := routers.NewAuthRouter(deps)
	authRouter.RegisterRoutes()

	log.Println("[rest:auth_signup] auth routers registered successfully")
}
