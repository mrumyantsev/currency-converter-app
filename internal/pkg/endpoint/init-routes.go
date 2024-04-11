package endpoint

import "github.com/labstack/echo/v4"

func (e *Endpoint) InitRoutes(echo *echo.Echo) {
	echo.GET("/currencies", e.Currencies)
}
