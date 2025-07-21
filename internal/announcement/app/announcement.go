package app

import (
	"github.com/1URose/marketplace/internal/announcement/transport/rest"
	"github.com/1URose/marketplace/internal/common/app"
	"log"
)

func Run(deps *app.Deps) {
	log.Println("[announcement] registering routers")

	rest.RegisterRoutes(deps)
	log.Println("[announcement] routers registered successfully")
}
