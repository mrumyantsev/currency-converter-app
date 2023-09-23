package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mrumyantsev/currency-converter/internal/pkg/config"

	"github.com/mrumyantsev/fastlog"
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

func (c *HttpClient) GetCurrencyData() []byte {
	fastlog.Debug("begin http data retrieving...")

	startTime := time.Now()

	url, err := url.Parse(c.config.CurrencySourceUrl)
	if err != nil {
		fastlog.Fatal("cannot parse url", err)
	}

	requ := c.createRequest(url)

	resp, err := c.client.Do(requ)
	if err != nil {
		fastlog.Fatal("cannot send request to server", err)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		fastlog.Fatal("cannot read data from response body", err)
	}
	defer resp.Body.Close()

	elapsedTime := time.Since(startTime)

	fastlog.Debug(fmt.Sprintf("getting http data time overall: %s", elapsedTime))

	return data
}

func (c *HttpClient) createRequest(url *url.URL) *http.Request {
	fastlog.Debug(fmt.Sprintf("using %s request protocol", c.config.HttpRequestProtocol))
	fastlog.Debug(fmt.Sprintf("using user-agent header: %s", c.config.FakeUserAgentHeaderValue))

	return &http.Request{
		Method: "GET",
		URL:    url,
		Proto:  c.config.HttpRequestProtocol,
		Header: map[string][]string{
			USER_AGENT_HEADER_NAME: {c.config.FakeUserAgentHeaderValue},
		},
	}
}
