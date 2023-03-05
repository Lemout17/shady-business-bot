package main

import (
	"fmt"
	"os"

	"github.com/Lawliet18/shady-business-bot/internal/service"
	"github.com/Lawliet18/shady-business-bot/internal/tgbot"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

func main() {
	_ = godotenv.Load()
	log := zerolog.New(zerolog.NewConsoleWriter())

	bot := tgbot.New(log, tgbot.Config{
		// WebhookURL: "",
		Token: os.Getenv("TOKEN"),
	})

	svc := service.New(log, fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT")))

	var eg errgroup.Group

	eg.Go(func() error {
		return bot.Start()
	})

	eg.Go(func() error {
		return svc.Start()
	})

	if err := eg.Wait(); err != nil {
		log.Fatal().Err(err).Msg("something went wrong")
	}

	log.Info().Msg("gracefully exitting")
}
