package application

import (
	"github.com/kaellybot/kaelly-twitter/controllers"
	"github.com/kaellybot/kaelly-twitter/services"
)

type Application interface {
	Run()
	Shutdown()
}

type Impl struct {
	controller controllers.Controller
	service    services.Service
}
