package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-twitter/models/constants"
	"github.com/kaellybot/kaelly-twitter/repositories/twitteraccounts"
	"github.com/kaellybot/kaelly-twitter/services/twitter"
	"github.com/kaellybot/kaelly-twitter/utils/databases"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() (*Impl, error) {
	// misc
	db, err := databases.New()
	if err != nil {
		return nil, err
	}

	broker, err := amqp.New(constants.RabbitMQClientID, viper.GetString(constants.RabbitMQAddress), nil)
	if err != nil {
		return nil, err
	}

	// repositories
	twitterRepo := twitteraccounts.New(db)

	// services
	twitterService, err := twitter.New(twitterRepo, broker)
	if err != nil {
		return nil, err
	}

	return &Impl{
		twitterService: twitterService,
		broker:         broker,
	}, nil
}

func (app *Impl) Run() error {
	return app.twitterService.DispatchNewTweets()
}

func (app *Impl) Shutdown() {
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
