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

func New() (*Application, error) {
	// misc
	db, err := databases.New()
	if err != nil {
		return nil, err
	}

	broker, err := amqp.New(constants.RabbitMQClientId, viper.GetString(constants.RabbitMqAddress), nil)
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

	return &Application{
		twitterService: twitterService,
		broker:         broker,
	}, nil
}

func (app *Application) Run() error {
	return app.twitterService.DispatchNewTweets()
}

func (app *Application) Shutdown() {
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
