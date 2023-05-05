package constants

import "github.com/rs/zerolog"

const (
	ConfigFileName = ".env"

	// MySQL URL with the following format: HOST:PORT.
	MySQLURL = "MYSQL_URL"

	// MySQL user.
	MySQLUser = "MYSQL_USER"

	// MySQL password.
	MySQLPassword = "MYSQL_PASSWORD"

	// MySQL database name.
	MySQLDatabase = "MYSQL_DATABASE"

	// RabbitMQ address.
	RabbitMQAddress = "RABBITMQ_ADDRESS"

	// Bearer Token used to consume the Twitter GraphQL API.
	TwitterBearerToken = "TWITTER_BEARER_TOKEN" // #nosec G101

	// Number of tweets retrieved per call.
	TwitterTweetCount = "TWEET_COUNT"

	// Timeout to retrieve tweets in seconds.
	TwitterTimeout = "HTTP_TIMEOUT"

	// Metric port.
	MetricPort = "METRIC_PORT"

	// Zerolog values from [trace, debug, info, warn, error, fatal, panic].
	LogLevel = "LOG_LEVEL"

	// Boolean; used to register commands at development guild level or globally.
	Production = "PRODUCTION"

	defaultMySQLURL           = "localhost:3306"
	defaultMySQLUser          = ""
	defaultMySQLPassword      = ""
	defaultMySQLDatabase      = "kaellybot"
	defaultRabbitMQAddress    = "amqp://localhost:5672"
	defaultTwitterBearerToken = ""
	defaultTwitterTweetCount  = 20
	defaultTwitterTimeout     = 60
	defaultMetricPort         = 2112
	defaultLogLevel           = zerolog.InfoLevel
	defaultProduction         = false
)

func GetDefaultConfigValues() map[string]any {
	return map[string]any{
		MySQLURL:           defaultMySQLURL,
		MySQLUser:          defaultMySQLUser,
		MySQLPassword:      defaultMySQLPassword,
		MySQLDatabase:      defaultMySQLDatabase,
		RabbitMQAddress:    defaultRabbitMQAddress,
		TwitterBearerToken: defaultTwitterBearerToken,
		TwitterTweetCount:  defaultTwitterTweetCount,
		TwitterTimeout:     defaultTwitterTimeout,
		MetricPort:         defaultMetricPort,
		LogLevel:           defaultLogLevel.String(),
		Production:         defaultProduction,
	}
}
