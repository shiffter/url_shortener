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
### Создание короткого урла
* Request
```
curl -X POST "http://localhost:8081/get_short_url" -d "short_url=sun.com"
```
* Response
``` json 
{
    "short_url": "tVzcw7MPjG",
    "status": 200,
    "Error": ""
}
```

### Получение оригинального урла

* Request
```shell
curl -X GET "http://localhost:8081/get_short_url?short_url=tVzcw7MPjG"
```
* Response
```json
{
  "original_url": "sun.com",
  "status": 200,
  "Error": ""
}
```