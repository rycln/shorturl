package handlers

import (
	"context"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/rycln/shorturl/internal/app/logger"
	"go.uber.org/zap"
)

type pingStorager interface {
	Ping(context.Context) error
}

type Ping struct {
	strg pingStorager
}

func NewPing(strg pingStorager) *Ping {
	return &Ping{
		strg: strg,
	}
}

func (p *Ping) PingDB(c *fiber.Ctx) error {
	err := p.strg.Ping(c.UserContext())
	if err != nil {
		logger.Log.Info("path:"+c.Path()+", "+"func:PingContext()",
			zap.Error(err),
		)
		return c.SendStatus(http.StatusInternalServerError)
	}
	return c.SendStatus(http.StatusOK)
}
