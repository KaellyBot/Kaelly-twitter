package twitteraccounts

import (
	"github.com/kaellybot/kaelly-twitter/models/entities"
	"github.com/kaellybot/kaelly-twitter/utils/databases"
)

type Repository interface {
	GetTwitterAccounts() ([]entities.TwitterAccount, error)
	Save(twitterAccount entities.TwitterAccount) error
}

type Impl struct {
	db databases.MySQLConnection
}
