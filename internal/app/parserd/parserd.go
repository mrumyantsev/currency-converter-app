package parserd

import (
	"errors"
	"fmt"
	"time"

	"github.com/mrumyantsev/currency-converter/internal/pkg/config"
	dbstorage "github.com/mrumyantsev/currency-converter/internal/pkg/db-storage"
	fsops "github.com/mrumyantsev/currency-converter/internal/pkg/fs-ops"
	httpclient "github.com/mrumyantsev/currency-converter/internal/pkg/http-client"
	httpserver "github.com/mrumyantsev/currency-converter/internal/pkg/http-server"
	memstorage "github.com/mrumyantsev/currency-converter/internal/pkg/mem-storage"
	"github.com/mrumyantsev/currency-converter/internal/pkg/models"
	timechecks "github.com/mrumyantsev/currency-converter/internal/pkg/time-checks"
	"github.com/mrumyantsev/currency-converter/internal/pkg/utils"
	xmlparser "github.com/mrumyantsev/currency-converter/internal/pkg/xml-parser"

	"github.com/mrumyantsev/fastlog"
)

type ParserD struct {
	config     *config.Config
	fsOps      *fsops.FsOps
	httpClient *httpclient.HttpClient
	xmlParser  *xmlparser.XmlParser
	timeChecks *timechecks.TimeChecks
	memStorage *memstorage.MemStorage
	dbStorage  *dbstorage.DbStorage
	httpServer *httpserver.HttpServer
}

func New() *ParserD {
	cfg := config.New()

	cfg.Init()

	fastlog.IsEnableDebugLogs = cfg.IsEnableDebugLogs

	memStorage := memstorage.New()

	return &ParserD{
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

func (p *ParserD) SaveCurrencyDataToFile() {
	data, err := p.httpClient.GetCurrencyData()
	if err != nil {
		fastlog.Error("cannot get currencies from web", err)
	}

	err = p.fsOps.OverwriteCurrencyDataFile(data)
	if err != nil {
		fastlog.Error("cannot write currencies to file", err)
	}

	fastlog.Info("currency data saved in file: " + p.config.CurrencySourceFile)
}

func (p *ParserD) Run() {
	var (
		timeToNextUpdate *time.Duration
		err              error
	)

	for {
		p.updateCurrencyDataInStorages()

		timeToNextUpdate = p.timeChecks.GetTimeToNextUpdate()

		fastlog.Info("next update will occur after " +
			(*timeToNextUpdate).Round(time.Second).String())

		if !p.httpServer.GetIsRunning() {
			go func() {
				err = p.httpServer.Run()
				if err != nil {
					fastlog.Error("cannot run http server", err)
				}
			}()
		}

		time.Sleep(*timeToNextUpdate)
	}
}

func (p *ParserD) updateCurrencyDataInStorages() {
	var (
		latestUpdateDatetime  *models.UpdateDatetime
		latestCurrencyStorage *models.CurrencyStorage
		isNeedUpdate          bool
		currentDatetime       string = time.Now().Format(time.RFC3339)
		err                   error
	)

	latestUpdateDatetime, err = p.getLatestUpdateDatetime()
	if err != nil {
		fastlog.Error("cannot get latest update datetime", err)
	}

	isNeedUpdate, err = p.timeChecks.IsNeedForUpdateDb(latestUpdateDatetime)
	if err != nil {
		fastlog.Error("cannot check is need update for db or not", err)
	}

	if isNeedUpdate {
		fastlog.Info("stored currency data is outdated. update is needed")

		latestCurrencyStorage, err = p.getParsedDataFromSource()
		if err != nil {
			fastlog.Error("cannot get parsed data from source", err)
		}

		latestUpdateDatetime, err = p.putDatetimeInDb(currentDatetime)
		if err != nil {
			fastlog.Error("cannot put datetime in db", err)
		}

		err = p.putCurrenciesToDb(latestCurrencyStorage, latestUpdateDatetime)
		if err != nil {
			fastlog.Error("cannot put currencies in db", err)
		}
	} else {
		fastlog.Info("stored currency data is up to date. no need for update")

		latestCurrencyStorage, err = p.getCurrenciesFromDb(latestUpdateDatetime)
		if err != nil {
			fastlog.Error("cannot get data from db", err)
		}
	}

	p.memStorage.SetUpdateDatetime(latestUpdateDatetime)
	p.memStorage.SetCurrencyStorage(latestCurrencyStorage)
}

func (p *ParserD) getLatestUpdateDatetime() (*models.UpdateDatetime, error) {
	fastlog.Info("getting latest update datetime from db...")

	p.connectToDb()
	defer p.disconnectFromDb()

	latestUpdateDatetime, err := p.dbStorage.GetLastDatetime()
	if err != nil {
		return nil, utils.DecorateError("cannot get current update datetime", err)
	}

	return latestUpdateDatetime, err
}

func (p *ParserD) getParsedDataFromSource() (*models.CurrencyStorage, error) {
	var (
		data []byte
		err  error
	)

	fastlog.Info("initiating parsing of xml currency data...")

	if p.config.IsReadCurrencyDataFromFile {
		fastlog.Info("getting data from local file...")

		data, err = p.fsOps.GetCurrencyData()
		if err != nil {
			return nil, utils.DecorateError("cannot get currencies from file", err)
		}
	} else {
		fastlog.Info("getting data from web...")

		data, err = p.httpClient.GetCurrencyData()
		if err != nil {
			return nil, utils.DecorateError("cannot get curencies from web", err)
		}
	}

	err = replaceCommasWithDots(data)
	if err != nil {
		return nil, utils.DecorateError("cannot replace commas in data", err)
	}

	currencyStorage, err := p.xmlParser.Parse(data)
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

func (p *ParserD) putDatetimeInDb(datetime string) (*models.UpdateDatetime, error) {
	p.connectToDb()
	defer p.disconnectFromDb()

	updateDatetime, err := p.dbStorage.InsertDatetime(datetime)
	if err != nil {
		return nil, utils.DecorateError("cannot insert datetime into db", err)
	}

	return updateDatetime, nil
}

func (p *ParserD) putCurrenciesToDb(currencyStorage *models.CurrencyStorage, updateDatetime *models.UpdateDatetime) error {
	fastlog.Info("putting currency data to db...")

	p.connectToDb()
	defer p.disconnectFromDb()

	err := p.dbStorage.InsertCurrencies(
		currencyStorage,
		updateDatetime,
	)
	if err != nil {
		return utils.DecorateError("cannot insert currencies into db", err)
	}

	return nil
}

func (p *ParserD) getCurrenciesFromDb(updateDatetime *models.UpdateDatetime) (*models.CurrencyStorage, error) {
	fastlog.Info("getting latest currency data from db...")

	p.connectToDb()
	defer p.disconnectFromDb()

	latestCurrencyStorage, err := p.dbStorage.GetCurrencies(
		updateDatetime.Id)
	if err != nil {
		return nil, utils.DecorateError("cannot get currencies from db", err)
	}

	return latestCurrencyStorage, err
}

func (p *ParserD) connectToDb() {
	err := p.dbStorage.Connect()
	if err != nil {
		fastlog.Error("cannot connect to db", err)
	}
}

func (p *ParserD) disconnectFromDb() {
	err := p.dbStorage.Disconnect()
	if err != nil {
		fastlog.Error("cannot disconnect from db", err)
	}
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
