package services

import (
	twitterscraper "github.com/n0madic/twitter-scraper"
)

type Service interface {
	RetrieveTweets() (map[string][]*twitterscraper.Tweet, error)
}

type Impl struct {
	tweetCount int
	username   string
	password   string
	scraper    *twitterscraper.Scraper
}
