package config

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib"
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

	DbDriver   string `envconfig:"DB_DRIVER" default:"postgres"`
	DbHostname string `envconfig:"DB_HOSTNAME" default:"localhost"`
	DbPort     string `envconfig:"DB_PORT" default:"5432"`
	DbUsername string `envconfig:"DB_USERNAME" default:"postgres"`
	DbPassword string `envconfig:"DB_PASSWORD" default:""`
	DbDatabase string `envconfig:"DB_DATABASE" default:"currency_storage"`
	DbSSLMode  string `envconfig:"DB_SSLMODE" default:"disable"`

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
		return lib.DecorateError("cannot populate struct with environment variables", err)
	}

	return nil
}
