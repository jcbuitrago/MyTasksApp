# Etapa builder
FROM golang:1.22 AS builder
WORKDIR /app

# Copiamos solo go.mod y go.sum primero para cachear dependencias
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copiamos el resto del código
COPY backend/. .

# Compilamos
RUN CGO_ENABLED=0 GOOS=linux go build -o backend-app .

# Etapa runtime mínima
FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=builder /app/backend-app /app/backend-app
EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/app/backend-app"]