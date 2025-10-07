# Stage 1: Build
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY ./go.mod go.sum ./
COPY ./internal/ internal/
COPY ./pkg/ pkg/
COPY ./services/player/ services/player/
RUN go build -o player-service services/player/main.go

# Stage 2: Run
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/player-service .
EXPOSE 7004
CMD ["./player-service"]
