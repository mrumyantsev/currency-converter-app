package server

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4/middleware"

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
	"github.com/mrumyantsev/logx"

	"github.com/mrumyantsev/logx/log"
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

func New() *App {
	cfg := config.New()

	if err := cfg.Init(); err != nil {
		log.Error("could not initialize configuration", err)
	}

	log.ApplyConfig(&logx.Config{
		IsDisableDebugLogs: !cfg.IsEnableDebugLogs,
	})

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
	}
}

func (a *App) Run() {
	if isUserWantSave() {
		a.SaveCurrencyDataToFile()

		return
	}

	log.Info("service started")

	err := a.database.Connect()
	if err != nil {
		log.Fatal("could not connect to database", err)
	}

	log.Debug("database connection opened")

	isShutdown := false

	go func() {
		if err = a.server.Start(); (err != nil) && !isShutdown {
			log.Fatal("could not start http server", err)
		}
	}()

	go a.workLoop()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit // Waiting for a signal to the graceful shutdown

	log.Info("shutdown signal read")

	isShutdown = true

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err = a.server.Shutdown(ctx); err != nil {
		log.Fatal("could not shutdown http server", err)
	}

	log.Debug("http server shut down")

	if err = a.database.Disconnect(); err != nil {
		log.Fatal("could not disconnect from database", err)
	}

	log.Debug("database connection closed")

	log.Info("service gracefully shut down")
}

func (a *App) SaveCurrencyDataToFile() {
	data, err := a.endpoint.CurrenciesFromSource()
	if err != nil {
		log.Fatal("could not get currencies from web", err)
	}

	if err = a.fsOps.OverwriteCurrencyDataFile(data); err != nil {
		log.Fatal("could not write currencies to file", err)
	}

	log.Info("currency data saved in file: " + a.config.CurrencySourceFile)
}

func (a *App) workLoop() {
	var (
		timeToNextUpdate time.Duration
		err              error
	)

	for {
		if err = a.updateCurrencyDataInStorages(); err != nil {
			log.Fatal("could not update currency data in storages", err)
		}

		timeToNextUpdate, err = a.timeChecks.TimeToNextUpdate()
		if err != nil {
			log.Fatal("could not get time to next update", err)
		}

		log.Info("next update will occur after " +
			(timeToNextUpdate).Round(time.Second).String())

		if err = a.calculateOutputData(); err != nil {
			log.Fatal("could not calculate output data", err)
		}

		time.Sleep(timeToNextUpdate)
	}
}

func (a *App) updateCurrencyDataInStorages() error {
	currentDatetime := time.Now().Format(time.RFC3339)

	var (
		latestUpdateDatetime models.UpdateDatetime
		latestCurrencies     models.Currencies
		isNeedUpdate         bool
		err                  error
	)

	log.Info("checking latest update datetime...")

	latestUpdateDatetime, err = a.service.UpdateDatetime.GetLatest()
	if err != nil {
		return errlib.Wrap("could not get current update datetime", err)
	}

	isNeedUpdate, err = a.timeChecks.IsNeedForUpdateDb(&latestUpdateDatetime)
	if err != nil {
		return errlib.Wrap("could not check is need update for db or not", err)
	}

	if isNeedUpdate {
		log.Info("data is outdated")
		log.Info("initializing update process...")

		if latestCurrencies, err = a.parsedDataFromSource(); err != nil {
			return errlib.Wrap("could not get parsed data from source", err)
		}

		log.Info("saving data...")

		latestUpdateDatetime, err = a.service.UpdateDatetime.Create(currentDatetime)
		if err != nil {
			return errlib.Wrap("could not insert datetime into db", err)
		}

		err = a.service.Currencies.Create(latestCurrencies, latestUpdateDatetime.Id)
		if err != nil {
			return errlib.Wrap("could not insert currencies into db", err)
		}
	}

	latestCurrencies, err = a.service.Currencies.GetLatest(latestUpdateDatetime.Id)
	if err != nil {
		return errlib.Wrap("could not get currencies from db", err)
	}

	a.memCache.SetUpdateDatetime(&latestUpdateDatetime)
	a.memCache.SetCurrencies(&latestCurrencies)

	log.Info("data is now up to date")

	return nil
}

func (a *App) parsedDataFromSource() (models.Currencies, error) {
	var (
		currencies   models.Currencies
		currencyData []byte
		err          error
	)

	log.Info("getting new data...")

	if a.config.IsReadCurrencyDataFromFile {
		log.Debug("getting data from local file...")

		if currencyData, err = a.fsOps.CurrencyData(); err != nil {
			return currencies, errlib.Wrap("could not get currencies from file", err)
		}
	} else {
		log.Debug("getting data from web...")

		if currencyData, err = a.endpoint.CurrenciesFromSource(); err != nil {
			return currencies, errlib.Wrap("could not get curencies from web", err)
		}
	}

	if err = replaceCommasWithDots(currencyData); err != nil {
		return currencies, errlib.Wrap("could not replace commas in data", err)
	}

	log.Info("parsing data...")

	if currencies, err = a.xmlParser.Parse(currencyData); err != nil {
		return currencies, errlib.Wrap("could not parse data", err)
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

	log.Info("calculate output data...")

	for _, currency := range currencies.Currencies {
		ratio, err = calculateRatio(currency.Value, currency.Multiplier)
		if err != nil {
			return errlib.Wrap("could not calculate currency rate", err)
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
		return output, errlib.Wrap("could not parse string to float", err)
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

func isUserWantSave() bool {
	f := flag.Bool("s", false, "Save currency data to a local file")

	flag.Parse()

	return *f
}
