# Stage 1: Build
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY ./go.mod go.sum ./
COPY ./internal/ internal/
COPY ./pkg/ pkg/
COPY ./services/profile/ services/profile/
RUN go build -o profile-service services/profile/main.go

# Stage 2: Run
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/profile-service .
EXPOSE 5001
CMD ["./profile-service"]
