package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Lawliet18/shady-business-bot/internal/message"
	"github.com/Lawliet18/shady-business-bot/internal/service"
	"github.com/Lawliet18/shady-business-bot/internal/tgbot"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

const (
	ShadyBusinessChatID = -1001820130859
)

func main() {
	_ = godotenv.Load()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	token := os.Getenv("TOKEN")

	msgChan := make(chan message.Message, 1)

	bot := tgbot.New(log, tgbot.Config{
		Token:            token,
		ChatID:           ShadyBusinessChatID,
		NotificationChan: msgChan,
	})

	addr := fmt.Sprintf("0.0.0.0:%s", os.Getenv("PORT"))
	svc := service.New(log, addr, msgChan)

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Info().Msg("gracefully shutting down")
		cancel()
	}()

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		return bot.Start(ctx)
	})

	eg.Go(func() error {
		return svc.Start(ctx)
	})

	if err := eg.Wait(); err != nil {
		log.Fatal().Err(err).Msg("something went wrong")
		return
	}

	log.Info().Msg("shut down was gracefull")
}
