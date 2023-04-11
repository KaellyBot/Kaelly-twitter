module github.com/kaellybot/kaelly-twitter

go 1.18

// replace github.com/kaellybot/kaelly-amqp => /home/kaysoro/git/kaelly-amqp

require (
	github.com/kaellybot/kaelly-amqp v0.0.1-beta5
	github.com/rs/zerolog v1.26.1
	github.com/tidwall/gjson v1.14.4
)

require (
	github.com/golang/protobuf v1.5.0 // indirect
	github.com/rabbitmq/amqp091-go v1.3.4 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
