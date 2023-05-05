package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-twitter/services/twitter"
)

type Interface interface {
	Run() error
	Shutdown()
}

type Application struct {
	twitterService twitter.Service
	broker         amqp.MessageBrokerInterface
}
