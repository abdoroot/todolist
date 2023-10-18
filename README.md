# TodoList App with Go and Gin Framework

This is a simple TodoList application built using the Go programming language and the Gin web framework. The application provides a simple user interface for managing your tasks and also exposes a RESTful API for task management. It utilizes PostgreSQL and MySQL databases for storing task data and implements JWT-based authentication for API access.

## Features

- **Task Management**: This application allows you to perform the following tasks:
  - Add new tasks
  - Remove tasks
  - Edit tasks

- **User Authentication**: User authentication is implemented using JSON Web Tokens (JWT) to secure the API endpoints.

- **Database Support**: This application supports both PostgreSQL and MySQL databases for task data storage.

## Prerequisites

Before you can run this application, you'll need to have the following prerequisites installed:

- [Go](https://golang.org/doc/install)
- [Gin Framework](https://github.com/gin-gonic/gin)
- [PostgreSQL](https://www.postgresql.org/download/) or [MySQL](https://dev.mysql.com/downloads/)

## Setup

1. Clone the repository:

   ```shell
   git clone https://github.com/abdoroot/todolist
   ```

2. Change the directory to your project:

   ```shell
   cd todolist
   ```

3. Create a configuration file `config.go`:

4. Install the required dependencies:

   ```shell
   go mod tidy
   ```

5. Run the application:

   ```shell
   go run main.go
   ```

The application will start on the specified port (default is 8080). You can access the UI by opening a web browser and navigating to `http://localhost:8080`. The RESTful API is available at `http://localhost:8080/api`.

## API Endpoints

### Authentication

- `POST /api/login`: Login with a JSON request containing a username and password to receive a JWT token for authentication.

### Task Management

- `GET /api/tasks`: Retrieve a list of tasks.
- `POST /api/tasks`: Create a new task.
- `PUT /api/tasks/:id`: Update an existing task by ID.
- `DELETE /api/tasks/:id`: Delete a task by ID.

## Database Schema

The database schema includes a table to store tasks. You can find the SQL schema in the `schema` directory for both PostgreSQL and MySQL.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.

## Acknowledgments

- This application is built with the Go programming language and the Gin web framework.
- Authentication is implemented using JWT (JSON Web Tokens).
- Database support includes PostgreSQL and MySQL.

Please feel free to contribute to this project, report issues, or suggest improvements.