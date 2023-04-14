package twitteraccounts

import (
	"github.com/kaellybot/kaelly-twitter/models/entities"
	"github.com/kaellybot/kaelly-twitter/utils/databases"
)

func New(db databases.MySQLConnection) *TwitterAccountRepositoryImpl {
	return &TwitterAccountRepositoryImpl{db: db}
}

func (repo *TwitterAccountRepositoryImpl) GetTwitterAccounts() ([]entities.TwitterAccount, error) {
	var twitterAccounts []entities.TwitterAccount
	response := repo.db.GetDB().Model(&entities.TwitterAccount{}).Find(&twitterAccounts)
	return twitterAccounts, response.Error
}

func (repo *TwitterAccountRepositoryImpl) Save(twitterAccount entities.TwitterAccount) error {
	return repo.db.GetDB().Save(&twitterAccount).Error
}
