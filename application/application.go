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
	broker := amqp.New(constants.RabbitMQClientID, viper.GetString(constants.RabbitMQAddress))
	db := databases.New()

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
		db:             db,
	}, nil
}

func (app *Impl) Run() error {
	if err := app.db.Run(); err != nil {
		return err
	}

	if err := app.broker.Run(); err != nil {
		return err
	}

	return app.twitterService.DispatchNewTweets()
}

func (app *Impl) Shutdown() {
	app.broker.Shutdown()
	app.db.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
