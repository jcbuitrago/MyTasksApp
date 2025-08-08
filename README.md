# MyTasksApp Deployment Guide

## Prerequisites

- [Docker](https://docs.docker.com/get-docker/) installed
- [Docker Compose](https://docs.docker.com/compose/install/) installed

## Project Structure

<!-- ...existing code or add your project structure here... -->
|-- backend
|   |-- Dockerfile
|   |-- main.go
|   |-- go.mod
|-- frontend
|   |-- Dockerfile
|   |-- index.html
|-- db
|   |-- init.sql
|--docker-compose.yml
|--.env

## Docker Compose Deployment

1. **Create/Edit `docker-compose.yml`**  
   Ensure you have a `docker-compose.yml` file in the project root. Example:

   ```yaml
   services:
  db:
    image: postgres:15
    restart: always
    environment:
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
      POSTGRES_DB: ${POSTGRES_DB}
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./db/init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - app-network

  backend:
    build: 
      context: ./backend
    container_name: backend
    depends_on:
      - db
    ports:
      - "8080:8080"
    environment:
      DB_HOST: db
      DB_PORT: 5432
      DB_USER: ${POSTGRES_USER}
      DB_PASSWORD: ${POSTGRES_PASSWORD}
      DB_NAME: ${POSTGRES_DB}
    networks:
      - app-network

  frontend:
    build: 
      context: ./frontend
    container_name: frontend
    ports:
      - "3000:80"
    networks:
      - app-network
    depends_on:
      - backend

  networks:
    app-network:

  volumes:
    pgdata:
   ```

   Adjust the service names, ports, and environment variables as needed for your project.

   Ensure you have a `.env` file in the project root. Example:

  ```.env
  POSTGRES_USER=<user>
  POSTGRES_PASSWORD=<password>
  POSTGRES_DB=<db_name>
  ```

2. **Build and Start Services**

   ```sh
   docker-compose up --build
   ```

   This will build and start all services defined in `docker-compose.yml`.

3. **Stopping Services**

   ```sh
   docker-compose down (--rmi all #if you want to delete all images)
   ```

## Testing Each Service

- **App Service**  
  Access the app backend at [http://localhost:8080/api/hello] (or the port you mapped).

  To run backend individually:
  
  cd backend
  docker build -t backend-app .
  docker run -d -p 8080:8080 --name backend backend-app

  To stop:

  docker stop backend && docker rm backend



- **Frontend Service**  
  To connect to the Frontend for testing access [http://localhost:3000]:
  
  cd frontend
  docker build -t frontend-app .
  docker run -d -p 3000:80 --name frontend frontend-app

  To stop:

  docker stop frontend && docker rm frontend



## Environment Variables

- Set environment variables in the `docker-compose.yml` file under each service.
- For sensitive data, consider using a `.env` file and reference it in `docker-compose.yml` with `env_file: .env`.

## Logs

To view logs for all services:
```sh
docker-compose logs
```
Or for a specific service:
```sh
docker-compose logs app
```
