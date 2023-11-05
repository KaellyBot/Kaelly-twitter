package application

import (
	"github.com/kaellybot/kaelly-twitter/controllers"
	"github.com/kaellybot/kaelly-twitter/services"
	"github.com/rs/zerolog/log"
)

func New() *Impl {
	twitterService := services.New()
	twitterController := controllers.New(twitterService)

	return &Impl{
		service:    twitterService,
		controller: twitterController,
	}
}

func (app *Impl) Run() {
	app.controller.Run()
}

func (app *Impl) Shutdown() {
	log.Info().Msgf("Application is no longer running")
}
