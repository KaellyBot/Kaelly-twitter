package models

const (
	Name             = "KaellyBot"
	TwitterUserAgent = Name
	// Don't stress out, this is not a leak but a default bearer token used by Twitter ;)
	TwitterBearerToken = "AAAAAAAAAAAAAAAAAAAAANRILgAAAAAAnNwIzUejRCOuH5E6I8xnZz4puTs%3D1Zv7ttfk8LF81IUq16cHjhLTvJu4FA33AGWWjCpTnA"
	TwitterTweetCount  = 20
	TwitterTimeout     = 60
	RabbitMqAddress    = "amqp://localhost:5672"
	RabbitMqClientId   = "Kaelly-Twitter"
)
