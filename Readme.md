# Cервис, предоставляющий API по созданию сокращённых ссылок.
___
## Функционал

Сервис принимает следующие запросы по http и grpc:
1. Метод Post,который сохраняет оригинальный URL в базе и возвращает сокращённый.
2. Метод Get, который принимает сокращённый URL и возвращает оригинальный.

## Запуск
```shell
# Клонируем репозиторий
> git clone https://github.com/shiffter/url_shortener
```
```shell
# Запуск с PostgreSQL
> make postgres
```
```shell
# Запуск с in-memory решением
> make inmemory
```

## Запросы
### Общий вид

`GET /get_short_url?short_url=value`
```
POST /create_short_url
body: {"original_url": "https://github.com"}
```