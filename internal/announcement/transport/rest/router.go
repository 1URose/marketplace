package rest

import (
	routers "github.com/1URose/marketplace/internal/announcement/transport/rest/routes"
	"github.com/1URose/marketplace/internal/common/app"
	"log"
)

func RegisterRoutes(deps *app.Deps) {
	log.Println("[rest:announcement] registering announcement routers")

	adRoute := routers.NewAdRoute(deps)

	adRoute.RegisterRoutes()

	log.Println("[rest:announcement] announcement routers registered successfully")
}
