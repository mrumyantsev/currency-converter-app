package server

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mrumyantsev/currency-converter/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter/internal/pkg/consts"
	dbstorage "github.com/mrumyantsev/currency-converter/internal/pkg/db-storage"
	fsops "github.com/mrumyantsev/currency-converter/internal/pkg/fs-ops"
	httpclient "github.com/mrumyantsev/currency-converter/internal/pkg/http-client"
	httpserver "github.com/mrumyantsev/currency-converter/internal/pkg/http-server"
	memstorage "github.com/mrumyantsev/currency-converter/internal/pkg/mem-storage"
	"github.com/mrumyantsev/currency-converter/internal/pkg/models"
	timechecks "github.com/mrumyantsev/currency-converter/internal/pkg/time-checks"
	"github.com/mrumyantsev/currency-converter/internal/pkg/utils"
	xmlparser "github.com/mrumyantsev/currency-converter/internal/pkg/xml-parser"

	"github.com/mrumyantsev/multilog"
	"github.com/mrumyantsev/multilog/log"
)

type Server struct {
	config     *config.Config
	fsOps      *fsops.FsOps
	httpClient *httpclient.HttpClient
	xmlParser  *xmlparser.XmlParser
	timeChecks *timechecks.TimeChecks
	memStorage *memstorage.MemStorage
	dbStorage  *dbstorage.DbStorage
	httpServer *httpserver.HttpServer
}

func New() *Server {
	cfg := config.New()

	err := cfg.Init()
	if err != nil {
		log.Error("cannot initialize configuration", err)
	}

	logCfg := &multilog.Config{
		IsDisableDebugLogs: !cfg.IsEnableDebugLogs,
		TimeFormat:         time.RFC1123Z,
	}

	log.ApplyConfig(logCfg)

	memStorage := memstorage.New()

	return &Server{
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

func (s *Server) SaveCurrencyDataToFile() {
	data, err := s.httpClient.GetCurrencyData()
	if err != nil {
		log.Error("cannot get currencies from web", err)
	}

	err = s.fsOps.OverwriteCurrencyDataFile(data)
	if err != nil {
		log.Error("cannot write currencies to file", err)
	}

	log.Info("currency data saved in file: " + s.config.CurrencySourceFile)
}

func (s *Server) Run() {
	var (
		timeToNextUpdate *time.Duration
		err              error
	)

	for {
		err = s.updateCurrencyDataInStorages()
		if err != nil {
			log.Error("cannot update currency data in storages", err)
		}

		timeToNextUpdate, err = s.timeChecks.GetTimeToNextUpdate()
		if err != nil {
			log.Error("cannot get time to next update", err)
		}

		log.Info("next update will occur after " +
			(*timeToNextUpdate).Round(time.Second).String())

		err = s.calculateOutputData()
		if err != nil {
			log.Error("cannot calculate output data", err)
		}

		if !s.httpServer.GetIsRunning() {
			go func() {
				err = s.httpServer.Run()
				if err != nil {
					log.Error("cannot run http server", err)
				}
			}()
		}

		time.Sleep(*timeToNextUpdate)
	}
}

func (s *Server) updateCurrencyDataInStorages() error {
	var (
		latestUpdateDatetime  *models.UpdateDatetime
		latestCurrencyStorage *models.CurrencyStorage
		isNeedUpdate          bool
		currentDatetime       string = time.Now().Format(time.RFC3339)
		err                   error
	)

	err = s.dbStorage.Connect()
	if err != nil {
		return utils.DecorateError("cannot connect to db to do data update", err)
	}

	log.Info("checking latest update datetime...")

	latestUpdateDatetime, err = s.dbStorage.GetLatestUpdateDatetime()
	if err != nil {
		return utils.DecorateError("cannot get current update datetime", err)
	}

	isNeedUpdate, err = s.timeChecks.IsNeedForUpdateDb(latestUpdateDatetime)
	if err != nil {
		return utils.DecorateError("cannot check is need update for db or not", err)
	}

	if isNeedUpdate {
		log.Info("data is outdated")
		log.Info("initializing update process...")

		latestCurrencyStorage, err = s.getParsedDataFromSource()
		if err != nil {
			return utils.DecorateError("cannot get parsed data from source", err)
		}

		log.Info("saving data...")

		latestUpdateDatetime, err = s.dbStorage.InsertUpdateDatetime(currentDatetime)
		if err != nil {
			return utils.DecorateError("cannot insert datetime into db", err)
		}

		err = s.dbStorage.InsertCurrencies(latestCurrencyStorage, latestUpdateDatetime.Id)
		if err != nil {
			return utils.DecorateError("cannot insert currencies into db", err)
		}
	}

	latestCurrencyStorage, err = s.dbStorage.GetLatestCurrencies(latestUpdateDatetime.Id)
	if err != nil {
		return utils.DecorateError("cannot get currencies from db", err)
	}

	s.memStorage.SetUpdateDatetime(latestUpdateDatetime)
	s.memStorage.SetCurrencyStorage(latestCurrencyStorage)

	log.Info("data is now up to date")

	err = s.dbStorage.Disconnect()
	if err != nil {
		return utils.DecorateError("cannot disconnect from db to do data update", err)
	}

	return nil
}

func (s *Server) getParsedDataFromSource() (*models.CurrencyStorage, error) {
	var (
		currencyData []byte
		err          error
	)

	log.Info("getting new data...")

	if s.config.IsReadCurrencyDataFromFile {
		log.Debug("getting data from local file...")

		currencyData, err = s.fsOps.GetCurrencyData()
		if err != nil {
			return nil, utils.DecorateError("cannot get currencies from file", err)
		}
	} else {
		log.Debug("getting data from web...")

		currencyData, err = s.httpClient.GetCurrencyData()
		if err != nil {
			return nil, utils.DecorateError("cannot get curencies from web", err)
		}
	}

	err = replaceCommasWithDots(currencyData)
	if err != nil {
		return nil, utils.DecorateError("cannot replace commas in data", err)
	}

	log.Info("parsing data...")

	currencyStorage, err := s.xmlParser.Parse(currencyData)
	if err != nil {
		return nil, utils.DecorateError("cannot parse data", err)
	}

	return currencyStorage, nil
}

func replaceCommasWithDots(data []byte) error {
	const (
		START_DATA_INDEX int  = 100
		CHAR_COMMA       byte = ','
		CHAR_DOT         byte = '.'
	)

	if data == nil {
		return errors.New("data is empty")
	}

	lengthOfData := len(data)

	for i := START_DATA_INDEX; i < lengthOfData; i++ {
		if data[i] == CHAR_COMMA {
			data[i] = CHAR_DOT
		}
	}

	return nil
}

func (s *Server) calculateOutputData() error {
	var (
		currencyStorage      *models.CurrencyStorage     = s.memStorage.GetCurrencyStorage()
		calculatedCurrencies []models.CalculatedCurrency = make(
			[]models.CalculatedCurrency,
			consts.LENGTH_OF_CURRENCIES_SCLICE_INITIAL,
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
			return utils.DecorateError("cannot calculate currency rate", err)
		}

		calculatedCurrency.Name = currency.Name
		calculatedCurrency.CharCode = currency.CharCode
		calculatedCurrency.Ratio = *ratio

		calculatedCurrencies = append(calculatedCurrencies, calculatedCurrency)
	}

	s.memStorage.SetCalculatedCurrency(calculatedCurrencies)

	return nil
}

func calculateRatio(currencyValue *string, currencyMultiplier *int) (*string, error) {
	const (
		FLOAT_BIT_SIZE          int  = 64
		FLOAT_COMMON_FORMAT     byte = 'f'
		FLOAT_MAXIMUM_PRECISION int  = -1
	)

	var (
		value      float64
		multiplier float64
		result     float64
		output     string
		err        error
	)

	value, err = strconv.ParseFloat(*currencyValue, FLOAT_BIT_SIZE)
	if err != nil {
		return nil, utils.DecorateError("cannot parse string to float", err)
	}

	multiplier = float64(*currencyMultiplier)

	result = 1 / (value / multiplier)

	output = strconv.FormatFloat(
		result,
		FLOAT_COMMON_FORMAT,
		FLOAT_MAXIMUM_PRECISION,
		FLOAT_BIT_SIZE,
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
