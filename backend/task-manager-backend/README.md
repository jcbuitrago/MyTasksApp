# Task Manager Backend

This is a RESTful API for a Task Manager application built using Go. The application implements CRUD operations for users, categories, and tasks, following the MVC design pattern.

## Project Structure

```
task-manager-backend
├── cmd
│   └── server
│       └── main.go            # Entry point of the application
├── config
│   └── config.go             # Configuration settings for the application
├── controllers
│   ├── category_controller.go # Handles HTTP requests related to categories
│   ├── task_controller.go     # Handles HTTP requests related to tasks
│   └── user_controller.go     # Handles HTTP requests related to users
├── models
│   ├── category.go            # Defines the Category model
│   ├── task.go                # Defines the Task model
│   └── user.go                # Defines the User model
├── repositories
│   ├── category_repository.go  # Interacts with the database for categories
│   ├── task_repository.go      # Interacts with the database for tasks
│   └── user_repository.go      # Interacts with the database for users
├── routes
│   └── routes.go              # Sets up application routes
├── services
│   ├── category_service.go     # Business logic for categories
│   ├── task_service.go         # Business logic for tasks
│   └── user_service.go         # Business logic for users
├── utils
│   └── response.go             # Utility functions for formatting responses
├── go.mod                       # Go module definition file
├── go.sum                       # Checksums for module dependencies
└── README.md                    # Documentation for the project
```

## Setup Instructions

1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd task-manager-backend
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Set up environment variables:**
   Create a `.env` file in the root directory and add your database configuration:
   ```
   DB_HOST=your_db_host
   DB_PORT=your_db_port
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=your_db_name
   ```

4. **Run the application:**
   ```
   go run cmd/server/main.go
   ```

## Usage

- The API will be available at `http://localhost:8080`.
- You can access the following endpoints:
  - **Users**
    - `GET /api/users` - Retrieve all users
    - `POST /api/users` - Create a new user
    - `GET /api/users/:id` - Retrieve a user by ID
    - `PUT /api/users/:id` - Update a user by ID
    - `DELETE /api/users/:id` - Delete a user by ID
  - **Categories**
    - `GET /api/categories` - Retrieve all categories
    - `POST /api/categories` - Create a new category
    - `GET /api/categories/:id` - Retrieve a category by ID
    - `PUT /api/categories/:id` - Update a category by ID
    - `DELETE /api/categories/:id` - Delete a category by ID
  - **Tasks**
    - `GET /api/tasks` - Retrieve all tasks
    - `POST /api/tasks` - Create a new task
    - `GET /api/tasks/:id` - Retrieve a task by ID
    - `PUT /api/tasks/:id` - Update a task by ID
    - `DELETE /api/tasks/:id` - Delete a task by ID

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License

This project is licensed under the MIT License.