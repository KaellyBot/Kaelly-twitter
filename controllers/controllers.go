package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kaellybot/kaelly-twitter/models/dtos"
	"github.com/kaellybot/kaelly-twitter/models/mappers"
	"github.com/kaellybot/kaelly-twitter/services"
	"github.com/rs/zerolog/log"
)

func New(service services.Service) *Impl {
	controller := &Impl{
		r:       gin.Default(),
		service: service,
	}

	controller.r.GET("/tweets", func(c *gin.Context) {
		log.Info().Msgf("Responding to /tweets call...")
		c.JSON(http.StatusOK, controller.retrieveTweets())
	})

	return controller
}

func (controller *Impl) Run() {
	go func() {
		controller.r.Run()
	}()
}

func (controller *Impl) retrieveTweets() map[string][]dtos.TweetDto {
	tweets, err := controller.service.RetrieveTweets()
	if err != nil {
		log.Error().Err(err).
			Msgf("Cannot retrieve tweets, returning empty collection...")
		return nil
	}
	return mappers.MapTweets(tweets)
}
