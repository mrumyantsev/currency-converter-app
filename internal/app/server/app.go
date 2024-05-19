package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"

	"errors"
	"strconv"
	"time"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/database"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/endpoint"
	fsops "github.com/mrumyantsev/currency-converter-app/internal/pkg/fs-ops"
	memcache "github.com/mrumyantsev/currency-converter-app/internal/pkg/mem-cache"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/models"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/repository"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/server"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/service"
	timechecks "github.com/mrumyantsev/currency-converter-app/internal/pkg/time-checks"
	xmlparser "github.com/mrumyantsev/currency-converter-app/internal/pkg/xml-parser"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib/errlib"
)

type App struct {
	config     *config.Config
	fsOps      *fsops.FsOps
	xmlParser  *xmlparser.XmlParser
	timeChecks *timechecks.TimeChecks
	memCache   *memcache.MemCache
	database   *database.Database
	service    *service.Service
	endpoint   *endpoint.Endpoint
	server     *server.Server
}

func New() (*App, error) {
	cfg := config.New()

	if err := cfg.Init(); err != nil {
		return nil, errlib.Wrap(err, "could not initialize configuration")
	}

	memCache := memcache.New()

	db := database.New(cfg)

	repository := repository.New(cfg, db)

	service := service.New(cfg, repository)

	endpoint := endpoint.New(cfg, memCache, service)

	mwCors := middleware.CORS()

	server := server.New(cfg, endpoint, mwCors)

	return &App{
		config:     cfg,
		fsOps:      fsops.New(cfg),
		xmlParser:  xmlparser.New(cfg),
		timeChecks: timechecks.New(cfg),
		memCache:   memCache,
		database:   db,
		service:    service,
		endpoint:   endpoint,
		server:     server,
	}, nil
}

func (a *App) Run() error {
	log.Info().Msg("service started")

	err := a.database.Connect()
	if err != nil {
		return errlib.Wrap(err, "could not connect to database")
	}

	log.Debug().Msg("database connection opened")

	goErr := make(chan error, 1)

	isShutdown := false

	go func() {
		if err := a.server.Start(); (err != nil) && !isShutdown {
			goErr <- errlib.Wrap(err, "could not start http server")
		}
	}()

	go func() {
		if err := a.workLoop(); err != nil {
			goErr <- errlib.Wrap(err, "could not proceed work loop")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-goErr:
		return err
	case <-quit:
		break
	}

	// Graceful shutdown

	log.Info().Msg("shutdown signal read")

	isShutdown = true

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err = a.server.Shutdown(ctx); err != nil {
		return errlib.Wrap(err, "could not shutdown http server")
	}

	log.Debug().Msg("http server shut down")

	if err = a.database.Disconnect(); err != nil {
		return errlib.Wrap(err, "could not disconnect from database")
	}

	log.Debug().Msg("database connection closed")

	log.Info().Msg("service gracefully shut down")

	return nil
}

func (a *App) SaveCurrencyDataToFile() error {
	data, err := a.endpoint.CurrenciesFromSource.CurrenciesFromSource()
	if err != nil {
		return errlib.Wrap(err, "could not get currencies from web")
	}

	if err = a.fsOps.OverwriteCurrencyDataFile(data); err != nil {
		return errlib.Wrap(err, "could not write currencies to file")
	}

	log.Info().Msg("currency data saved in file: " + a.config.CurrencySourceFile)

	return nil
}

func (a *App) workLoop() error {
	var (
		timeToNextUpdate time.Duration
		err              error
	)

	for {
		if err = a.updateCurrencyDataInStorages(); err != nil {
			return errlib.Wrap(err, "could not update currency data in storages")
		}

		timeToNextUpdate, err = a.timeChecks.TimeToNextUpdate()
		if err != nil {
			return errlib.Wrap(err, "could not get time to next update")
		}

		log.Info().Msg("next update will occur after " +
			(timeToNextUpdate).Round(time.Second).String())

		if err = a.calculateOutputData(); err != nil {
			return errlib.Wrap(err, "could not calculate output data")
		}

		time.Sleep(timeToNextUpdate)
	}

	return nil
}

func (a *App) updateCurrencyDataInStorages() error {
	currentDatetime := time.Now().Format(time.RFC3339)

	var (
		latestUpdateDatetime models.UpdateDatetime
		latestCurrencies     models.Currencies
		isNeedUpdate         bool
		err                  error
	)

	log.Info().Msg("checking latest update datetime...")

	latestUpdateDatetime, err = a.service.UpdateDatetime.GetLatest()
	if err != nil {
		return errlib.Wrap(err, "could not get current update datetime")
	}

	isNeedUpdate, err = a.timeChecks.IsNeedForUpdateDb(&latestUpdateDatetime)
	if err != nil {
		return errlib.Wrap(err, "could not check is need update for db or not")
	}

	if isNeedUpdate {
		log.Info().Msg("data is outdated")
		log.Info().Msg("initializing update process...")

		if latestCurrencies, err = a.parsedDataFromSource(); err != nil {
			return errlib.Wrap(err, "could not get parsed data from source")
		}

		log.Info().Msg("saving data...")

		latestUpdateDatetime, err = a.service.UpdateDatetime.Create(currentDatetime)
		if err != nil {
			return errlib.Wrap(err, "could not insert datetime into db")
		}

		err = a.service.Currencies.Create(latestCurrencies, latestUpdateDatetime.Id)
		if err != nil {
			return errlib.Wrap(err, "could not insert currencies into db")
		}
	}

	latestCurrencies, err = a.service.Currencies.GetLatest(latestUpdateDatetime.Id)
	if err != nil {
		return errlib.Wrap(err, "could not get currencies from db")
	}

	a.memCache.SetUpdateDatetime(&latestUpdateDatetime)
	a.memCache.SetCurrencies(&latestCurrencies)

	log.Info().Msg("data is now up to date")

	return nil
}

func (a *App) parsedDataFromSource() (models.Currencies, error) {
	var (
		currencies   models.Currencies
		currencyData []byte
		err          error
	)

	log.Info().Msg("getting new data...")

	if a.config.IsReadCurrencyDataFromFile {
		log.Debug().Msg("getting data from local file...")

		if currencyData, err = a.fsOps.CurrencyData(); err != nil {
			return currencies, errlib.Wrap(err, "could not get currencies from file")
		}
	} else {
		log.Debug().Msg("getting data from web...")

		if currencyData, err = a.endpoint.CurrenciesFromSource.CurrenciesFromSource(); err != nil {
			return currencies, errlib.Wrap(err, "could not get curencies from web")
		}
	}

	if err = replaceCommasWithDots(currencyData); err != nil {
		return currencies, errlib.Wrap(err, "could not replace commas in data")
	}

	log.Info().Msg("parsing data...")

	if currencies, err = a.xmlParser.Parse(currencyData); err != nil {
		return currencies, errlib.Wrap(err, "could not parse data")
	}

	return currencies, nil
}

func replaceCommasWithDots(data []byte) error {
	const (
		startDataIndex = 100
		charComma      = ','
		charDot        = '.'
	)

	if data == nil {
		return errors.New("data is empty")
	}

	lengthOfData := len(data)

	for i := startDataIndex; i < lengthOfData; i++ {
		if data[i] == charComma {
			data[i] = charDot
		}
	}

	return nil
}

func (a *App) calculateOutputData() error {
	currencies := a.memCache.Currencies()

	calculatedCurrencies := make(
		[]models.CalculatedCurrency,
		0,
		len(currencies.Currencies),
	)

	var (
		calculatedCurrency models.CalculatedCurrency
		ratio              string
		err                error
	)

	log.Info().Msg("calculate output data...")

	for _, currency := range currencies.Currencies {
		ratio, err = calculateRatio(currency.Value, currency.Multiplier)
		if err != nil {
			return errlib.Wrap(err, "could not calculate currency rate")
		}

		calculatedCurrency.Name = currency.Name
		calculatedCurrency.CharCode = currency.CharCode
		calculatedCurrency.Ratio = ratio

		calculatedCurrencies = append(calculatedCurrencies, calculatedCurrency)
	}

	a.memCache.SetCalculatedCurrencies(calculatedCurrencies)

	return nil
}

func calculateRatio(currencyValue string, currencyMultiplier int) (string, error) {
	const (
		floatBitSize   = 64
		floatFormat    = 'f'
		floatPrecision = -1
	)

	var (
		value      float64
		multiplier float64
		result     float64
		output     string
		err        error
	)

	value, err = strconv.ParseFloat(currencyValue, floatBitSize)
	if err != nil {
		return output, errlib.Wrap(err, "could not parse string to float")
	}

	multiplier = float64(currencyMultiplier)

	result = 1 / (value / multiplier)

	output = strconv.FormatFloat(
		result,
		floatFormat,
		floatPrecision,
		floatBitSize,
	)

	return output, nil
}
