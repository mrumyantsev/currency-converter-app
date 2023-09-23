package xmlparser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/mrumyantsev/currency-converter/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter/internal/pkg/consts"
	"github.com/mrumyantsev/currency-converter/internal/pkg/models"

	"github.com/mrumyantsev/fastlog"
	"golang.org/x/net/html/charset"
)

type XmlParser struct {
	config *config.Config
}

func New(cfg *config.Config) *XmlParser {
	return &XmlParser{
		config: cfg,
	}
}

func (p *XmlParser) Parse(data []byte) *models.CurrencyStorage {
	fastlog.Debug("begin parsing data...")

	var (
		startTime       time.Time     = time.Now()
		buffer          *bytes.Buffer = bytes.NewBuffer(data)
		decoder         *xml.Decoder  = xml.NewDecoder(buffer)
		currencyStorage *models.CurrencyStorage
	)

	decoder.CharsetReader = charset.NewReaderLabel

	if p.config.IsUseMultithreadedParsing {
		fastlog.Debug("using multithreaded parsing")
		currencyStorage = p.getParsedDataMultiThreaded(decoder)
	} else {
		fastlog.Debug("using singlethreaded parsing")
		currencyStorage = p.getParsedDataSingleThreaded(decoder)
	}

	elapsedTime := time.Since(startTime)

	fastlog.Debug(fmt.Sprintf("parsing time overall: %s", elapsedTime))

	return currencyStorage
}

func (p *XmlParser) getParsedDataMultiThreaded(decoder *xml.Decoder) *models.CurrencyStorage {
	const (
		CURRENCY_START_ELEMENT_NAME string = "Valute"
	)

	var (
		currencyStorage models.CurrencyStorage = models.CurrencyStorage{
			Currencies: make(
				[]models.Currency,
				consts.LEN_OF_CURRENCIES_SCLICE_INITIAL,
				consts.CAP_OF_CURRENCIES_SCLICE_INITIAL,
			),
		}
		currency     models.Currency
		token        xml.Token
		startElement xml.StartElement
		ok           bool
		err          error
	)

	for {
		token, err = decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}

			fastlog.Fatal("cannot decode xml element", err)
		}

		if token == nil {
			break
		}

		startElement, ok = token.(xml.StartElement)
		if !ok {
			continue
		}

		if startElement.Name.Local == CURRENCY_START_ELEMENT_NAME {
			decoder.DecodeElement(&currency, &startElement)
			currencyStorage.Currencies = append(currencyStorage.Currencies, currency)
		}
	}

	return &currencyStorage
}

func (p *XmlParser) getParsedDataSingleThreaded(decoder *xml.Decoder) *models.CurrencyStorage {
	var (
		currencyStorage models.CurrencyStorage = models.CurrencyStorage{
			Currencies: make(
				[]models.Currency,
				consts.LEN_OF_CURRENCIES_SCLICE_INITIAL,
				consts.CAP_OF_CURRENCIES_SCLICE_INITIAL,
			),
		}
		err error
	)

	err = decoder.Decode(&currencyStorage)
	if err != nil {
		fastlog.Fatal("cannot decode xml data", err)
	}

	return &currencyStorage
}
