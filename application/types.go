package application

import (
	"errors"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-twitter/services/twitter"
)

var (
	ErrCannotInstanciateApp = errors.New("Cannot instanciate application")
)

type ApplicationInterface interface {
	Run() error
	Shutdown()
}

type Application struct {
	twitter twitter.TwitterServiceInterface
	broker  amqp.MessageBrokerInterface
}
