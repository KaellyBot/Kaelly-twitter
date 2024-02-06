package constants

const (
	UrlClassic = "https://twitter.com"
	UrlPreview = "https://vxtwitter.com"
)

type TwitterAccount struct {
	Username string
}

func GetTwitterAccount() TwitterAccount {
	return TwitterAccount{Username: "BungieHelp"}
}
