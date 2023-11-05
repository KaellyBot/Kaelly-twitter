package constants

const (
	UrlClassic = "https://twitter.com"
	UrlPreview = "https://vxtwitter.com"
)

type TwitterAccount struct {
	Locale   string
	Username string
}

func GetTwitterAccounts() []TwitterAccount {
	return []TwitterAccount{
		{
			Locale:   "FR",
			Username: "DOFUSfr",
		},
		{
			Locale:   "EN",
			Username: "DOFUS_EN",
		},
		{
			Locale:   "ES",
			Username: "ES_DOFUS",
		},
	}
}
