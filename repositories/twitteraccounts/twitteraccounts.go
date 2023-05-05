package twitteraccounts

import (
	"github.com/kaellybot/kaelly-twitter/models/entities"
	"github.com/kaellybot/kaelly-twitter/utils/databases"
)

func New(db databases.MySQLConnection) *Impl {
	return &Impl{db: db}
}

func (repo *Impl) GetTwitterAccounts() ([]entities.TwitterAccount, error) {
	var twitterAccounts []entities.TwitterAccount
	response := repo.db.GetDB().Model(&entities.TwitterAccount{}).Find(&twitterAccounts)
	return twitterAccounts, response.Error
}

func (repo *Impl) Save(twitterAccount entities.TwitterAccount) error {
	return repo.db.GetDB().Save(&twitterAccount).Error
}
