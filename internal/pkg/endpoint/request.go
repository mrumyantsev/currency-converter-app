package endpoint

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mrumyantsev/currency-converter-app/pkg/lib/errlib"
	"github.com/mrumyantsev/logx/log"
)

const (
	methodGet       = "GET"
	headerUserAgent = "User-Agent"
)

func (e *Endpoint) CurrenciesFromSource() ([]byte, error) {
	startTime := time.Now()

	url, err := url.Parse(e.config.CurrencySourceUrl)
	if err != nil {
		return nil, errlib.Wrap("could not parse url", err)
	}

	req := e.request(url, methodGet)

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, errlib.Wrap("could not send request to server", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errlib.Wrap("could not read data from response body", err)
	}
	defer func() { _ = resp.Body.Close() }()

	elapsedTime := time.Since(startTime)

	log.Debug("getting http data time overall: " + elapsedTime.String())

	return data, nil
}

func (e *Endpoint) request(url *url.URL, method string) *http.Request {
	log.Debug(fmt.Sprintf("using %s protocol in request", e.config.HttpRequestProtocol))
	log.Debug(fmt.Sprintf("using user-agent header: %s", e.config.FakeUserAgentHeaderValue))

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
