package twitter

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-twitter/models/constants"
	"github.com/kaellybot/kaelly-twitter/repositories/twitteraccounts"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
)

func New(twitterAccountsRepo twitteraccounts.TwitterAccountRepository,
	broker *amqp.MessageBroker) (*TwitterServiceImpl, error) {

	return &TwitterServiceImpl{
		tweetCount:          viper.GetInt(constants.TwitterTweetCount),
		token:               viper.GetString(constants.TwitterBearerToken),
		twitterAccountsRepo: twitterAccountsRepo,
		broker:              broker,
		client: http.Client{
			Timeout: time.Duration(viper.GetInt(constants.TwitterTimeout)) * time.Second,
		},
	}, nil
}

func (service *TwitterServiceImpl) CheckTweets() error {
	log.Info().Msgf("Retrieving tweets from Twitter...")

	twitterAccounts, err := service.twitterAccountsRepo.GetTwitterAccounts()
	if err != nil {
		return err
	}

	responses := make(map[amqp.Language][]Tweet, 0)
	guestToken, err := service.getGuestToken()
	if err != nil {
		return err
	}

	for _, account := range twitterAccounts {
		tweets, err := service.getUserTweets(service.token, guestToken, account.Id)
		if err != nil {
			return err
		}
		responses[account.Language] = tweets
	}

	// TODO
	log.Info().Msgf("Tweets: %v", responses)

	return nil
}

func (service *TwitterServiceImpl) getGuestToken() (string, error) {
	req, err := http.NewRequest("HEAD", twitterURL, nil)
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
		return "", fmt.Errorf("Cannot consume twitter API, guest_token could not be retrieved: %d", res.StatusCode)
	}

	return "", errCookieNotFound
}

func (service *TwitterServiceImpl) getUserTweets(bearerToken, guestToken, userId string) ([]Tweet, error) {
	req, err := http.NewRequest("GET", twitterAPIURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))
	req.Header.Set(headerGuestToken, guestToken)

	q := req.URL.Query()
	q.Add(variablesParameter, getVariables(userId, service.tweetCount))
	q.Add(featuresParameter, getFeatures())
	req.URL.RawQuery = q.Encode()

	resp, err := service.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return castResponse(userId, resp)
	} else {
		return nil, fmt.Errorf("Cannot consume twitter API (userId=%s): %d", userId, resp.StatusCode)
	}
}

func castResponse(userId string, entity *http.Response) ([]Tweet, error) {
	body, err := ioutil.ReadAll(entity.Body)
	if err != nil {
		return nil, fmt.Errorf("cannot read response body: %w", err)
	}

	rootNode := gjson.ParseBytes(body)
	dataNode := rootNode.Get("data.user.result.timeline.timeline.instructions.#(type==\"TimelineAddEntries\").entries.#.content.itemContent.tweet_results.result")

	var tweets []Tweet
	for _, tweetData := range dataNode.Array() {
		if !isEntryOriginalTweet(tweetData) {
			continue
		}

		restId := tweetData.Get("rest_id").String()
		createdAtStr := tweetData.Get("legacy.created_at").String()
		createdAt, err := time.Parse(time.RubyDate, createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("cannot parse tweet created_at: %w", err)
		}

		tweets = append(tweets, Tweet{
			URL:       fmt.Sprintf("%s/%s/status/%s", twitterURL, userId, restId),
			CreatedAt: createdAt,
		})
	}

	return tweets, nil
}

func getVariables(userId string, tweetCount int) string {
	return fmt.Sprintf(`{"userId":"%s",
                "count":%d,
                "includePromotedContent":false,
                "withVoice":false,
                "withDownvotePerspective":false,
                "withReactionsMetadata":false,
                "withReactionsPerspective":false}`,
		userId, tweetCount)
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
