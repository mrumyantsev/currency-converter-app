package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/utils"
)

type Config struct {
	IsEnableDebugLogs            bool   `envconfig:"ENABLE_DEBUG_LOGS" default:"false"`
	IsReadCurrencyDataFromFile   bool   `envconfig:"READ_CURRENCIES_FROM_FILE" default:"false"`
	CurrencySourceUrl            string `envconfig:"CURRENCIES_SOURCE_URL" default:"https://www.cbr.ru/scripts/XML_daily.asp"`
	CurrencySourceFile           string `envconfig:"CURRENCIES_SOURCE_FILE" default:"./test/currencies.xml"`
	HttpRequestProtocol          string `envconfig:"HTTP_REQUEST_PROTOCOL" default:"HTTP/2"`
	FakeUserAgentHeaderValue     string `envconfig:"FAKE_USER_AGENT_HEADER_VALUE" default:"Mozilla/5.0 (X11; Linux x86_64)"`
	IsUseMultithreadedParsing    bool   `envconfig:"USE_MULTITHREADED_PARSING" default:"true"`
	TimeWhenNeedToUpdateCurrency string `envconfig:"TIME_WHEN_NEED_TO_UPDATE_CURRENCY" default:"13:30:00"`

	StorageDriver       string `envconfig:"STORAGE_DRIVER" default:"postgres"`
	StorageHostname     string `envconfig:"STORAGE_HOSTNAME" default:"localhost"`
	StoragePort         string `envconfig:"STORAGE_PORT" default:"5432"`
	StorageUser         string `envconfig:"STORAGE_USER" default:"postgres"`
	StoragePassword     string `envconfig:"STORAGE_PASSWORD" default:""`
	StorageDatabaseName string `envconfig:"STORAGE_DATABASE_NAME" default:"currency_storage"`
	StorageSSLMode      string `envconfig:"STORAGE_SSL_MODE" default:"disable"`

	HttpServerListenIp   string `envconfig:"HTTP_SERVER_LISTEN_IP" default:"0.0.0.0"`
	HttpServerListenPort string `envconfig:"HTTP_SERVER_LISTEN_PORT" default:"8080"`
}

// New Config constructor.
func New() *Config {
	return &Config{}
}

// Init initialization from environment variables
func (c *Config) Init() error {
	err := envconfig.Process("", c)
	if err != nil {
		return utils.DecorateError("cannot populate struct with environment variables", err)
	}

	return nil
}
