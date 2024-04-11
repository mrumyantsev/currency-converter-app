package endpoint

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib/errlib"
	"github.com/mrumyantsev/logx/log"
)

func (e *Endpoint) Currencies(ctx echo.Context) error {
	calculatedCurrencies := e.memCache.CalculatedCurrencies()

	if err := ctx.JSON(http.StatusOK, calculatedCurrencies); err != nil {
		errMsg := "could not send reponse data"

		log.Error(errMsg, err)

		return errlib.Wrap(errMsg, err)
	}

	return nil
}
