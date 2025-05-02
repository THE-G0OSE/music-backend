что нужно для запуска?
postgresql и чистая бд в нём, golang

установка зависимостей
```
go mod tidy
```

создай файл .env по образу и подобию .env-example в корне проекта
```
DB_HOST="your host"
DB_PORT="your port"
DB_USER="your username"
DB_PASSWORD="your password"
DB_NAME="your db name"
```

запуск
```
go run main.go
```
