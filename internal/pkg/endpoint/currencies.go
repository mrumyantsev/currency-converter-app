package endpoint

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	memcache "github.com/mrumyantsev/currency-converter-app/internal/pkg/mem-cache"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/service"
	"github.com/mrumyantsev/go-errlib"
	"github.com/rs/zerolog/log"
)

type CurrenciesEndpoint struct {
	config   *config.Config
	memCache *memcache.MemCache
	service  service.Currencies
}

func NewCurrenciesEndpoint(cfg *config.Config, mc *memcache.MemCache, svc service.Currencies) *CurrenciesEndpoint {
	return &CurrenciesEndpoint{
		config:   cfg,
		memCache: mc,
		service:  svc,
	}
}

func (e *CurrenciesEndpoint) Currencies(ctx echo.Context) error {
	calculatedCurrencies := e.memCache.CalculatedCurrencies()

	if err := ctx.JSON(http.StatusOK, calculatedCurrencies); err != nil {
		errMsg := "could not send reponse data"

		log.Error().Err(err).Msg(errMsg)

		return errlib.Wrap(err, errMsg)
	}

	return nil
}
