package xmlparser

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"time"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/models"
	"github.com/mrumyantsev/go-errlib"
	"github.com/rs/zerolog/log"

	"golang.org/x/net/html/charset"
)

const (
	firstXmlElement = "Valute"
)

type XmlParser struct {
	config *config.Config
}

func New(cfg *config.Config) *XmlParser {
	return &XmlParser{config: cfg}
}

func (p *XmlParser) Parse(data []byte) (models.Currencies, error) {
	startTime := time.Now()
	buffer := bytes.NewBuffer(data)
	decoder := xml.NewDecoder(buffer)

	var (
		currencies models.Currencies
		err        error
	)

	decoder.CharsetReader = charset.NewReaderLabel

	if p.config.IsUseMultithreadedParsing {
		log.Debug().Msg("using multithreaded parsing")

		currencies, err = p.parsedDataMultiThreaded(decoder)
		if err != nil {
			return currencies, errlib.Wrap(err, "could not do multithreaded parsing")
		}
	} else {
		log.Debug().Msg("using singlethreaded parsing")

		currencies, err = p.parsedDataSingleThreaded(decoder)
		if err != nil {
			return currencies, errlib.Wrap(err, "could not do singlethreaded parsing")
		}
	}

	elapsedTime := time.Since(startTime)

	log.Debug().Msg(fmt.Sprintf("parsing time overall: %s", elapsedTime))

	return currencies, nil
}

func (p *XmlParser) parsedDataMultiThreaded(decoder *xml.Decoder) (models.Currencies, error) {
	currencies := models.Currencies{
		Currencies: make([]models.Currency, 0, p.config.InitialCurrenciesCapacity),
	}

	var (
		currency     models.Currency
		token        xml.Token
		startElement xml.StartElement
		ok           bool
		err          error
	)

	for {
		if token, err = decoder.Token(); err != nil {
			if err == io.EOF {
				break
			}

			return currencies, errlib.Wrap(err, "could not decode xml element")
		}
		if token == nil {
			break
		}

		if startElement, ok = token.(xml.StartElement); !ok {
			continue
		}

		if startElement.Name.Local == firstXmlElement {
			decoder.DecodeElement(&currency, &startElement)

			currencies.Currencies = append(currencies.Currencies, currency)
		}
	}

	return currencies, nil
}

func (p *XmlParser) parsedDataSingleThreaded(decoder *xml.Decoder) (models.Currencies, error) {
	currencies := models.Currencies{
		Currencies: make([]models.Currency, 0, p.config.InitialCurrenciesCapacity),
	}

	if err := decoder.Decode(&currencies); err != nil {
		return currencies, errlib.Wrap(err, "could not decode xml data")
	}

	return currencies, nil
}
