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

INSERT INTO public.multipliers (id, multiplier)
VALUES
(0, 1),
(1, 10),
(2, 100),
(3, 1000),
(4, 10000),
(5, 100000),
(6, 1000000),
(7, 10000000);

INSERT INTO public.info (num_code, char_code, multiplier_id, name)
VALUES
(036, 'AUD', 0, 'Австралийский доллар'),
(944, 'AZN', 0, 'Азербайджанский манат'),
(826, 'GBP', 0, 'Фунт стерлингов Соединенного королевства'),
(051, 'AMD', 2, 'Армянский драм'),
(933, 'BYN', 0, 'Белорусский рубль'),
(975, 'BGN', 0, 'Болгарский лев'),
(986, 'BRL', 0, 'Бразильский реал'),
(348, 'HUF', 2, 'Венгерский форинт'),
(704, 'VND', 4, 'Вьетнамский донг'),
(344, 'HKD', 0, 'Гонконгский доллар'),
(981, 'GEL', 0, 'Грузинский лари'),
(208, 'DKK', 0, 'Датская крона'),
(784, 'AED', 0, 'Дирхам ОАЭ'),
(840, 'USD', 0, 'Доллар США'),
(978, 'EUR', 0, 'Евро'),
(818, 'EGP', 1, 'Египетский фунт'),
(356, 'INR', 1, 'Индийская рупия'),
(360, 'IDR', 4, 'Индонезийская рупия'),
(398, 'KZT', 2, 'Казахстанский тенге'),
(124, 'CAD', 0, 'Канадский доллар'),
(634, 'QAR', 0, 'Катарский риал'),
(417, 'KGS', 1, 'Киргизский сом'),
(156, 'CNY', 0, 'Китайский юань'),
(498, 'MDL', 1, 'Молдавский лей'),
(554, 'NZD', 0, 'Новозеландский доллар'),
(578, 'NOK', 1, 'Норвежская крона'),
(985, 'PLN', 0, 'Польский злотый'),
(946, 'RON', 0, 'Румынский лей'),
(960, 'XDR', 0, 'Единица специальных прав заимствования (СДР)'),
(702, 'SGD', 0, 'Сингапурский доллар'),
(972, 'TJS', 1, 'Таджикский сомони'),
(764, 'THB', 1, 'Таиландский бат'),
(949, 'TRY', 1, 'Турецкая лира'),
(934, 'TMT', 0, 'Новый туркменский манат'),
(860, 'UZS', 4, 'Узбекский сум'),
(980, 'UAH', 1, 'Украинская гривна'),
(203, 'CZK', 1, 'Чешская крона'),
(752, 'SEK', 1, 'Шведская крона'),
(756, 'CHF', 0, 'Швейцарский франк'),
(941, 'RSD', 2, 'Сербский динар'),
(710, 'ZAR', 1, 'Южноафриканский рэнд'),
(410, 'KRW', 3, 'Южнокорейская вона'),
(392, 'JPY', 2, 'Японская иена');

INSERT INTO public.update_datetimes (id, update_datetime)
VALUES
(1, '2023-09-20 16:46:18');
ALTER SEQUENCE update_datetimes_id_seq RESTART WITH 2;

INSERT INTO public.currency_values (currency_value, update_datetime_id, info_num_code)
VALUES
('62.3374', 1, 036),
('56.8336', 1, 944),
('119.8247', 1, 826),
('25.0077', 1, 051),
('29.6363', 1, 933),
('52.9206', 1, 975),
('19.8915', 1, 986),
('26.9339', 1, 348),
('40.1251', 1, 704),
('12.3725', 1, 344),
('36.5808', 1, 981),
('13.8840', 1, 208),
('26.3054', 1, 784),
('96.6172', 1, 840),
('103.3699', 1, 978),
('31.2742', 1, 818),
('11.6473', 1, 356),
('62.8159', 1, 360),
('20.4936', 1, 398),
('71.9628', 1, 124),
('26.5432', 1, 634),
('10.8914', 1, 417),
('13.2097', 1, 156),
('53.6229', 1, 498),
('57.4293', 1, 554),
('90.0851', 1, 578),
('22.2518', 1, 985),
('20.7864', 1, 946),
('127.5386', 1, 960),
('70.7611', 1, 702),
('88.1455', 1, 972),
('26.7205', 1, 764),
('35.7602', 1, 949),
('27.6049', 1, 934),
('79.4024', 1, 860),
('26.1603', 1, 980),
('42.3760', 1, 203),
('86.5940', 1, 752),
('107.6755', 1, 756),
('88.0460', 1, 941),
('51.0106', 1, 710),
('72.6390', 1, 410),
('65.3702', 1, 392);
