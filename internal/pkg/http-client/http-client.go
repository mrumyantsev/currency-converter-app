package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/utils"

	"github.com/mrumyantsev/logx/log"
)

const (
	USER_AGENT_HEADER_NAME = "User-Agent"
)

type HttpClient struct {
	config *config.Config
	client *http.Client
}

func New(cfg *config.Config) *HttpClient {
	httpClient := &HttpClient{
		config: cfg,
		client: &http.Client{},
	}

	return httpClient
}

func (c *HttpClient) GetCurrencyData() ([]byte, error) {
	startTime := time.Now()

	url, err := url.Parse(c.config.CurrencySourceUrl)
	if err != nil {
		return nil, utils.DecorateError("cannot parse url", err)
	}

	requ := c.createRequest(url)

	resp, err := c.client.Do(requ)
	if err != nil {
		return nil, utils.DecorateError("cannot send request to server", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, utils.DecorateError("cannot read data from response body", err)
	}
	defer resp.Body.Close()

	elapsedTime := time.Since(startTime)

	log.Debug(fmt.Sprintf("getting http data time overall: %s", elapsedTime))

	return data, nil
}

func (c *HttpClient) createRequest(url *url.URL) *http.Request {
	log.Debug(fmt.Sprintf("using %s protocol in request", c.config.HttpRequestProtocol))
	log.Debug(fmt.Sprintf("using user-agent header: %s", c.config.FakeUserAgentHeaderValue))

	return &http.Request{
		Method: "GET",
		URL:    url,
		Proto:  c.config.HttpRequestProtocol,
		Header: map[string][]string{
			USER_AGENT_HEADER_NAME: {c.config.FakeUserAgentHeaderValue},
		},
	}
}
