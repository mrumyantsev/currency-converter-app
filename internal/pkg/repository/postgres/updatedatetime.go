package postgres

import (
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/config"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/database"
	"github.com/mrumyantsev/currency-converter-app/internal/pkg/models"
	"github.com/mrumyantsev/go-errlib"
)

type UpdateDatetimeRepository struct {
	config   *config.Config
	database *database.Database
}

func NewUpdateDatetimeRepository(cfg *config.Config, db *database.Database) *UpdateDatetimeRepository {
	return &UpdateDatetimeRepository{
		config:   cfg,
		database: db,
	}
}

func (r *UpdateDatetimeRepository) Create(datetime string) (models.UpdateDatetime, error) {
	query := `INSERT INTO public.update_datetimes (update_datetime)
VALUES
($1)
RETURNING id;
	`

	updateDatetime := models.UpdateDatetime{
		UpdateDatetime: datetime,
	}

	stmt, err := r.database.Prepare(query)
	if err != nil {
		return updateDatetime, errlib.Wrap(err, "could not prepare statement for inserting datetime")
	}

	row := stmt.QueryRow(datetime)
	if err != nil {
		return updateDatetime, errlib.Wrap(err, "could not execute inserting state of datetime")
	}

	row.Scan(&updateDatetime.Id)

	return updateDatetime, nil
}

func (r *UpdateDatetimeRepository) GetLatest() (models.UpdateDatetime, error) {
	query := `SELECT id, update_datetime
FROM public.update_datetimes
WHERE id = (
	SELECT MAX(ID)
	FROM public.update_datetimes
);
	`

	var updateDatetime models.UpdateDatetime

	rows, err := r.database.Query(query)
	if err != nil {
		return updateDatetime, errlib.Wrap(err, "could not perform select of update datetimes")
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		err = rows.Scan(&updateDatetime.Id, &updateDatetime.UpdateDatetime)
		if err != nil {
			return updateDatetime, errlib.Wrap(err, "could not scan from a row")
		}
	}

	return updateDatetime, nil
}
