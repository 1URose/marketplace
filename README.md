# Marketplace Api

## О проекте

Marketplace — REST-сервис для создания и просмотра объявлений с поддержкой аутентификации.
Основные функции:

* Регистрация и вход пользователей (JWT + Redis для хранения сессий).
* CRUD объявлениями: создание (только для аутентифицированных) и публичный просмотр со страницей, сортировкой и фильтрацией.
* Swagger UI для интерактивной документации API.

---

## Архитектура и модули

* **cmd/marketplace** — точка входа, инициализация сервера Gin, CORS, Swagger, graceful shutdown.
* **internal/common**

  * `config` — загрузка настроек из `.env` (Gin, JWT, PostgreSQL, Redis, лимиты объявлений).
  * `db` — клиенты PostgreSQL (pgxpool) и Redis (go-redis).
  * `jwt` — генерация/валидация access и refresh токенов.
  * `logger` — централизованный логгер с привязкой к файлам и строкам.
  * `password` — bcrypt-хеширование и проверка пароля.
* **internal/user\_profile** — управление пользователями (просмотр всех пользователей, для тестирования).
* **internal/auth\_signup** — регистрация, логин, обновление токенов; хранение refresh-токенов в Redis.
* **internal/announcement** — объявления:

  * `domain/ad` — сущность Ad и фильтр AdFilter.
  * `use_cases` — бизнес-логика создания и получения объявлений с учётом пагинации.
  * `infrastructure/repository/postgresql` — реализация AdRepository на PostgreSQL.
  * `transport/rest/ad` — HTTP-хендлеры через Gin + валидация (минимальные/максимальные длины, типы изображений через HEAD-запрос).

---

## Технологии и решения

* **Язык**: Go 1.24 (без CGO)
* **Web-фреймворк**: Gin
* **БД**: PostgreSQL 16 (pgxpool) + Liquibase для миграций
* **Кеш / сессии**: Redis (go-redis v8)
* **Аутентификация**: JWT (golang-jwt/jwt v4), хранение refresh-токенов в Redis
* **Документация API**: Swagger (swaggo/gin-swagger)
* **Логирование**: собственный middleware с указанием файла и строки
* **Контейнеризация**: Docker + Docker Compose
* **CI/CD**: образ билдится на Go-builder-этапе (multistage Dockerfile)

---

## Быстрый старт

1. **Склонировать репозиторий**

   ```bash
   git clone https://github.com/1URose/marketplace.git
   cd marketplace
   ```

2. **Создать файл окружения**
   В корне проекта лежит пример `marketplace/.env`. При необходимости скорректируйте параметры (пароли, порты).

3. **Запустить все сервисы через Docker Compose**

   ```bash
   docker-compose -f docker/docker-compose.yml up -d
   ```

4. **Проверить состояние**

   * Сервисы PostgreSQL, Redis и GO-сервис запустятся в docker-сети `shared`.
   * Swagger UI доступен по адресу:

     ```
     http://localhost:8000/swagger/index.html
     ```

5. **Использование API**

   * **Регистрация**: `POST /auth/signup`
   * **Логин**: `POST /auth/login` → получите `accessToken` и `refreshToken`
   * **Создать объявление**: `POST /ad` с заголовком
     `Authorization: Bearer <accessToken>`
   * **Список объявлений**: `GET /ads?page=1&sort_by=price&sort_order=asc&min_price=100&max_price=1000`

---

Остались вопросы? Загляните в Swagger или в исходники модулей — архитектура построена по принципам “чистой” (слоистой) структуры: domain → use\_cases → repository → transport. Удачной работы!
