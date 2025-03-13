# # Start with Golang builder image
# FROM golang:1.23-alpine AS builder

# # Set environment variables
# ENV GO111MODULE=on \
#     CGO_ENABLED=0 \
#     GOOS=linux \
#     GOARCH=amd64

# # Set working directory
# WORKDIR /app

# # Copy go mod and sum files
# COPY go.mod go.sum ./

# # Download dependencies
# RUN go mod download

# # Copy application files
# COPY . .

# # Build the application
# RUN go build -o order-service .

# # Create a minimal runtime image
# FROM alpine:latest

# # Install Postgres client to run migrations
# RUN apk add --no-cache postgresql-client

# # Set environment variables (to be overridden at runtime)
# ENV DATABASE_URL="postgres://postgres:123456@localhost:5432/orders_db?sslmode=disable" \
#     KAFKA_BROKERS="kafka:9092" \
#     SERVER_PORT="8008"

# # Set working directory
# WORKDIR /root/

# # Copy built binary and init.sql from builder
# COPY --from=builder /app/order-service .
# COPY --from=builder /app/db/migrations/init.sql /root/init.sql

# # Expose the application port
# EXPOSE 8008

# # Command to run the executable after running the migrations
# CMD ["sh", "-c", "psql $DATABASE_URL -f /root/init.sql && ./order-service"]

#----------------------------------------------------------------------------------------------------------------------

# Start with Golang builder image
FROM golang:1.23-alpine AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy application files
COPY . .

# Build the application
RUN go build -o order-service .

# Create a minimal runtime image
FROM alpine:latest

# Install Postgres client to run migrations
RUN apk add --no-cache postgresql-client

# Set environment variables (to be overridden at runtime)
ENV DATABASE_URL=${DATABASE_URL} \
    KAFKA_BROKERS=${KAFKA_BROKERS} \
    SERVER_PORT=${SERVER_PORT}

# Set working directory
WORKDIR /root/

# Copy built binary and init.sql from builder
COPY --from=builder /app/order-service .
COPY --from=builder /app/db/init.sql /root/init.sql

# Expose the application port
EXPOSE ${SERVER_PORT}

# Command to run the executable after running the migrations
CMD ["sh", "-c", "psql $DATABASE_URL -f /root/init.sql && ./order-service"]
# CMD ["sh", "-c", " ./order-service"]

#----------------------------------------------------------------------------------------------------------------------
