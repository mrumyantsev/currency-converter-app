package parser

import (
	"errors"
	"fmt"
	"time"

	"github.com/mrumyantsev/currency-converter/internal/pkg/config"
	fsops "github.com/mrumyantsev/currency-converter/internal/pkg/fs-ops"
	httpclient "github.com/mrumyantsev/currency-converter/internal/pkg/http-client"
	"github.com/mrumyantsev/currency-converter/internal/pkg/models"
	"github.com/mrumyantsev/currency-converter/internal/pkg/storage"
	"github.com/mrumyantsev/currency-converter/internal/pkg/utils"
	xmlparser "github.com/mrumyantsev/currency-converter/internal/pkg/xml-parser"

	"github.com/mrumyantsev/fastlog"
)

type Parser struct {
	config     *config.Config
	fsOps      *fsops.FsOps
	httpClient *httpclient.HttpClient
	xmlParser  *xmlparser.XmlParser
	storage    *storage.Storage
}

func New() *Parser {
	cfg := config.New()

	cfg.Init()

	fastlog.IsEnableDebugLogs = cfg.IsEnableDebugLogs

	return &Parser{
		config:     cfg,
		fsOps:      fsops.New(cfg),
		httpClient: httpclient.New(cfg),
		xmlParser:  xmlparser.New(cfg),
		storage:    storage.New(cfg),
	}
}

func (p *Parser) Run() {
	var data []byte

	if p.config.IsReadCurrencyDataFromFile {
		data = p.fsOps.GetCurrencyData()
	} else {
		data = p.httpClient.GetCurrencyData()
	}

	err := replaceCommasWithDots(data)
	if err != nil {
		fastlog.Error("cannot replace commas in data", err)
	}

	currencies := p.xmlParser.Parse(data)

	showOnScreen(currencies)
}

func replaceCommasWithDots(data []byte) error {
	fastlog.Debug("replacing commas with dots in data...")

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

// Shows data on screen.
func showOnScreen(currencyStorage *models.CurrencyStorage) {
	for _, currency := range currencyStorage.Currencies {
		// break
		fmt.Println(currency.NumCode, fmt.Sprintf("%+T", currency.NumCode))
		fmt.Println(" ", currency.CharCode, fmt.Sprintf("%+T", currency.CharCode))
		fmt.Println(" ", currency.Multiplier, fmt.Sprintf("%+T", currency.Multiplier))
		fmt.Println(" ", currency.Name, fmt.Sprintf("%+T", currency.Name))
		fmt.Println(" ", currency.CurrencyValue, fmt.Sprintf("%+T", currency.CurrencyValue))
	}
}

func (p *Parser) SaveCurrencyDataFile() {
	fastlog.Info("currency data saved on file: " + p.config.CurrencySourceFile)

	data := p.httpClient.GetCurrencyData()

	p.fsOps.OverwriteCurrencyDataFile(data)
}

func (p *Parser) DoWithDB() {
	err := p.storage.Connect()
	if err != nil {
		fastlog.Error("cannot connect to db", err)
	}
	defer func() {
		err = p.storage.Disconnect()
		if err != nil {
			fastlog.Error("cannot disconnect from db", err)
		}
	}()

	currencyStorage, err := p.storage.GetCurrencies(1)
	if err != nil {
		fastlog.Error("cannot get currencies with desired id", err)
	}

	showOnScreen(currencyStorage)
}

func (p *Parser) isNeedForUpdateDb() bool {
	err := p.storage.Connect()
	if err != nil {
		fastlog.Error("cannot connect to db", err)
	}
	defer func() {
		err = p.storage.Disconnect()
		if err != nil {
			fastlog.Error("cannot disconnect from db", err)
		}
	}()

	dbRes, err := p.storage.GetLastDatetime()
	if err != nil {
		fastlog.Error("cannot get last update datetime from db", err)
	}

	todayUpdateDatetime, err := p.getTodayUpdateDatetime()
	if err != nil {
		fastlog.Error("cannot get today update datetime", err)
	}

	lastUpdateDatetime, err := time.Parse(time.RFC3339, dbRes.UpdateDatetime)
	if err != nil {
		fastlog.Error("cannot parse update time from db", err)
	}

	if lastUpdateDatetime.Before(*todayUpdateDatetime) {
		return true
	}

	return false
}

func (p *Parser) getTodayUpdateDatetime() (*time.Time, error) {
	updateTime, err := time.Parse("15:04:05", p.config.TimeWhenNeedToUpdateCurrency)
	if err != nil {
		return nil, utils.DecorateError("cannot parse update time from config", err)
	}

	todayYear, todayMonth, todayDay := time.Now().Date()

	todayUpdateDatetime := time.Date(
		todayYear,
		todayMonth,
		todayDay,
		updateTime.Hour(),
		updateTime.Minute(),
		updateTime.Second(),
		0, // zero nanoseconds
		time.Now().Location(),
	)

	return &todayUpdateDatetime, nil
}
