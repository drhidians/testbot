# Telegram Bot API
 
Сервис который представляет из себя Telegram-бота и HTTP API.
Обновление бот  получает через вебхук.
Бот и API запущены на одном порту и доступны по разным URL.

Конфигурацию сервиса определяется двумя способами:
1. переменные окружение (все переменные должны иметь префикс с названия сервиса, например `TESTBOT_*`)
2. параметры командной строки

Общие параметры:

- `db` - строка для подключения к базе данных (`postgres://user:pass@host:port/dbname`)
- `db-max-open-conns` - максимальный размер пула подключений к БД;
- `db-max-idle-conns` - максимальное кол-во простаювающих соеденений к БД;
- `secret` - секрет для подписи JWT-токенов;
- `addr` - адресс на котром будет запущен сервер (`:8000`, `localhost:8000`, `...`)
- `domain` - домен на который будет установлен вебхук и который нужно использовать при построений абсолютных URL.

При запуска миграции базы данных  автоматически применяться (если нужно).

Скиллы:

- Golang 1.13.3
- [PostgreSQL 12](https://hub.docker.com/_/postgres)
- [Telegram Bot API](https://core.telegram.org/bots/api)

 - организация проекта следует принципам [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
 - добавлен логирование;
 - unit и/или интеграционные тесты;
 - следую принципам [12 Factor Apps](https://12factor.net)


 - не используются глобальные переменные и синглтоны;
 - функция которая осуществляет сетевое взаимодействие  принимает контекст (`context.Context`);
 - не использую HTTP-фреймворки. Использую роутер chi.

## Бот

Бот  по команде `/start` добавляет пользователя в базу данных (id, имя, фамилия, username, первая фотография профиля, дата первого взаимодействия с ботом, дата последнего обновление пользователя), обновляет эту информацию при любом взаимодействии с ботом (если необходимо) и отвечает [JWT-токеном](http://jwt.io) который в дальнейшем может быть использован для взаимодействия с API.

При запуске бот,  проверяет установлен ли вебхук на нужный URL, и если нет то устанавливать его.

Путь для вебхука: `/bot/webhook`

Для того чтобы убедится в том что запрос пришел именно от Telegram, нужно проверить вхождение IP-адреса клиента в подсети `149.154.160.0/20` и `91.108.4.0/22`.

Параметры:

- `bot-token` - токен бота полученный у [@BotFather](http://telega.one/BotFather);
- - домен на который будет установлен вебхук (подразумевается что перед сервис всегда будет стоять реверс-прокси на котором настроен HTTPS, для локальной разработки следует использовать [ngrok](https://ngrok.com), например);
- `bot-webhook-max-conns` - параметр который передается Telegram Bot API при [установке вебхук](https://core.telegram.org/bots/api#setwebhook)а, представляет из себя максимальное количество параллельных HTTP-соединении от Telegram-сервера;

## API

### Получение информации о боте

Метод публичный и доступный без авторизации. Нужен для того чтобы не хардкодить @username бота на фронтенде. Фактически это поля из результата вызова метода `[getMe](https://core.telegram.org/bots/api#getme)`

```http
GET /bot
```

```json5
{
    // Telegram ID бота
    "id": 12345,

    // @username бота
    "username": "TestTaskBot",

    // Имя бота
    "name": "Test Task Bot"
}
```

### Получение информации о текущем пользователе

Данный метод требует авторизации в виде JWT-токена полученого в ответе на команду `/start`.

Токен может быть передан в виде HTTP-заголовка (`Authorization: Bearer <token>`) или query-параметра (`?token=<token>`).

```http
GET /user
```

```json5
{
    // ID пользователя в сервисе
    "id": 1234,

    // ID пользователя в Telegram
    "externalId": 1939494,

    // Username пользователя в Telegram.
    "username": "durov",

    // Имя пользователя (first name + last name)
    "name": "Pavel Durov",

    // ссылка на аватар пользователя (детальней - ниже)
    "avatar": "https://<domain>/media/<file_id>",

    // Язык пользователя (IETF)
    "language": "ru-RU",

    // Дата первого взаимодействия пользователя с ботом (unix time)
    "joinedAt": 1572887840,

    // Дата последнего обновления данных пользователя
    "updatedAt": null,
}
```

### Проксирование медиа

Чтобы не хранить автарки пользователя, их нужно запрашивать по требованию, "стримить" клиенту ответ Telegram Bot API (сам файл).
реализация HTTP-кэширования на основе `file_id`, который как известно **уникальный** и **привязан к конкретному файлу**.

```http
GET https://{domain}/media/{file_id}
```

# TO DO

- unit test's bot service
- unit test's user service
