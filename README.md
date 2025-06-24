# Auth service 🔐
Микросервис аутентификации и авторизации с использованием gRPC и PostgreSQL
[![Go Version](https://img.shields.io/badge/go-1.24+-blue)](https://golang.org/)

## Установка ⚙️

### Запуск с использованием docker-compose 🐳
В корне проекта создать файл `prod.env` и заполнить   
его по аналогии с `.env.example`

```ENV=development
GRPC_ADDRESS=0.0.0.0:40051
POSTGRES_DSN=host=postgres port=5432 dbname=authservice user=postgres password=postgres
```
