# Stage 1: Build the Go application
FROM golang:1.17 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the source code into the container
COPY . .

# Initialize Go module and tidy dependencies
RUN go mod init userservice && go mod tidy

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Use a smaller base image for running the application
FROM alpine:latest  

# Add certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set the working directory inside the container
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/main .

# Expose the port the application runs on
EXPOSE 9000

# Set the command to run the application
CMD ["./main"]
