# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .
RUN go mod tidy 
RUN go get github.com/gin-gonic/gin 
RUN go build -o backend-app

# Stage 2: Run
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/backend-app .

EXPOSE 8080

CMD ["./backend-app"]