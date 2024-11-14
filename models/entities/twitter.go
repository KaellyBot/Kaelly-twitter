package entities

import (
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
)

type TwitterAccount struct {
	ID         string `gorm:"primaryKey"`
	Name       string
	Game       amqp.Game
	Locale     amqp.Language
	LastUpdate time.Time `gorm:"not null; default:current_timestamp"`
}
