package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-twitter/services/twitter"
)

type Application interface {
	Run() error
	Shutdown()
}

type Impl struct {
	twitterService twitter.Service
	broker         amqp.MessageBroker
}
