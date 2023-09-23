package storage

const (
	QUERY_INSERT_CURRENCIES = `
		INSERT INTO %s (
			currency_value,
			update_datetime_id,
			info_num_code
		) VALUES (
			?, ?, ?
		);
	`
)
