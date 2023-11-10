package webhooks

import (
	"github.com/gtuk/discordwebhook"
	"github.com/kaellybot/kaelly-twitter/models/constants"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func SendWebhookMessage(content string) {
	webhookUrl := viper.GetString(constants.DiscordWebhookUrl)
	username := constants.ExternalName
	err := discordwebhook.SendMessage(webhookUrl, discordwebhook.Message{
		Username: &username,
		Content:  &content,
	})
	if err != nil {
		log.Error().Err(err).Msgf("An error occurred while sending webhook message, continuing...")
	}
}
