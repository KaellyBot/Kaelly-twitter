package twitter

import (
	"errors"
	"net/http"
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-twitter/repositories/twitteraccounts"
)

const (
	twitterURL            = "https://twitter.com"
	twitterAPIURL         = "https://twitter.com/i/api/graphql/BeHK76TOCY3P8nO-FWocjA/UserTweets"
	cookieGuestToken      = "gt"
	headerGuestToken      = "x-guest-token"
	variablesParameter    = "variables"
	featuresParameter     = "features"
	twitterEntryTypeTweet = "Tweet"
	routingkey            = "news.twitter"
)

var (
	errCookieNotFound = errors.New("cookie cannot be found")
)

type Service interface {
	DispatchNewTweets() error
}

type Impl struct {
	tweetCount          int
	token               string
	broker              amqp.MessageBroker
	client              http.Client
	twitterAccountsRepo twitteraccounts.Repository
}

type Tweet struct {
	ID        string
	URL       string
	CreatedAt time.Time
}
