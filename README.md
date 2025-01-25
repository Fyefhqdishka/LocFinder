# LocFinder

LocFinder - это серверное приложение, предоставляющее услуги определения местоположения по IP-адресу, реализованное на языке Go. Оно включает в себя слой базы данных для хранения данных о местоположении, фронтенд на основе React и Docker-установку. Проект также поддерживает изящное завершение работы сервера и миграцию базы данных.

## Характеристики
-  **Поиск местоположения**: Получение данных о местоположении для IP-адреса клиента или предоставленного IP-адреса.
- **CRUD-операции**: Полная поддержка CRUD для управления записями о местоположении.
- **Миграция баз данных**: Встроенные скрипты миграции с поддержкой составных индексов.
- **React Frontend**: Простой пользовательский интерфейс, расположенный в internal/ui.
- **Докеризованное развертывание**: Удобный запуск приложения и его сервисов с помощью Docker Compose.
- **Докеризованное развертывание**: Удобный запуск приложения и его сервисов с помощью Docker Compose.

## Endpoints
| Endpoint                     | Method   | Description                        |
|------------------------------|----------|------------------------------------|
| `/location`                  | `GET`    | Get location by client IP.         |
| `/location/{ip}`             | `GET`    | Get location for a provided IP.    |
| `/location/{ip}`             | `PUT`    | Update location for a provided IP. |
| `/location/{ip}`             | `DELETE` | Delete location for a provided IP. |
| `/locations`                 | `GET`    | Get all stored locations.          |

### Поддержка CORS
Приложение включает поддержку CORS для `http://localhost:5173`, позволяя использовать такие методы, как `GET`, `POST`, `PUT`, `DELETE` и `OPTIONS`.

## Структура файла
- **`cmd/main.go`**: Точка входа в приложение с обработкой плавного отключения Graceful Shutdown.
- **`internal/app/app.go`**: Настройка и инициализация сервера.
- **`pkg/routes`**: Содержит логику регистрации маршрутов.
- **`internal/handlers`**: Обработчики для конечных точек HTTP.
- **`internal/ui`**: Фронтенд на основе React.
- **`.env`**: Переменные окружения для настройки.

## База данных
### Миграции
Миграция схем баз данных осуществляется с помощью [Goose](https://github.com/pressly/goose). В комплект входят следующие скрипты миграции:

#### Up Migration
```sql
-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ip_address VARCHAR(45) UNIQUE NOT NULL,
    country VARCHAR(100),
    city VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_ip_country_city ON locations(ip_address, country, city);
CREATE INDEX idx_created_at ON locations(created_at DESC);
-- +goose StatementEnd
```

#### Down Migration
```sql
-- +goose Down
-- +goose StatementBegin
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP TABLE IF EXISTS locations;
-- +goose StatementEnd
```

## Makefile
Use the `Makefile` for common tasks:
```make
up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

run-tests:
	go test -v ./internal/handlers ./internal/service
```

## Использование

### Предварительные условия
- [Docker](https://www.docker.com/)
- [Go](https://golang.org/)
- Node.js (для React frontend)

### Шаги для запуска
1. Клонируйте репозиторий:
   ```bash
   git clone https://github.com/Fyefhqdishka/LocFinder.git
   cd LocFinder
   ```

2. Настройте файл `.env`:
   Обновите файл `.env`, добавив в него соответствующие значения. Образцы значений приведены в файле.

3. Запустите приложение:
   ```bash
   make up
   ```

4. Получите доступ к пользовательскому интерфейсу React:
   Перейдите по адресу `http://localhost:5173` в вашем браузере.

5. Используйте конечные точки API:
   Протестируйте конечные точки с помощью таких инструментов, как [Postman](https://www.postman.com/) или `curl`.

6. Запустите тесты:
   ```bash
   make run-tests
   ```

### Остановка приложения
Чтобы остановить работу служб, используйте:
```bash
make down
```

