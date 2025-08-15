# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files first for better caching
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN go build -o backend-app

# Stage 2: Run
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/backend-app .

EXPOSE 8080

CMD ["./backend-app"]