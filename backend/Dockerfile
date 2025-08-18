# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod file first
COPY go.mod ./

# Initialize go.sum and download dependencies
RUN go mod tidy
RUN go mod download

# Copy the rest of the source code
COPY . .

# Ensure all dependencies are available
RUN go mod tidy

# Build the application
RUN go build -o backend-app

# Stage 2: Run
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/backend-app .

EXPOSE 8080

CMD ["./backend-app"]