# backend/Dockerfile
FROM golang:1.22 AS builder
WORKDIR /app

# dentro de backend/Dockerfile
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o backend-app .

FROM gcr.io/distroless/base-debian12
WORKDIR /appx
COPY --from=builder /app/backend-app /app/backend-app
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/backend-app"]
