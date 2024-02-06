package services

import (
	twitterscraper "github.com/n0madic/twitter-scraper"
)

type Service interface {
	RetrieveTweets() ([]*twitterscraper.Tweet, error)
}

type Impl struct {
	tweetCount   int
	loginErrored bool
	username     string
	password     string
	scraper      *twitterscraper.Scraper
}
