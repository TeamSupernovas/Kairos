# Start with the official Golang image
FROM golang:1.20

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum, then download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o notification-service .

# Expose the application port
EXPOSE 8080

# Command to run the application
CMD ["./notification-service"]
