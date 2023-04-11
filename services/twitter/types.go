package twitter

import (
	"errors"
	"net/http"

	amqp "github.com/kaellybot/kaelly-amqp"
)

const (
	twitterURL            = "https://twitter.com"
	twitterAPIURL         = "https://twitter.com/i/api/graphql/BeHK76TOCY3P8nO-FWocjA/UserTweets"
	cookieGuestToken      = "gt"
	headerGuestToken      = "x-guest-token"
	variablesParameter    = "variables"
	featuresParameter     = "features"
	twitterEntryTypeTweet = "Tweet"
)

var (
	errCookieNotFound = errors.New("Cookie cannot be found")
)

type TwitterServiceInterface interface {
	CheckTweets() error
}

type TwitterService struct {
	tweetCount int
	token      string
	broker     amqp.MessageBrokerInterface
	client     http.Client
}
