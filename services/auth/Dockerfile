# Stage 1: Build
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY ./go.mod go.sum ./
COPY ./internal/ internal/
COPY ./pkg/ pkg/
COPY ./services/auth/ services/auth/
RUN go build -o auth-service services/auth/main.go

# Stage 2: Run
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/auth-service .
EXPOSE 5001
CMD ["./auth-service"]
