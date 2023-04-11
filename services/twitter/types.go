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
)

var (
	errCookieNotFound = errors.New("Cookie cannot be found")
)

type TwitterService interface {
	CheckTweets() error
}

type TwitterServiceImpl struct {
	tweetCount          int
	token               string
	broker              amqp.MessageBrokerInterface
	client              http.Client
	twitterAccountsRepo twitteraccounts.TwitterAccountRepository
}

type Tweet struct {
	URL       string
	CreatedAt time.Time
}
