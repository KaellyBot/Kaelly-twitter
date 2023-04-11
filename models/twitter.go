package models

import (
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
)

var (
	TwitterIDs = map[amqp.Language]string{
		amqp.Language_FR: "72272795",
	}
)

type Tweet struct {
	URL       string
	CreatedAt time.Time
}
