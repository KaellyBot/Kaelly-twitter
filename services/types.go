package services

import (
	twitterscraper "github.com/imperatrona/twitter-scraper"
)

type Service interface {
	RetrieveTweets() (map[string][]*twitterscraper.Tweet, error)
}

type Impl struct {
	tweetCount   int
	loginErrored bool
	username     string
	password     string
	scraper      *twitterscraper.Scraper
}
