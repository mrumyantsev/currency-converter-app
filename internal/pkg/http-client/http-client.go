package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib/e"

	"github.com/mrumyantsev/logx/log"
)

const (
	headerUserAgent = "User-Agent"
	methodGet       = "GET"
)

type HttpClient struct {
	client *http.Client
	config *config.Config
}

func New(cfg *config.Config) *HttpClient {
	return &HttpClient{client: &http.Client{}, config: cfg}
}

func (c *HttpClient) CurrencyData() ([]byte, error) {
	startTime := time.Now()

	url, err := url.Parse(c.config.CurrencySourceUrl)
	if err != nil {
		return nil, e.Wrap("could not parse url", err)
	}

	requ := c.createRequest(url)

	resp, err := c.client.Do(requ)
	if err != nil {
		return nil, e.Wrap("could not send request to server", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, e.Wrap("could not read data from response body", err)
	}
	defer func() { _ = resp.Body.Close() }()

	elapsedTime := time.Since(startTime)

	log.Debug(fmt.Sprintf("getting http data time overall: %s", elapsedTime))

	return data, nil
}

func (c *HttpClient) createRequest(url *url.URL) *http.Request {
	log.Debug(fmt.Sprintf("using %s protocol in request", c.config.HttpRequestProtocol))
	log.Debug(fmt.Sprintf("using user-agent header: %s", c.config.FakeUserAgentHeaderValue))

	return &http.Request{
		Method: methodGet,
		URL:    url,
		Proto:  c.config.HttpRequestProtocol,
		Header: map[string][]string{
			headerUserAgent: {c.config.FakeUserAgentHeaderValue},
		},
	}
}
