package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-twitter/services/twitter"
)

type ApplicationInterface interface {
	Run() error
	Shutdown()
}

type Application struct {
	twitterService twitter.TwitterService
	broker         amqp.MessageBrokerInterface
}
