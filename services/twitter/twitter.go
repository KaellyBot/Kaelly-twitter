package twitter

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"sync"
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-twitter/models/constants"
	"github.com/kaellybot/kaelly-twitter/models/entities"
	"github.com/kaellybot/kaelly-twitter/repositories/twitteraccounts"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func New(twitterAccountsRepo twitteraccounts.Repository, broker amqp.MessageBroker) (*Impl, error) {
	return &Impl{
		tweetCount:          viper.GetInt(constants.TwitterTweetCount),
		token:               viper.GetString(constants.TwitterBearerToken),
		twitterAccountsRepo: twitterAccountsRepo,
		broker:              broker,
		client: http.Client{
			Timeout: time.Duration(viper.GetInt(constants.TwitterTimeout)) * time.Second,
		},
	}, nil
}

func (service *Impl) DispatchNewTweets() error {
	log.Info().Msgf("Retrieving tweets from Twitter...")

	twitterAccounts, err := service.twitterAccountsRepo.GetTwitterAccounts()
	if err != nil {
		return err
	}

	guestToken, err := service.getGuestToken()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, account := range twitterAccounts {
		wg.Add(1)
		go func(twitterAccount entities.TwitterAccount) {
			defer wg.Done()
			service.checkTwitterAccount(twitterAccount, guestToken)
		}(account)
	}

	wg.Wait()
	return nil
}

func (service *Impl) checkTwitterAccount(account entities.TwitterAccount, guestToken string) {
	log.Info().
		Str(constants.LogLanguage, account.Locale.String()).
		Str(constants.LogTwitterID, account.ID).
		Msgf("Reading tweets...")

	tweets, err := service.getUserTweets(service.token, guestToken, account.ID)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogTwitterID, account.ID).
			Msgf("Cannot retrieve tweet from the twitter account, ignored")
		return
	}

	publishedTweets := 0
	lastUpdate := account.LastUpdate
	for _, tweet := range tweets {
		if tweet.CreatedAt.UTC().After(lastUpdate.UTC()) {
			errPublish := service.publishTweet(tweet, account.Locale)
			if errPublish != nil {
				log.Error().Err(err).
					Str(constants.LogCorrelationID, tweet.ID).
					Str(constants.LogTwitterID, account.ID).
					Str(constants.LogTweetID, tweet.ID).
					Msgf("Impossible to publish tweet, breaking loop")
				break
			}

			account.LastUpdate = tweet.CreatedAt.UTC()
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

func (service *Impl) publishTweet(tweet Tweet, lg amqp.Language) error {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_NEWS_TWITTER,
		Language: lg,
		NewsTwitterMessage: &amqp.NewsTwitterMessage{
			Url:  tweet.URL,
			Date: timestamppb.New(tweet.CreatedAt),
		},
	}

	return service.broker.Publish(&message, amqp.ExchangeNews, routingkey, tweet.ID)
}

func (service *Impl) getGuestToken() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), service.client.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, twitterURL, nil)
	if err != nil {
		return "", err
	}

	res, err := service.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusOK {
		for _, cookie := range res.Cookies() {
			if cookie.Name == cookieGuestToken {
				return cookie.Value, nil
			}
		}
	} else {
		return "", fmt.Errorf("cannot consume twitter API, guest_token could not be retrieved: %d", res.StatusCode)
	}

	return "", errCookieNotFound
}

func (service *Impl) getUserTweets(bearerToken, guestToken, userID string) ([]Tweet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), service.client.Timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, twitterAPIURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	req.Header.Set(headerGuestToken, guestToken)

	q := req.URL.Query()
	q.Add(variablesParameter, getVariables(userID, service.tweetCount))
	q.Add(featuresParameter, getFeatures())
	req.URL.RawQuery = q.Encode()

	resp, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return castResponse(userID, resp)
	}

	return nil, fmt.Errorf("cannot consume twitter API (userID=%s): %d", userID, resp.StatusCode)
}

func castResponse(userID string, entity *http.Response) ([]Tweet, error) {
	body, err := io.ReadAll(entity.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %w", err)
	}

	rootNode := gjson.ParseBytes(body)
	dataNode := rootNode.Get("data.user.result.timeline.timeline.instructions.#(type==\"TimelineAddEntries\")" +
		".entries.#.content.itemContent.tweet_results.result")

	var tweets []Tweet
	for _, tweetData := range dataNode.Array() {
		if !isEntryOriginalTweet(tweetData) {
			continue
		}

		restID := tweetData.Get("rest_id").String()
		createdAtStr := tweetData.Get("legacy.created_at").String()
		createdAt, errParsing := time.Parse(time.RubyDate, createdAtStr)
		if errParsing != nil {
			return nil, fmt.Errorf("cannot parse tweet created_at: %w", err)
		}

		tweets = append(tweets, Tweet{
			ID:        restID,
			URL:       fmt.Sprintf("%s/%s/status/%s", twitterURL, userID, restID),
			CreatedAt: createdAt,
		})
	}

	sort.SliceStable(tweets, func(i, j int) bool {
		return tweets[i].CreatedAt.Before(tweets[j].CreatedAt)
	})

	return tweets, nil
}

func getVariables(userID string, tweetCount int) string {
	return fmt.Sprintf(`{"userID":"%s",
                "count":%d,
                "includePromotedContent":false,
                "withVoice":false,
                "withDownvotePerspective":false,
                "withReactionsMetadata":false,
                "withReactionsPerspective":false}`,
		userID, tweetCount)
}

func getFeatures() string {
	return "{\"blue_business_profile_image_shape_enabled\":false" +
		",\"responsive_web_graphql_exclude_directive_enabled\":false" +
		",\"verified_phone_label_enabled\":false" +
		",\"responsive_web_graphql_timeline_navigation_enabled\":false" +
		",\"responsive_web_graphql_skip_user_profile_image_extensions_enabled\":false" +
		",\"tweetypie_unmention_optimization_enabled\":false" +
		",\"vibe_api_enabled\":false" +
		",\"responsive_web_edit_tweet_api_enabled\":false" +
		",\"graphql_is_translatable_rweb_tweet_is_translatable_enabled\":false" +
		",\"view_counts_everywhere_api_enabled\":false" +
		",\"longform_notetweets_consumption_enabled\":false" +
		",\"tweet_awards_web_tipping_enabled\":false" +
		",\"freedom_of_speech_not_reach_fetch_enabled\":false" +
		",\"standardized_nudges_misinfo\":false" +
		",\"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled\":false" +
		",\"interactive_text_enabled\":false" +
		",\"responsive_web_text_conversations_enabled\":false" +
		",\"longform_notetweets_richtext_consumption_enabled\":false" +
		",\"responsive_web_enhance_cards_enabled\":false}"
}

func isEntryOriginalTweet(result gjson.Result) bool {
	return result.Get("__typename").String() == twitterEntryTypeTweet &&
		!result.Get("legacy.retweeted_status_result").Exists()
}
