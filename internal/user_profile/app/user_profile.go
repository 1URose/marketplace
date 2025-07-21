package app

import (
	"github.com/1URose/marketplace/internal/common/app"
	"github.com/1URose/marketplace/internal/user_profile/transport/rest"

	"log"
)

func Run(deps *app.Deps) {
	log.Println("[user_profile] registering routers")

	rest.RegisterRoutes(deps)

	log.Println("[user_profile] routers registered successfully")
}
