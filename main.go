package main

import (
	"flag"
	"fmt"

	"github.com/kaellybot/kaelly-twitter/application"
	"github.com/kaellybot/kaelly-twitter/models"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	twitterTimeout  int
	rabbitMqAddress string
	token           string
)

func init() {
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		file = short
		return fmt.Sprintf("%s:%d", file, line)
	}
	log.Logger = log.With().Caller().Logger()

	flag.IntVar(&twitterTimeout, "rssTimeout", models.TwitterTimeout, "Timeout to retrieve tweets in seconds")
	flag.StringVar(&rabbitMqAddress, "rabbitMqAddress", models.RabbitMqAddress, "RabbitMQ address")
	flag.StringVar(&token, "token", "", "Twitter Bot Token")
	flag.Parse()
}

func main() {
	app, err := application.New(token, models.RabbitMqClientId, rabbitMqAddress, twitterTimeout)
	if err != nil {
		log.Fatal().Err(err).Msgf("Shutting down after failing to instanciate application")
	}

	app.Run()

	log.Info().Msgf("Gracefully shutting down %s...", models.Name)
	app.Shutdown()
}
