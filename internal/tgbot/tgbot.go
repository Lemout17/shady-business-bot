package tgbot

import (
	"net/url"

	"github.com/NicoNex/echotron/v3"
	"github.com/rs/zerolog"
)

type Bot struct {
	log zerolog.Logger

	webhookURL string
	token      string
}

type Config struct {
	WebhookURL string
	Token      string
}

func New(log zerolog.Logger, config Config) *Bot {
	return &Bot{
		log:        log,
		webhookURL: config.WebhookURL,
		token:      config.Token,
	}
}

func (bot *Bot) Start() error {
	api := echotron.NewAPI(bot.token)

	updates, err := bot.updatesChannel()
	if err != nil {
		return err
	}

	for u := range updates {
		if u.Message.Text == "/start" {
			api.SendMessage("Hello world", u.ChatID(), nil)
		}
	}

	return nil
}

func (bot *Bot) updatesChannel() (<-chan *echotron.Update, error) {
	if bot.webhookURL == "" {
		return echotron.PollingUpdates(bot.token), nil
	}

	u, err := url.ParseRequestURI(bot.webhookURL)
	if err != nil {
		return nil, err
	}
	u.Path = bot.token

	return echotron.WebhookUpdates(u.String(), bot.token), nil
}
