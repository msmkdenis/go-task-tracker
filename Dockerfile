# docker build --tag go-task-tracker .
# docker run -d -p 7540:7540 -v host-volume:/db go-task-tracker
FROM golang:1.22-alpine AS builder

# Создаем и переходим в директорию приложения.
WORKDIR /app

# Копируем go.mod и go.sum для загрузки зависимостей.
COPY go.mod go.sum ./

# Загружаем все зависимости. Зависимости будут кэшированы, если файлы go.mod и go.sum не были изменены.
RUN go mod download

# Копируем исходный код из текущего каталога в рабочий каталог внутри контейнера.
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o task-tracker ./cmd/tracker/main.go

FROM alpine:3.19

# перемещаем исполняемый и другие файлы в нужную директорию
WORKDIR /app/

COPY --from=builder --chown=app:app app .

ENV SERVER_ADDRESS=":7540"
ENV TOKEN_NAME="token"
ENV SECRET_KEY="secret"
ENV TOKEN_TTL=8
ENV SALT="super-salty-salt"
ENV DBFILE="/db/scheduler.db"

EXPOSE 7540

CMD ["./task-tracker"]