package rest

import (
	"github.com/1URose/marketplace/internal/common/app"
	"github.com/1URose/marketplace/internal/user_profile/transport/rest/routes"
	"log"
)

func RegisterRoutes(deps *app.Deps) {
	log.Println("[rest:user_profile] registering user_profile routers")

	userRoute := routes.NewUserRoute(deps)

	userRoute.RegisterRoutes()

	log.Println("[rest:user_profile] user_profile routers registered")
}
