# Stage 1: Build
FROM golang:1.24.5-alpine AS builder
WORKDIR /app
COPY ./go.mod go.sum ./
COPY ./internal/ internal/
COPY ./pkg/ pkg/
COPY ./services/auction/ services/auction/
RUN go build -o auction-service services/auction/main.go

# Stage 2: Run
FROM alpine:3.22
WORKDIR /app
COPY --from=builder /app/auction-service .
EXPOSE 7003
CMD ["./auction-service"]
