Currency Converter - это приложение, позволяющее получать данные об иностранных валютах, сохранять их базу (для возможного построения статистики), и конвертировать из одной в другую, используя удобный веб-интерфейс. Данные о валютах берутся из государственных источников страны РФ. Они форматируются из одного формата в другой, более подходящий для работы клиентской части. За обновлением данных следит внутренний планировщик, который запускает обновление из источников каждый день в указанное время (по умолчанию 13:30). Данные из открытых источников обновляются раз в сутки, так что можно настроить желаемое время их получения.

Существует возможность сохранить данные, получаемы из сети в файл на диске (или указать файл как источник данных для обновления). Сохранение в файл производится командой `parserd --save`. Для запуска основного функционала приложения параметры не указываются.

Клиентский код приложения не производит сортировку данных (они приходят к нему уже отсортированными). Он также следит за обновлениями и проверяет, доступен ли сервер для получения данных. По умолчанию запрос к серверу повторяется каждые 5 минут.

![screenshot](https://github.com/mrumyantsev/currency-converter/assets/36193247/9c675152-4ffd-4f4c-b733-25b04c9a4996)

*Вид приложения в браузере*

![console](https://github.com/mrumyantsev/currency-converter/assets/36193247/272c0aeb-785e-4287-bbc1-832901277411)

*Логи в консоли серверного приложения*

![database_schema](https://github.com/mrumyantsev/currency-converter/assets/36193247/a28ec201-abec-4516-bf6b-81bb8a6241ac)

*Схема таблиц в базе данных*
