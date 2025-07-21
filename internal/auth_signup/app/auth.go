package app

import (
	"github.com/1URose/marketplace/internal/auth_signup/transport/rest"
	"github.com/1URose/marketplace/internal/common/app"
	"log"
)

func Run(deps *app.Deps) {
	log.Println("[auth_signup] registering auth routers")

	rest.RegisterRoutes(deps)

	log.Println("[auth_signup] auth routers registered successfully")
}
