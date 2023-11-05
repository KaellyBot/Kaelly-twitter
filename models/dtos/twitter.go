package dtos

import "time"

type TweetDto struct {
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
}
