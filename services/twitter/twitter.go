package twitter

import (
	"net/http"
	"sort"
	"sync"
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-twitter/models/constants"
	"github.com/kaellybot/kaelly-twitter/models/entities"
	"github.com/kaellybot/kaelly-twitter/repositories/twitteraccounts"
	twitterscraper "github.com/n0madic/twitter-scraper"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func New(twitterAccountsRepo twitteraccounts.Repository, broker amqp.MessageBroker) (*Impl, error) {
	return &Impl{
		authToken:           viper.GetString(constants.TwitterAuthToken),
		csrfToken:           viper.GetString(constants.TwitterCSRFToken),
		tweetCount:          viper.GetInt(constants.TwitterTweetCount),
		twitterAccountsRepo: twitterAccountsRepo,
		broker:              broker,
		scraper:             twitterscraper.New(),
	}, nil
}

func (service *Impl) DispatchNewTweets() error {
	log.Info().Msgf("Retrieving tweets from Twitter...")

	twitterAccounts, err := service.twitterAccountsRepo.GetTwitterAccounts()
	if err != nil {
		return err
	}

	service.scraper.SetCookies([]*http.Cookie{
		{
			Name:  "auth_token",
			Value: service.authToken,
		},
		{
			Name:  "ct0",
			Value: service.csrfToken,
		},
	})

	var wg sync.WaitGroup
	for _, account := range twitterAccounts {
		wg.Add(1)
		go func(twitterAccount entities.TwitterAccount) {
			defer wg.Done()
			service.checkTwitterAccount(twitterAccount)
		}(account)
	}

	wg.Wait()
	return nil
}

func (service *Impl) checkTwitterAccount(account entities.TwitterAccount) {
	log.Info().
		Str(constants.LogLanguage, account.Locale.String()).
		Str(constants.LogTwitterID, account.ID).
		Msgf("Reading tweets...")

	tweets, _, err := service.scraper.FetchTweets(account.Name, service.tweetCount, "")
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogTwitterID, account.ID).
			Msgf("Cannot retrieve tweets from account, ignored")
		return
	}

	publishedTweets := 0
	lastUpdate := account.LastUpdate
	tweets = service.keepInterestingTweets(tweets)

	for _, tweet := range tweets {
		utcDate := time.Unix(tweet.Timestamp, 0).UTC()
		if utcDate.After(lastUpdate.UTC()) {
			errPublish := service.publishTweet(tweet, account.Game, account.Locale)
			if errPublish != nil {
				log.Error().Err(err).
					Str(constants.LogCorrelationID, tweet.ID).
					Str(constants.LogTwitterID, account.ID).
					Str(constants.LogTweetID, tweet.ID).
					Msgf("Impossible to publish tweet, breaking loop")
				break
			}

			account.LastUpdate = utcDate
			err = service.twitterAccountsRepo.Save(account)
			if err != nil {
				log.Error().Err(err).
					Str(constants.LogCorrelationID, tweet.ID).
					Str(constants.LogTwitterID, account.ID).
					Str(constants.LogTweetID, tweet.ID).
					Msgf("Impossible to update account, breaking loop; this tweet might be published again next time")
				break
			}

			publishedTweets++
		}
	}

	log.Info().
		Str(constants.LogLanguage, account.Locale.String()).
		Str(constants.LogTwitterID, account.ID).
		Int(constants.LogTweetNumber, publishedTweets).
		Msgf("Tweet(s) read and published")
}

func (service *Impl) publishTweet(tweet *twitterscraper.Tweet, game amqp.Game, lg amqp.Language) error {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_NEWS_TWITTER,
		Language: lg,
		Game:     game,
		NewsTwitterMessage: &amqp.NewsTwitterMessage{
			Url:  tweet.PermanentURL,
			Date: timestamppb.New(time.Unix(tweet.Timestamp, 0).UTC()),
		},
	}

	return service.broker.Emit(&message, amqp.ExchangeNews, routingkey, tweet.ID)
}

func (service *Impl) keepInterestingTweets(tweets []*twitterscraper.Tweet) []*twitterscraper.Tweet {
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

	return result
}
