package storage

import (
	"database/sql"
	"fmt"

	"github.com/mrumyantsev/currency-converter/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter/internal/pkg/consts"
	"github.com/mrumyantsev/currency-converter/internal/pkg/models"
	"github.com/mrumyantsev/currency-converter/internal/pkg/utils"

	_ "github.com/lib/pq"
	"github.com/mrumyantsev/fastlog"
)

type Storage struct {
	config *config.Config
	conn   *sql.DB
}

func New(cfg *config.Config) *Storage {
	storage := &Storage{
		config: cfg,
	}

	return storage
}

// Connects to database.
func (s *Storage) Connect() error {
	var (
		psqlInfo = fmt.Sprintf("host=%s port=%s user=%s "+
			"password=%s dbname=%s sslmode=%s",
			s.config.StorageConnectHostname,
			s.config.StorageConnectPort,
			s.config.StorageConnectUser,
			s.config.StorageConnectPassword,
			s.config.StorageDatabaseName,
			s.config.StorageSSLMode,
		)

		err error
	)

	s.conn, err = sql.Open(s.config.StorageDriver, psqlInfo)
	if err != nil {
		fastlog.Fatal("cannot connect to db", err)
	}

	fastlog.Debug("openning db connection")

	return nil
}

// Disconnects from database.
func (s *Storage) Disconnect() error {
	err := s.conn.Close()
	if err != nil {
		fastlog.Fatal("cannot disconnect from db", err)
	}

	fastlog.Debug("closing db connection")

	return nil
}

func (s *Storage) GetLastDatetime() (*models.UpdateDatetime, error) {
	query := `
		SELECT id, update_datetime
		FROM public.update_datetimes
		WHERE id = (
			SELECT MAX(ID)
			FROM public.update_datetimes
		);
	`

	rows, err := s.conn.Query(query)
	if err != nil {
		return nil, utils.DecorateError("cannot perform select of update datetimes", err)
	}
	defer rows.Close()

	res := models.UpdateDatetime{}

	for rows.Next() {
		err = rows.Scan(&res.Id, &res.UpdateDatetime)
		if err != nil {
			return nil, utils.DecorateError("cannot scan from a row", err)
		}
	}

	return &res, nil
}

func (s *Storage) InsertDatetime(datetime string) (*int, error) {
	query := `
		INSERT INTO public.update_datetimes (update_datetime)
		VALUES
		($1)
		RETURNING id;
	`

	stmt, err := s.conn.Prepare(query)
	if err != nil {
		return nil, utils.DecorateError("cannot prepare statement for inserting datetime", err)
	}

	row := stmt.QueryRow(datetime)
	if err != nil {
		return nil, utils.DecorateError("cannot execute inserting state of datetime", err)
	}

	var id int

	row.Scan(&id)

	return &id, nil
}

func (s *Storage) GetCurrencies(updateDatetimeId int) (*models.CurrencyStorage, error) {
	query := `
		SELECT
			public.info.num_code,
			public.info.char_code,
			public.multipliers.multiplier,
			public.info.name,
			public.currency_values.currency_value
		FROM public.multipliers
		JOIN public.info
		  ON public.multipliers.id = public.info.multiplier_id
		JOIN public.currency_values
		  ON public.info.num_code = public.currency_values.info_num_code
		WHERE public.currency_values.update_datetime_id = $1
		ORDER BY public.info.num_code;
	`

	stmt, err := s.conn.Prepare(query)
	if err != nil {
		return nil, utils.DecorateError("cannot prepare statement for getting currencies", err)
	}

	rows, err := stmt.Query(updateDatetimeId)
	if err != nil {
		return nil, utils.DecorateError("cannot perform select of currencies", err)
	}
	defer rows.Close()

	var (
		currencyStorage models.CurrencyStorage = models.CurrencyStorage{
			Currencies: make(
				[]models.Currency,
				consts.LEN_OF_CURRENCIES_SCLICE_INITIAL,
				consts.CAP_OF_CURRENCIES_SCLICE_INITIAL,
			),
		}
		currency models.Currency
	)

	for rows.Next() {
		err = rows.Scan(
			&currency.NumCode,
			&currency.CharCode,
			&currency.Multiplier,
			&currency.Name,
			&currency.CurrencyValue,
		)
		if err != nil {
			return nil, utils.DecorateError("cannot scan currency entry from a row", err)
		}

		currencyStorage.Currencies = append(currencyStorage.Currencies, currency)
	}

	return &currencyStorage, nil
}

func (s *Storage) InsertCurrencies(currencyStorage *models.CurrencyStorage, updateDatetimeId int) error {
	query := `
		INSERT INTO public.currency_values
		(currency_value, update_datetime_id, info_num_code)
		VALUES
		($1,$2,$3)
	`

	extendCurrenciesQuery(
		&query,
		4, // it is means next placeholder: $4
		0,
		len(currencyStorage.Currencies)-1,
	)

	query += ";"

	stmt, err := s.conn.Prepare(query)
	if err != nil {
		return utils.DecorateError("cannot prepare statement for inserting currencies", err)
	}

	// new struct ?

	_, err = stmt.Exec()
	if err != nil {
		return utils.DecorateError("cannot execute inserting of currencies", err)
	}

	return nil
}

func extendCurrenciesQuery(query *string, startPlaceholder int, startLine int, endLine int) {
	for i := startLine; i < endLine; i++ {
		*query += fmt.Sprintf(",($%d,$%d,$%d)", startPlaceholder, startPlaceholder+1, startPlaceholder+2)
		startPlaceholder += 3
	}
}
