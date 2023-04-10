package twitter

import (
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
)

type TwitterServiceInterface interface {
	CheckTweets()
}

type TwitterService struct {
	token   string
	broker  amqp.MessageBrokerInterface
	timeout time.Duration
}
