# 1. Сборка фронтенда (React/Vite)
FROM node:20-alpine AS frontend-builder
WORKDIR /web
# Копируем package.json и устанавливаем зависимости
COPY web/gopher-status-web/package*.json ./
RUN npm install
# Копируем исходники фронта и билдим
COPY web/gopher-status-web/ .
RUN npm run build

# 2. Сборка бэкенда (Go)
FROM golang:1.25.5-alpine AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o monitor ./cmd/server/main.go

# 3. Финальный образ
FROM alpine:latest
WORKDIR /app

EXPOSE 8080 5051

# Копируем бинарник
COPY --from=backend-builder /app/monitor .

# !!! ВАЖНО: Копируем собранный фронтенд (папку dist) !!!
# Обратите внимание: путь назначения (/app/web/dist) должен совпадать с тем,
# где ваш Go-сервер ожидает найти статику.
COPY --from=frontend-builder /web/dist ./web/dist

CMD ["./monitor"]
