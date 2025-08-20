# Proyecto de Nivelación - Gestor de Tareas (MISO)

## Integrantes

|Nombre|Correo|Codigo|
|------|------|------|
|Francisco Santamaría|F.santamaria@uniandes.edu.co|202022134|
|David Octavio Ibarra Muñoz| d.ibarra@uniandes.edu.co| 202014446|
|Luis Fernando Ruiz| lf.ruizo1@uniandes.edu.co| 202211513|
|Juan Camilo Buitrago|Correo | Codigo

Este repositorio contiene el desarrollo del **Proyecto No. 0 - Ejercicio de Nivelación** para el curso **"Desarrollo de Software en la Nube"** de la Maestría en Ingeniería de Software (MISO) de la Universidad de los Andes.

El objetivo es construir una aplicación web completa para la gestión de tareas, implementando una API REST en Go, una base de datos PostgreSQL y una interfaz de usuario web, todo orquestado a través de contenedores Docker.

## 🚀 Stack Tecnológico

  * **Backend**: Go (Golang)
  * **Base de Datos**: PostgreSQL 15
  * **Frontend**: HTML5, CSS3, JavaScript (Vanilla)
  * **Servidor Web**: Nginx
  * **Orquestación**: Docker & Docker Compose

## ✨ Funcionalidades Principales

La aplicación cumple con los siguientes requisitos funcionales:

### 1\. Gestión de Usuarios y Autenticación

  * Creación de cuentas de usuario con credenciales.
  * Inicio y cierre de sesión mediante tokens de autenticación (JWT).
  * Asignación de un avatar por defecto.

### 2\. Gestión de Categorías

  * Crear, visualizar, actualizar y eliminar categorías para organizar las tareas (ej. "Trabajo", "Hogar").

### 3\. Gestión de Tareas

  * Crear tareas asociadas a una categoría.
  * Visualizar la lista de tareas con filtros por categoría y/o estado.
  * Actualizar el estado de una tarea ("Sin Empezar", "Empezada", "Finalizada").
  * Modificar la descripción y la fecha tentativa de finalización.
  * Eliminar tareas.
  * Registro automático de la fecha de creación de cada tarea.

## 📋 Requisitos Previos

  * [**Docker**](https://docs.docker.com/engine/install/)
  * [**Docker Compose**](https://docs.docker.com/compose/install/)

## ⚙️ Guía de Inicio Rápido

### 1\. Clonar el Repositorio

```bash
git clone <URL-del-repositorio>
cd MyTasksApp
```

### 2\. Configurar Variables de Entorno

Crea un archivo `.env` en la raíz del proyecto con el siguiente contenido:

```env
# Credenciales para la Base de Datos
POSTGRES_USER=myuser
POSTGRES_PASSWORD=mypassword
POSTGRES_DB=mytasks_db

# Secreto para la autenticación del Backend (JWT)
JWT_SECRET=un_secreto_muy_largo_y_seguro_aqui
```

### 3\. Construir y Ejecutar la Aplicación

Ejecuta el siguiente comando en la raíz del proyecto para construir las imágenes y levantar los servicios:

```bash
docker-compose up --build
```

Para ejecutar en segundo plano, utiliza la bandera `-d`.

## 🌐 Acceso a la Aplicación

  * **🖥️ Frontend (Aplicación Web)**:
      * [**http://localhost:3000**](http://localhost:3000/)
  * **⚙️ Backend (API)**:
      * [**http://localhost:8080**](http://localhost:8080)
  * **🗃️ Base de Datos (PostgreSQL)**:
      * **Host**: `localhost`
      * **Puerto**: `5433`
      * **Credenciales**: Las definidas en el archivo `.env`.

## 🛑 Detener la Aplicación

```bash
docker-compose down
```

## 📄 Endpoints de la API REST

La documentación completa de la API se encuentra en la colección de Postman del proyecto. A continuación, se resumen los endpoints principales:

#### Usuarios

  * `POST /usuarios`: Crear un nuevo usuario.
  * `POST /usuarios/iniciar-sesion`: Iniciar sesión y obtener un token.

#### Categorías

  * `GET /categorias`: Obtener todas las categorías.
  * `POST /categorias`: Crear una nueva categoría.
  * `DELETE /categorias/{id}`: Eliminar una categoría.

#### Tareas

  * `GET /tareas/usuario`: Obtener todas las tareas del usuario autenticado.
  * `GET /tareas/{id}`: Obtener una tarea específica por su ID.
  * `POST /tareas`: Crear una nueva tarea.
  * `PUT /tareas/{id}`: Actualizar una tarea existente.
  * `DELETE /tareas/{id}`: Eliminar una tarea.
