# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

# Enable Go modules and optimize for static builds
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the full source
COPY . .

# Build the Go binary
RUN go build -o rating-service ./main.go

# Stage 2: Minimal runtime image
FROM alpine:latest

# Install PostgreSQL client to run migrations
RUN apk add --no-cache postgresql-client

# Set working directory
WORKDIR /root/

# Copy the binary and migrations
COPY --from=builder /app/rating-service .
COPY --from=builder /app/db/migrations /root/migrations

# # Environment variables passed in at runtime
# ENV DATABASE_URL=${DATABASE_URL} \
#     KAFKA_BROKERS=${KAFKA_BROKERS} \
#     KAFKA_NOTIFICATION_TOPIC=${KAFKA_NOTIFICATION_TOPIC} \
#     SERVER_PORT=${SERVER_PORT:-8080}

# # Expose server port
# EXPOSE ${SERVER_PORT}

# Run migrations and start the app
CMD ["sh", "-c", "psql $DATABASE_URL -f /root/migrations/000001_create_rating_table.up.sql && ./rating-service"]
