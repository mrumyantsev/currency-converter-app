package postgres

import (
	"fmt"

	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/database"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/models"
	"github.com/mrumyantsev/currency-converter-app/pkg/lib/errlib"
)

type CurrenciesRepository struct {
	config   *config.Config
	database *database.Database
}

func NewCurrenciesRepository(cfg *config.Config, db *database.Database) *CurrenciesRepository {
	return &CurrenciesRepository{
		config:   cfg,
		database: db,
	}
}

func (r *CurrenciesRepository) Create(currencies models.Currencies, updateDatetimeId int) error {
	query := `INSERT INTO public.currency_values
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

	stmt, err := r.database.Prepare(query)
	if err != nil {
		return errlib.Wrap("could not prepare statement for inserting currencies", err)
	}

	entries := []any{}

	for _, currency := range currencies.Currencies {
		entries = append(
			entries,
			currency.Value,
			updateDatetimeId,
			currency.NumCode,
		)
	}

	if _, err = stmt.Exec(entries...); err != nil {
		return errlib.Wrap("could not execute inserting of currencies", err)
	}

	return nil
}

func (r *CurrenciesRepository) GetLatest(updateDatetimeId int) (models.Currencies, error) {
	query := `SELECT
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

	currencies := models.Currencies{
		Currencies: make(
			[]models.Currency,
			0,
			r.config.InitialCurrenciesCapacity,
		),
	}

	stmt, err := r.database.Prepare(query)
	if err != nil {
		return currencies, errlib.Wrap("could not prepare statement for getting currencies", err)
	}

	rows, err := stmt.Query(updateDatetimeId)
	if err != nil {
		return currencies, errlib.Wrap("could not perform select of currencies", err)
	}
	defer func() { _ = rows.Close() }()

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
			return currencies, errlib.Wrap("could not scan currency entry from a row", err)
		}

		currencies.Currencies = append(
			currencies.Currencies,
			currency,
		)
	}

	return currencies, nil
}

func extendCurrenciesQuery(query *string, startPlaceholder int, startLine int, endLine int) {
	for i := startLine; i < endLine; i++ {
		*query += fmt.Sprintf(",($%d,$%d,$%d)", startPlaceholder, startPlaceholder+1, startPlaceholder+2)
		startPlaceholder += 3
	}
}
