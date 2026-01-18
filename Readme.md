# ECM пример e-commerce сайта

## Запуск и отладка
Для запуска приложения воспользуйтесь стандартным билдером docker compose
```
docker compose up --force-recreate --build -d
```

### Domain mapping
```
http://localhost:3000 - frontend
http://ecm_back:8080 - backend
http://ecm-postgres-1:5432 - psql
```

## Стек
- Go
- Node.js
- Postgres
- Docker

