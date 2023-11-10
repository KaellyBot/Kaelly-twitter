package services

import (
	"sort"

	"github.com/kaellybot/kaelly-twitter/models/constants"
	"github.com/kaellybot/kaelly-twitter/webhooks"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() *Impl {
	return &Impl{
		tweetCount:   viper.GetInt(constants.TwitterTweetCount),
		username:     viper.GetString(constants.TwitterUsername),
		password:     viper.GetString(constants.TwitterPassword),
		scraper:      twitterscraper.New(),
		loginErrored: false,
	}
}

func (service *Impl) RetrieveTweets() (map[string][]*twitterscraper.Tweet, error) {
	log.Info().Msgf("Retrieving tweets from Twitter...")

	err := service.scraper.Login(service.username, service.password)
	if err != nil {
		if !service.loginErrored {
			service.loginErrored = true
			webhooks.SendWebhookMessage(":warning: Twitter login failed, manual action required!")
		}

		return nil, err
	}
	service.loginErrored = false
	defer service.scraper.Logout()

	result := make(map[string][]*twitterscraper.Tweet)
	for _, account := range constants.GetTwitterAccounts() {
		tweets, err := service.checkTwitterAccount(account)
		if err != nil {
			log.Error().Err(err).
				Str(constants.LogTwitterID, account.Username).
				Msgf("Cannot retrieve tweets from account, continuing...")
			continue
		}
		result[account.Locale] = tweets
	}

	return result, nil
}

func (service *Impl) checkTwitterAccount(account constants.TwitterAccount) ([]*twitterscraper.Tweet, error) {
	log.Info().
		Str(constants.LogLanguage, account.Locale).
		Str(constants.LogTwitterID, account.Username).
		Msgf("Reading tweets...")

	tweets, _, err := service.scraper.FetchTweets(account.Username, service.tweetCount, "")
	if err != nil {
		return nil, err
	}

	result := make([]*twitterscraper.Tweet, 0)
	for _, tweet := range tweets {
		// Exclude RTs
		if tweet.RetweetedStatus != nil {
			continue
		}

		result = append(result, tweet)
	}

	sort.SliceStable(result, func(i, j int) bool {
		return result[i].Timestamp < result[j].Timestamp
	})

	return result, nil
}
