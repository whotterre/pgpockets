# Multistage Dockerfile 
# Stage 1: Build app
FROM golang:1.22.4-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -a -o /app/cmd/app

# Stage 2: Run app
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/cmd/app /app/app

CMD ["/app/app"]

EXPOSE 8080
