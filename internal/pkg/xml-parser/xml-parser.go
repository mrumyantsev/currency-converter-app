package xmlparser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/consts"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/models"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib/e"

	"github.com/mrumyantsev/logx/log"
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

func (p *XmlParser) Parse(data []byte) (*models.CurrencyStorage, error) {

	var (
		startTime       time.Time     = time.Now()
		buffer          *bytes.Buffer = bytes.NewBuffer(data)
		decoder         *xml.Decoder  = xml.NewDecoder(buffer)
		currencyStorage *models.CurrencyStorage
		err             error
	)

	decoder.CharsetReader = charset.NewReaderLabel

	if p.config.IsUseMultithreadedParsing {
		log.Debug("using multithreaded parsing")

		currencyStorage, err = p.getParsedDataMultiThreaded(decoder)
		if err != nil {
			return nil, e.Wrap("could not do multithreaded parsing", err)
		}
	} else {
		log.Debug("using singlethreaded parsing")

		currencyStorage, err = p.getParsedDataSingleThreaded(decoder)
		if err != nil {
			return nil, e.Wrap("could not do singlethreaded parsing", err)
		}
	}

	elapsedTime := time.Since(startTime)

	log.Debug(fmt.Sprintf("parsing time overall: %s", elapsedTime))

	return currencyStorage, nil
}

func (p *XmlParser) getParsedDataMultiThreaded(decoder *xml.Decoder) (*models.CurrencyStorage, error) {
	const (
		CURRENCY_START_ELEMENT_NAME string = "Valute"
	)

	var (
		currencyStorage models.CurrencyStorage = models.CurrencyStorage{
			Currencies: make(
				[]models.Currency,
				consts.LENGTH_OF_CURRENCIES_SCLICE_INITIAL,
				consts.CAPACITY_OF_CURRENCIES_SCLICE_INITIAL,
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

			return nil, e.Wrap("could not decode xml element", err)
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

	return &currencyStorage, nil
}

func (p *XmlParser) getParsedDataSingleThreaded(decoder *xml.Decoder) (*models.CurrencyStorage, error) {
	var (
		currencyStorage models.CurrencyStorage = models.CurrencyStorage{
			Currencies: make(
				[]models.Currency,
				consts.LENGTH_OF_CURRENCIES_SCLICE_INITIAL,
				consts.CAPACITY_OF_CURRENCIES_SCLICE_INITIAL,
			),
		}
		err error
	)

	err = decoder.Decode(&currencyStorage)
	if err != nil {
		return nil, e.Wrap("could not decode xml data", err)
	}

	return &currencyStorage, nil
}
