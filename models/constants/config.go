package constants

import "github.com/rs/zerolog"

const (
	ConfigFileName = ".env"

	// Username used to login to Twitter.
	TwitterUsername = "TWITTER_USERNAME"

	// Password used to login to Twitter.
	TwitterPassword = "TWITTER_PASSWORD"

	// Number of tweets retrieved per call.
	TwitterTweetCount = "TWEET_COUNT"

	// Metric port.
	MetricPort = "METRIC_PORT"

	// Zerolog values from [trace, debug, info, warn, error, fatal, panic].
	LogLevel = "LOG_LEVEL"

	// Boolean; used to register commands at development guild level or globally.
	Production = "PRODUCTION"

	defaultTwitterUsername   = ""
	defaultTwitterPassword   = ""
	defaultTwitterTweetCount = 20
	defaultMetricPort        = 2112
	defaultLogLevel          = zerolog.InfoLevel
	defaultProduction        = false
)

func GetDefaultConfigValues() map[string]any {
	return map[string]any{
		TwitterUsername:   defaultTwitterUsername,
		TwitterPassword:   defaultTwitterPassword,
		TwitterTweetCount: defaultTwitterTweetCount,
		MetricPort:        defaultMetricPort,
		LogLevel:          defaultLogLevel.String(),
		Production:        defaultProduction,
	}
}
