package twitter

import (
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/rs/zerolog/log"
)

func New(token string, timeout int, broker *amqp.MessageBroker) (*TwitterService, error) {
	return &TwitterService{
		token:   token,
		broker:  broker,
		timeout: time.Duration(timeout) * time.Second,
	}, nil
}

func (service *TwitterService) CheckTweets() {
	log.Info().Msgf("Twitter service is listening events...")

	// TODO
}
