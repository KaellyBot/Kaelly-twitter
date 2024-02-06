package mappers

import (
	"strings"
	"time"

	"github.com/kaellybot/kaelly-twitter/models/constants"
	"github.com/kaellybot/kaelly-twitter/models/dtos"
	twitterscraper "github.com/n0madic/twitter-scraper"
)

func MapTweets(tweets []*twitterscraper.Tweet) []dtos.TweetDto {
	result := make([]dtos.TweetDto, 0)

	for _, tweet := range tweets {
		result = append(result, dtos.TweetDto{
			Url:       mapUrl(tweet.PermanentURL),
			CreatedAt: time.Unix(tweet.Timestamp, 0).UTC(),
		})
	}

	return result
}

func mapUrl(url string) string {
	return strings.Replace(url, constants.UrlClassic, constants.UrlPreview, 1)
}
