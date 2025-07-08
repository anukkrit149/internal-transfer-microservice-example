# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o app

# Final stage
FROM alpine:3.17

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app .
COPY --from=builder /app/config/env.yaml ./config/

# Expose the application port
EXPOSE 3000

# Command to run the application
ENTRYPOINT ["./app"]
CMD ["api"]
