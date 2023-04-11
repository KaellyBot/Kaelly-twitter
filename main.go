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
	twitterTweetCount int
	twitterTimeout    int
	rabbitMqAddress   string
	twitterToken      string
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

	flag.IntVar(&twitterTweetCount, "tweetCount", models.TwitterTimeout, "Tweet count retrieved for each account")
	flag.IntVar(&twitterTimeout, "twitterTimeout", models.TwitterTimeout, "Timeout to retrieve tweets in seconds")
	flag.StringVar(&rabbitMqAddress, "rabbitMqAddress", models.RabbitMqAddress, "RabbitMQ address")
	flag.StringVar(&twitterToken, "twitterToken", models.TwitterBearerToken, "Twitter Bearer Token")
	flag.Parse()
}

func main() {
	app, err := application.New(twitterToken, models.RabbitMqClientId, rabbitMqAddress,
		twitterTweetCount, twitterTimeout)
	if err != nil {
		log.Fatal().Err(err).Msgf("Shutting down after failing to instanciate application")
	}

	app.Run()

	log.Info().Msgf("Gracefully shutting down %s...", models.Name)
	app.Shutdown()
}
