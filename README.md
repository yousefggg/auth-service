# Mountain Tour Auth Service

[![Go Report Card](https://goreportcard.com)](https://goreportcard.com)
![License](https://shields.io)

Сервис аутентификации и управления пользователями. Является ядром экосистемы **Mountain Tour**, отвечая за безопасность, регистрацию и разграничение прав доступа.

## 🚀 Технологический стек

*   **Language:** Go (Standard `net/http` для чистоты архитектуры)
*   **Database:** PostgreSQL + `pgx`
*   **Security:** `bcrypt` (хеширование), `golang-jwt` (токены)
*   **Logging:** [common-lib](https://github.com)
*   **Testing:** `testify`, `mockery`

## 🏗 Архитектура

Проект реализован согласно принципам **Clean Architecture**. Это обеспечивает независимость бизнес-логики от внешних фреймворков и баз данных.

### Слои приложения:
- **Domain**: Сущности (`User`) и интерфейсы репозиториев.
- **Usecase**: Бизнес-логика (регистрация, вход, валидация).
- **Repository**: Реализация доступа к PostgreSQL.
- **Delivery (Handler)**: HTTP-слой, обработка запросов и JSON.

---

## 🛠 Установка и запуск

### 1. Подготовка окружения
Создайте файл `.env` в корне проекта и укажите настройки:
```env
DB_URL=postgres://user:password@localhost:5432/mountain_tour
JWT_SECRET=your_secret_key
PORT=8080
```

### 2. Загрузка зависимостей
```bash
go mod download
go get github.com/yousefggg/common-lib@main
go mod tidy
```

### 3. Запуск сервиса
```bash
go run cmd/main.go
```

---

## 🛣 API Эндпоинты


| Метод | Путь | Описание | Доступ |
| :--- | :--- | :--- | :--- |
| `POST` | `/auth/register` | Регистрация нового пользователя | Public |
| `POST` | `/auth/login` | Вход и получение JWT-токена | Public |

---

## 🧪 Тестирование

Проект покрыт Unit-тестами с использованием моков.

*   **Usecase:** ~89% покрытия
*   **Delivery:** ~63% покрытия

Запуск всех тестов с отчетом о покрытии:
```bash
go test ./... -v -cover
```

## 📄 Лицензия

Данный проект распространяется под лицензией **MIT**.

## 👤 Контакты

*   **Автор:** Юсеф Муляев
*   **Email:** [mulaev2006@gmail.com](mailto:mulaev2006@gmail.com)
*   **GitHub:** [@yousefggg](https://github.com)