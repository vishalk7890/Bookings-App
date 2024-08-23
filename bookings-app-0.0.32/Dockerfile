# # Use the official Golang image as a build stage
# FROM golang:1.20 AS builder

# # Set the Current Working Directory inside the container
# WORKDIR /app

# # Copy go.mod and go.sum files
# COPY go.mod go.sum ./

# # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# # Copy the source code into the container
# COPY . .

# # Build the Go app
# RUN go build -o bookings-app ./cmd/web/main.go ./cmd/web/routes.go ./cmd/web/middleware.go && ls -l /app
# RUN chmod +x /app/bookings-app

# # Start a new stage from scratch
# FROM alpine:latest

# # Set the Current Working Directory inside the container
# WORKDIR /root/

# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /app/bookings-app .

# # Expose port 8080 to the outside world
# EXPOSE 8080

# # Command to run the executable
# CMD ["./bookings-app"]
# Use the official Golang image as a build stage
FROM golang:1.20 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go app
RUN go build -o bookings-app ./cmd/web/main.go  ./cmd/web/routes.go ./cmd/web/middleware.go

# Start a new stage from scratch with a more recent Ubuntu base image
FROM ubuntu:22.04

# Set the Current Working Directory inside the container
WORKDIR /root/

# Install necessary libraries
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/bookings-app .

# Ensure the binary is executable
RUN chmod +x ./bookings-app

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./bookings-app"]