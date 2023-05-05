package constants

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

	// Default values.
	defaultMySQLURLValue           = "localhost:3306"
	defaultMySQLUserValue          = ""
	defaultMySQLPasswordValue      = ""
	defaultMySQLDatabaseValue      = "kaellybot"
	defaultRabbitMQAddressValue    = "amqp://localhost:5672"
	defaultTwitterBearerTokenValue = ""
	defaultTwitterTweetCountValue  = 20
	defaultTwitterTimeoutValue     = 60
	defaultMetricPortValue         = 2112
	defaultLogLevelValue           = "info"
	defaultProductionValue         = false
)

func GetDefaultConfigValues() map[string]any {
	return map[string]any{
		MySQLURL:           defaultMySQLURLValue,
		MySQLUser:          defaultMySQLUserValue,
		MySQLPassword:      defaultMySQLPasswordValue,
		MySQLDatabase:      defaultMySQLDatabaseValue,
		RabbitMQAddress:    defaultRabbitMQAddressValue,
		TwitterBearerToken: defaultTwitterBearerTokenValue,
		TwitterTweetCount:  defaultTwitterTweetCountValue,
		TwitterTimeout:     defaultTwitterTimeoutValue,
		MetricPort:         defaultMetricPortValue,
		LogLevel:           defaultLogLevelValue,
		Production:         defaultProductionValue,
	}
}
