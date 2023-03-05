package service

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type Service struct {
	addr string
}

func New(log zerolog.Logger, addr string) *Service {
	return &Service{
		addr: addr,
	}
}

type requestArgs struct {
	Name  string
	Phone string
	Msg   string
}

func (svc *Service) Start() error {
	e := echo.New()
	e.POST("/", func(c echo.Context) error {
		var args requestArgs
		err := c.Bind(&args)
		if err != nil {
			return fmt.Errorf("bind args: %w", err)
		}

		return c.NoContent(200)
	})

	err := e.Start(svc.addr)
	if err != nil {
		return fmt.Errorf("run http server: %w", err)
	}

	return nil
}
