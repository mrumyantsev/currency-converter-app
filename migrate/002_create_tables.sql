CREATE TABLE IF NOT EXISTS public.multipliers (
	id         INTEGER NOT NULL UNIQUE,
	multiplier INTEGER NOT NULL,
		CONSTRAINT pk_multipliers PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.info (
	num_code      INTEGER    NOT NULL UNIQUE,
	char_code     VARCHAR(3) NOT NULL,
	multiplier_id INTEGER    NOT NULL,
	name          TEXT       NOT NULL,
		CONSTRAINT pk_info PRIMARY KEY (num_code),
		CONSTRAINT fk_info_multipliers FOREIGN KEY (multiplier_id)
			REFERENCES public.multipliers (id) MATCH SIMPLE
			ON UPDATE NO ACTION
			ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS public.update_datetimes (
	id              SERIAL                   NOT NULL UNIQUE,
    update_datetime TIMESTAMP WITH TIME ZONE NOT NULL,
		CONSTRAINT pk_update_datetimes PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS public.currency_values (
	id                 SERIAL        NOT NULL UNIQUE,
	currency_value     NUMERIC(8, 4) NOT NULL,
	update_datetime_id INTEGER       NOT NULL,
	info_num_code      INTEGER       NOT NULL,
		CONSTRAINT pk_currency_values PRIMARY KEY (id),
		CONSTRAINT fk_currency_values_update_datetimes FOREIGN KEY (update_datetime_id)
			REFERENCES public.update_datetimes (id) MATCH SIMPLE
			ON UPDATE NO ACTION
			ON DELETE CASCADE,
		CONSTRAINT fk_currency_values_info FOREIGN KEY (info_num_code)
			REFERENCES public.info (num_code) MATCH SIMPLE
			ON UPDATE NO ACTION
			ON DELETE CASCADE
);
