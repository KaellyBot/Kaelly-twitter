package twitteraccounts

import (
	"github.com/kaellybot/kaelly-twitter/models/entities"
	"github.com/kaellybot/kaelly-twitter/utils/databases"
)

type TwitterAccountRepository interface {
	GetTwitterAccounts() ([]entities.TwitterAccount, error)
	Save(twitterAccount entities.TwitterAccount) error
}

type TwitterAccountRepositoryImpl struct {
	db databases.MySQLConnection
}
