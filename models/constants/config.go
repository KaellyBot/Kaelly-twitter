package constants

import "github.com/rs/zerolog"

const (
	ConfigFileName = ".env"

	// MySQL URL with the following format: HOST:PORT
	MySqlUrl = "MYSQL_URL"

	// MySQL user
	MySqlUser = "MYSQL_USER"

	// MySQL password
	MySqlPassword = "MYSQL_PASSWORD"

	// MySQL database name
	MySqlDatabase = "MYSQL_DATABASE"

	// RabbitMQ address
	RabbitMqAddress = "RABBITMQ_ADDRESS"

	// Bearer Token used to consume the Twitter GraphQL API
	TwitterBearerToken = "TWITTER_BEARER_TOKEN"

	// Number of tweets retrieved per call
	TwitterTweetCount = "TWEET_COUNT"

	// Timeout to retrieve tweets in seconds
	TwitterTimeout = "HTTP_TIMEOUT"

	// Zerolog values from [trace, debug, info, warn, error, fatal, panic]
	LogLevel = "LOG_LEVEL"

	// Boolean; used to register commands at development guild level or globally.
	Production = "PRODUCTION"
)

var (
	DefaultConfigValues = map[string]interface{}{
		MySqlUrl:           "localhost:3306",
		MySqlUser:          "",
		MySqlPassword:      "",
		MySqlDatabase:      "kaellybot",
		RabbitMqAddress:    "amqp://localhost:5672",
		TwitterBearerToken: "",
		TwitterTweetCount:  20,
		TwitterTimeout:     60,
		LogLevel:           zerolog.InfoLevel.String(),
		Production:         false,
	}
)
