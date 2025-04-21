package endpoint

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/go-errlib"
	"github.com/rs/zerolog/log"
)

const (
	methodGet       = "GET"
	headerUserAgent = "User-Agent"
)

type CurrenciesFromSourceEndpoint struct {
	config *config.Config
	client *http.Client
}

func NewCurrenciesFromSourceEndpoint(cfg *config.Config) *CurrenciesFromSourceEndpoint {
	return &CurrenciesFromSourceEndpoint{
		config: cfg,
		client: new(http.Client),
	}
}

func (e *CurrenciesFromSourceEndpoint) CurrenciesFromSource() ([]byte, error) {
	startTime := time.Now()

	url, err := url.Parse(e.config.CurrencySourceUrl)
	if err != nil {
		return nil, errlib.Wrap(err, "could not parse url")
	}

	req := e.request(url, methodGet)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, errlib.Wrap(err, "could not send request to server")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errlib.Wrap(err, "could not read data from response body")
	}
	defer func() { _ = resp.Body.Close() }()

	elapsedTime := time.Since(startTime)

	log.Debug().Msg("getting http data time overall: " + elapsedTime.String())

	return data, nil
}

func (e *CurrenciesFromSourceEndpoint) request(url *url.URL, method string) *http.Request {
	log.Debug().Msg(fmt.Sprintf("using %s protocol in request", e.config.HttpRequestProtocol))
	log.Debug().Msg(fmt.Sprintf("using user-agent header: %s", e.config.FakeUserAgentHeaderValue))

	if method == "" {
		method = methodGet
	}

	return &http.Request{
		Method: method,
		URL:    url,
		Proto:  e.config.HttpRequestProtocol,
		Header: map[string][]string{
			headerUserAgent: {e.config.FakeUserAgentHeaderValue},
		},
	}
}
