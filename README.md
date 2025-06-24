# Auth service 🔐
Микросервис аутентификации и авторизации с использованием gRPC и PostgreSQL   
[![Go Version](https://img.shields.io/badge/go-1.24+-blue)](https://golang.org/)

## Установка ⚙️
### Подготовка конфигов
В корне проекта создать файлы `prod.env` `postgres.env` и заполнить   
его по аналогии с `.env.example` `postgres.env.example`

```env
# prod.env
ENV=development
GRPC_ADDRESS=0.0.0.0:40051
POSTGRES_DSN=host=postgres port=5432 dbname=authservice user=postgres password=postgres
```

```env
# postgres.env
POSTGRES_DB=authservice
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
```


### Запуск с использованием docker-compose 🐳

```shell
docker compose up --build
```

