package config

import (
	"errors"

	"github.com/kelseyhightower/envconfig"
	"github.com/mrumyantsev/go-errlib"
)

const (
	EnvPrefix = ""
)

// A Config is the application configuration structure.
type Config struct {
	IsEnableDebugLogs            bool   `envconfig:"ENABLE_DEBUG_LOGS" default:"false"`
	IsReadCurrencyDataFromFile   bool   `envconfig:"READ_CURRENCIES_FROM_FILE" default:"false"`
	CurrencySourceUrl            string `envconfig:"CURRENCIES_SOURCE_URL" default:"https://www.cbr.ru/scripts/XML_daily.asp"`
	CurrencySourceFile           string `envconfig:"CURRENCIES_SOURCE_FILE" default:"currencies.xml"`
	HttpRequestProtocol          string `envconfig:"HTTP_REQUEST_PROTOCOL" default:"HTTP/2"`
	FakeUserAgentHeaderValue     string `envconfig:"FAKE_USER_AGENT_HEADER_VALUE" default:"Mozilla/5.0 (X11; Linux x86_64)"`
	IsUseMultithreadedParsing    bool   `envconfig:"USE_MULTITHREADED_PARSING" default:"true"`
	TimeWhenNeedToUpdateCurrency string `envconfig:"TIME_WHEN_NEED_TO_UPDATE_CURRENCY" default:"13:30:00"`
	InitialCurrenciesCapacity    int    `envconfig:"INITIAL_CURRENCIES_CAPACITY" default:"50"`

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

// New creates an application configuration.
func New() *Config {
	return &Config{}
}

// Init initializes application configuration.
func (c *Config) Init() error {
	if err := envconfig.Process(EnvPrefix, c); err != nil {
		return errlib.Wrap(err, "could not populate config structure")
	}

	if c.DbPassword == "" {
		return errors.New("no database password specified")
	}

	return nil
}
