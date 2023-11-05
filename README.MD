# Movie API

Это простое RESTful API для управления коллекцией фильмов. Он позволяет выполнять операции CRUD (Create, Read, Update, Delete) с записями фильмов.

## Предварительные требования

Перед запуском этого приложения убедитесь, что у вас установлены следующие компоненты:

- Go (Golang): [Инструкции по установке](https://golang.org/doc/install)
- Gorilla Mux: Выполните команду `go get -u github.com/gorilla/mux`
- Zap: Выполните команду `go get -u go.uber.org/zap`

## Установка

1. Склонируйте этот репозиторий на свой компьютер:

git clone https://github.com/yourusername/movie-api.git

markdown
Copy code

2. Перейдите в каталог проекта:

cd movie-api

markdown
Copy code

3. Запустите следующую команду, чтобы запустить приложение:

go run main.go

bash
Copy code

Приложение будет доступно по адресу `http://localhost:8000`.

## Конечные точки (Endpoints)

- `GET /movies`: Получить список всех фильмов.
- `GET /movies/{id}`: Получить информацию о конкретном фильме по ID.
- `POST /movies`: Создать новую запись о фильме.
- `PUT /movies/{id}`: Обновить запись о фильме по ID.
- `DELETE /movies/{id}`: Удалить запись о фильме по ID.

## Использование

Вы можете использовать инструменты, такие как [Postman](https://www.postman.com/) или `curl`, чтобы взаимодействовать с API.

### Пример: Создание фильма

```http
POST http://localhost:8000/movies
Content-Type: application/json

{
    "Isbn": "1234567890",
    "Title": "Пример фильма",
    "Director": {
        "FirstName": "Иван",
        "LastName": "Иванов"
    }
}
Пример: Получение списка фильмов
http
Copy code
GET http://localhost:8000/movies
Вклад
Не стесняйтесь вносить вклад в этот проект, создавая проблемы или предлагая запросы на включение изменений (pull requests).