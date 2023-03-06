package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Lawliet18/shady-business-bot/internal/message"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Service struct {
	addr    string
	log     zerolog.Logger
	msgChan chan<- message.Message
}

func New(log zerolog.Logger, addr string, msgChan chan<- message.Message) *Service {
	return &Service{
		addr:    addr,
		log:     log,
		msgChan: msgChan,
	}
}

type requestArgs struct {
	Name  string
	Phone string
}

func (svc *Service) Start(ctx context.Context) error {
	e := echo.New()
	e.Any("/", func(c echo.Context) error {
		var args requestArgs
		err := c.Bind(&args)
		if err != nil {
			return fmt.Errorf("bind args: %w", err)
		}

		svc.msgChan <- message.Message{
			Name:  args.Name,
			Phone: args.Phone,
		}

		return c.NoContent(200)
	})

	go func() {
		<-ctx.Done()

		cancelCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := e.Shutdown(cancelCtx)
		if err != nil {
			svc.log.Err(err).Msg("shutdown server")
		}
	}()

	err := e.Start(svc.addr)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			return nil
		default:
			return fmt.Errorf("run http server: %w", err)
		}
	}

	return nil
}
