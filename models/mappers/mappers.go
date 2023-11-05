package mappers

import (
	"time"

	"github.com/kaellybot/kaelly-twitter/models/dtos"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

func MapTweets(localizedTweets map[string][]*twitterscraper.Tweet) map[string][]dtos.TweetDto {
	result := make(map[string][]dtos.TweetDto)

	for locale, tweets := range localizedTweets {
		result[locale] = mapTweets(tweets)
	}

	return result
}

func mapTweets(tweets []*twitterscraper.Tweet) []dtos.TweetDto {
	result := make([]dtos.TweetDto, 0)

	for _, tweet := range tweets {
		result = append(result, dtos.TweetDto{
			Url:       tweet.PermanentURL,
			CreatedAt: time.Unix(tweet.Timestamp, 0).UTC(),
		})
	}

	return result
}
