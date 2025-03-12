# Start from the official Golang image
FROM golang:1.21 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd/app/main.go

# Use a minimal base image for production
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Install required dependencies (if needed)
RUN apk --no-cache add ca-certificates

# Copy compiled binary from builder
COPY --from=builder /app/main /root/main

# Copy the environment file (optional, for reference)
COPY gds.env . 

# Expose application port
EXPOSE 9091

# Run the Go application and prevent container exit on failure
#CMD ["sh", "-c", "/root/main || tail -f /dev/null"]

CMD ["./main"]