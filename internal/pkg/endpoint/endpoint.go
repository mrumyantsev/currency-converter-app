package endpoint

import (
	"github.com/labstack/echo/v4"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	memcache "github.com/mrumyantsev/currency-converter-app/internal/pkg/mem-cache"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/service"
)

type CurrenciesFromSource interface {
	CurrenciesFromSource() ([]byte, error)
}

type Currencies interface {
	Currencies(ctx echo.Context) error
}

type Endpoint struct {
	CurrenciesFromSource CurrenciesFromSource
	Currencies           Currencies
}

func New(cfg *config.Config, mc *memcache.MemCache, svc *service.Service) *Endpoint {
	return &Endpoint{
		CurrenciesFromSource: NewCurrenciesFromSourceEndpoint(cfg),
		Currencies:           NewCurrenciesEndpoint(cfg, mc, svc.Currencies),
	}
}

func (e *Endpoint) InitRoutes(echo *echo.Echo) {
	echo.GET("/currencies", e.Currencies.Currencies)
}
