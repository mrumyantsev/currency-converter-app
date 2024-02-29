package parserd

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	dbstorage "github.com/mrumyantsev/currency-converter-app/internal/pkg/db-storage"
	fsops "github.com/mrumyantsev/currency-converter-app/internal/pkg/fs-ops"
	httpclient "github.com/mrumyantsev/currency-converter-app/internal/pkg/http-client"
	httpserver "github.com/mrumyantsev/currency-converter-app/internal/pkg/http-server"
	memstorage "github.com/mrumyantsev/currency-converter-app/internal/pkg/mem-storage"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/models"
	timechecks "github.com/mrumyantsev/currency-converter-app/internal/pkg/time-checks"
	xmlparser "github.com/mrumyantsev/currency-converter-app/internal/pkg/xml-parser"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib/e"
	"github.com/mrumyantsev/logx"

	"github.com/mrumyantsev/logx/log"
)

type App struct {
	config     *config.Config
	fsOps      *fsops.FsOps
	httpClient *httpclient.HttpClient
	xmlParser  *xmlparser.XmlParser
	timeChecks *timechecks.TimeChecks
	memStorage *memstorage.MemStorage
	dbStorage  *dbstorage.DbStorage
	httpServer *httpserver.HttpServer
}

func New() *App {
	cfg := config.New()

	err := cfg.Init()
	if err != nil {
		log.Error("could not initialize configuration", err)
	}

	log.ApplyConfig(&logx.Config{
		IsDisableDebugLogs: !cfg.IsEnableDebugLogs,
	})

	memStorage := memstorage.New()

	return &App{
		config:     cfg,
		fsOps:      fsops.New(cfg),
		httpClient: httpclient.New(cfg),
		xmlParser:  xmlparser.New(cfg),
		timeChecks: timechecks.New(cfg),
		memStorage: memStorage,
		dbStorage:  dbstorage.New(cfg),
		httpServer: httpserver.New(cfg, memStorage),
	}
}

func (a *App) SaveCurrencyDataToFile() {
	data, err := a.httpClient.GetCurrencyData()
	if err != nil {
		log.Fatal("could not get currencies from web", err)
	}

	err = a.fsOps.OverwriteCurrencyDataFile(data)
	if err != nil {
		log.Fatal("could not write currencies to file", err)
	}

	log.Info("currency data saved in file: " + a.config.CurrencySourceFile)
}

func (a *App) Run() {
	var (
		timeToNextUpdate *time.Duration
		err              error
	)

	for {
		err = a.updateCurrencyDataInStorages()
		if err != nil {
			log.Fatal("could not update currency data in storages", err)
		}

		timeToNextUpdate, err = a.timeChecks.GetTimeToNextUpdate()
		if err != nil {
			log.Fatal("could not get time to next update", err)
		}

		log.Info("next update will occur after " +
			(*timeToNextUpdate).Round(time.Second).String())

		err = a.calculateOutputData()
		if err != nil {
			log.Fatal("could not calculate output data", err)
		}

		if !a.httpServer.GetIsRunning() {
			go func() {
				err = a.httpServer.Run()
				if err != nil {
					log.Fatal("could not run http server", err)
				}
			}()
		}

		time.Sleep(*timeToNextUpdate)
	}
}

func (a *App) updateCurrencyDataInStorages() error {
	var (
		latestUpdateDatetime  *models.UpdateDatetime
		latestCurrencyStorage *models.CurrencyStorage
		isNeedUpdate          bool
		currentDatetime       string = time.Now().Format(time.RFC3339)
		err                   error
	)

	err = a.dbStorage.Connect()
	if err != nil {
		return e.Wrap("could not connect to db to do data update", err)
	}
	defer func() { _ = a.dbStorage.Disconnect() }()

	log.Info("checking latest update datetime...")

	latestUpdateDatetime, err = a.dbStorage.GetLatestUpdateDatetime()
	if err != nil {
		return e.Wrap("could not get current update datetime", err)
	}

	isNeedUpdate, err = a.timeChecks.IsNeedForUpdateDb(latestUpdateDatetime)
	if err != nil {
		return e.Wrap("could not check is need update for db or not", err)
	}

	if isNeedUpdate {
		log.Info("data is outdated")
		log.Info("initializing update process...")

		latestCurrencyStorage, err = a.getParsedDataFromSource()
		if err != nil {
			return e.Wrap("could not get parsed data from source", err)
		}

		log.Info("saving data...")

		latestUpdateDatetime, err = a.dbStorage.InsertUpdateDatetime(currentDatetime)
		if err != nil {
			return e.Wrap("could not insert datetime into db", err)
		}

		err = a.dbStorage.InsertCurrencies(latestCurrencyStorage, latestUpdateDatetime.Id)
		if err != nil {
			return e.Wrap("could not insert currencies into db", err)
		}
	}

	latestCurrencyStorage, err = a.dbStorage.GetLatestCurrencies(latestUpdateDatetime.Id)
	if err != nil {
		return e.Wrap("could not get currencies from db", err)
	}

	a.memStorage.SetUpdateDatetime(latestUpdateDatetime)
	a.memStorage.SetCurrencyStorage(latestCurrencyStorage)

	log.Info("data is now up to date")

	return nil
}

func (a *App) getParsedDataFromSource() (*models.CurrencyStorage, error) {
	var (
		currencyData []byte
		err          error
	)

	log.Info("getting new data...")

	if a.config.IsReadCurrencyDataFromFile {
		log.Debug("getting data from local file...")

		currencyData, err = a.fsOps.GetCurrencyData()
		if err != nil {
			return nil, e.Wrap("could not get currencies from file", err)
		}
	} else {
		log.Debug("getting data from web...")

		currencyData, err = a.httpClient.GetCurrencyData()
		if err != nil {
			return nil, e.Wrap("could not get curencies from web", err)
		}
	}

	err = replaceCommasWithDots(currencyData)
	if err != nil {
		return nil, e.Wrap("could not replace commas in data", err)
	}

	log.Info("parsing data...")

	currencyStorage, err := a.xmlParser.Parse(currencyData)
	if err != nil {
		return nil, e.Wrap("could not parse data", err)
	}

	return currencyStorage, nil
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
	var (
		currencyStorage      *models.CurrencyStorage     = a.memStorage.GetCurrencyStorage()
		calculatedCurrencies []models.CalculatedCurrency = make(
			[]models.CalculatedCurrency,
			0,
			len(currencyStorage.Currencies),
		)
		calculatedCurrency models.CalculatedCurrency
		ratio              *string
		err                error
	)

	log.Info("calculate output data...")

	for _, currency := range currencyStorage.Currencies {
		ratio, err = calculateRatio(&currency.CurrencyValue, &currency.Multiplier)
		if err != nil {
			return e.Wrap("could not calculate currency rate", err)
		}

		calculatedCurrency.Name = currency.Name
		calculatedCurrency.CharCode = currency.CharCode
		calculatedCurrency.Ratio = *ratio

		calculatedCurrencies = append(calculatedCurrencies, calculatedCurrency)
	}

	a.memStorage.SetCalculatedCurrency(calculatedCurrencies)

	return nil
}

func calculateRatio(currencyValue *string, currencyMultiplier *int) (*string, error) {
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

	value, err = strconv.ParseFloat(*currencyValue, floatBitSize)
	if err != nil {
		return nil, e.Wrap("could not parse string to float", err)
	}

	multiplier = float64(*currencyMultiplier)

	result = 1 / (value / multiplier)

	output = strconv.FormatFloat(
		result,
		floatFormat,
		floatPrecision,
		floatBitSize,
	)

	return &output, nil
}

// Prints data. For debugging purposes.
func printData(currencyStorage *models.CurrencyStorage) {
	for _, currency := range currencyStorage.Currencies {
		fmt.Println(currency.NumCode)
		fmt.Println(" ", currency.CharCode)
		fmt.Println(" ", currency.Multiplier)
		fmt.Println(" ", currency.Name)
		fmt.Println(" ", currency.CurrencyValue)
	}
}
