package dbstorage

import (
	"database/sql"
	"fmt"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/models"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib/e"

	_ "github.com/lib/pq"
)

const (
	initialCurrenciesCapacity = 50
)

type DbStorage struct {
	conn   *sql.DB // connection to database
	config *config.Config
}

func New(cfg *config.Config) *DbStorage {
	return &DbStorage{config: cfg}
}

// Connect connects to database.
func (s *DbStorage) Connect() error {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=%s",
		s.config.DbHostname,
		s.config.DbPort,
		s.config.DbUsername,
		s.config.DbPassword,
		s.config.DbDatabase,
		s.config.DbSSLMode,
	)

	var err error

	if s.conn, err = sql.Open(s.config.DbDriver, psqlInfo); err != nil {
		return e.Wrap("could not connect to db", err)
	}

	return nil
}

// Disconnect disconnects from database.
func (s *DbStorage) Disconnect() error {
	if err := s.conn.Close(); err != nil {
		return e.Wrap("could not disconnect from db", err)
	}

	return nil
}

func (s *DbStorage) LatestUpdateDatetime() (*models.UpdateDatetime, error) {
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
		return nil, e.Wrap("could not perform select of update datetimes", err)
	}
	defer func() { _ = rows.Close() }()

	updateDatetime := &models.UpdateDatetime{}

	for rows.Next() {
		err = rows.Scan(&updateDatetime.Id, &updateDatetime.UpdateDatetime)
		if err != nil {
			return nil, e.Wrap("could not scan from a row", err)
		}
	}

	return updateDatetime, nil
}

func (s *DbStorage) InsertUpdateDatetime(datetime string) (*models.UpdateDatetime, error) {
	query := `
		INSERT INTO public.update_datetimes (update_datetime)
		VALUES
		($1)
		RETURNING id;
	`

	stmt, err := s.conn.Prepare(query)
	if err != nil {
		return nil, e.Wrap("could not prepare statement for inserting datetime", err)
	}

	row := stmt.QueryRow(datetime)
	if err != nil {
		return nil, e.Wrap("could not execute inserting state of datetime", err)
	}

	updateDatetime := &models.UpdateDatetime{UpdateDatetime: datetime}

	row.Scan(&updateDatetime.Id)

	return updateDatetime, nil
}

func (s *DbStorage) LatestCurrencies(updateDatetimeId int) (*models.Currencies, error) {
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
		ORDER BY public.info.name;
	`

	stmt, err := s.conn.Prepare(query)
	if err != nil {
		return nil, e.Wrap("could not prepare statement for getting currencies", err)
	}

	rows, err := stmt.Query(updateDatetimeId)
	if err != nil {
		return nil, e.Wrap("could not perform select of currencies", err)
	}
	defer func() { _ = rows.Close() }()

	currencies := &models.Currencies{
		Currencies: make(
			[]models.Currency,
			0,
			initialCurrenciesCapacity,
		),
	}

	var currency models.Currency

	for rows.Next() {
		err = rows.Scan(
			&currency.NumCode,
			&currency.CharCode,
			&currency.Multiplier,
			&currency.Name,
			&currency.Value,
		)
		if err != nil {
			return nil, e.Wrap("could not scan currency entry from a row", err)
		}

		currencies.Currencies = append(
			currencies.Currencies,
			currency,
		)
	}

	return currencies, nil
}

func (s *DbStorage) InsertCurrencies(currencies *models.Currencies, updateDatetimeId int) error {
	query := `
		INSERT INTO public.currency_values
		(currency_value, update_datetime_id, info_num_code)
		VALUES
		($1,$2,$3)
	`

	currenciesLength := len(currencies.Currencies)

	extendCurrenciesQuery(
		&query,
		4, // means next placeholder ($4)
		0,
		currenciesLength-1,
	)

	query += ";"

	stmt, err := s.conn.Prepare(query)
	if err != nil {
		return e.Wrap("could not prepare statement for inserting currencies", err)
	}

	entries := []interface{}{}

	for _, currency := range currencies.Currencies {
		entries = append(
			entries,
			currency.Value,
			updateDatetimeId,
			currency.NumCode,
		)
	}

	if _, err = stmt.Exec(entries...); err != nil {
		return e.Wrap("could not execute inserting of currencies", err)
	}

	return nil
}

func extendCurrenciesQuery(query *string, startPlaceholder int, startLine int, endLine int) {
	for i := startLine; i < endLine; i++ {
		*query += fmt.Sprintf(",($%d,$%d,$%d)", startPlaceholder, startPlaceholder+1, startPlaceholder+2)
		startPlaceholder += 3
	}
}
