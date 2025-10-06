# Start from the official Golang image
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Use a minimal image for running
FROM alpine:latest
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose port (ubah jika aplikasi Anda menggunakan port lain)
EXPOSE 8080

# Run the binary
CMD ["./main"]
