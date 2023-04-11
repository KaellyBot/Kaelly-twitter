package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-twitter/services/twitter"
	"github.com/rs/zerolog/log"
)

func New(twitterToken, rabbitMqClientId, rabbitMqAddress string, tweetCount, twitterTimeout int) (*Application, error) {
	broker, err := amqp.New(rabbitMqClientId, rabbitMqAddress, []amqp.Binding{})
	if err != nil {
		log.Error().Err(err).Msgf("Failed to instantiate broker")
		return nil, ErrCannotInstanciateApp
	}

	twitter, err := twitter.New(twitterToken, tweetCount, twitterTimeout, broker)
	if err != nil {
		log.Error().Err(err).Msgf("Twitter service instantiation failed")
		return nil, err
	}

	return &Application{
		twitter: twitter,
		broker:  broker,
	}, nil
}

func (app *Application) Run() {
	err := app.twitter.CheckTweets()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to check tweets")
	}
}

func (app *Application) Shutdown() {
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
