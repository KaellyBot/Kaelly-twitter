package twitter

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-twitter/repositories/twitteraccounts"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

const (
	routingkey = "news.twitter"
)

type Service interface {
	DispatchNewTweets() error
}

type Impl struct {
	tweetCount          int
	username            string
	password            string
	broker              amqp.MessageBroker
	scraper             *twitterscraper.Scraper
	twitterAccountsRepo twitteraccounts.Repository
}
