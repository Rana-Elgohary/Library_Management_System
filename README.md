# Library Management System

## Overview

The Library Management System API allows you to manage authors and books within a library. This API supports CRUD operations (Create, Read, Update, Delete) for both authors and books, as well as search functionality for books by title. Built with Go using the Fiber framework and GORM for ORM, this system connects to a MySQL database.

## Features

- **Author Management:**
  - Create, Read, Update, and Delete authors
  - Soft delete authors

- **Book Management:**
  - Create, Read, Update, and Delete books
  - Soft delete books
  - Search books by title

## Getting Started

### Prerequisites

- Go 1.18+
- MySQL
- Fiber framework
- GORM ORM
- Swagger for API documentation

### Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/yourusername/Library_Management_System.git
    cd Library_Management_System
    ```

2. **Install dependencies:**

    ```bash
    go mod tidy
    ```

3. **Set up your MySQL database:**

    Ensure your MySQL server is running and create a database named `Library_Management_System_PyramakerzTask`. Update the DSN in `config/config.go` with your database credentials.

    ```go
    dsn := "root:password@tcp(127.0.0.1:3306)/Library_Management_System_PyramakerzTask?charset=utf8mb4&parseTime=True&loc=Local"
    ```

4. **Run the application:**

    ```bash
    go run main.go
    ```

    The application will start on `localhost:9090`.

### API Endpoints

#### Authors

- **Get All Authors:**
  - `GET /api/author`
  
- **Get Author by ID:**
  - `GET /api/author/:authorid`
  
- **Create Author:**
  - `POST /api/author`
  - Request body: `{ "name": "Author Name", "email": "author@example.com" }`
  
- **Update Author:**
  - `PUT /api/author/:authorid`
  - Request body: `{ "name": "Updated Name", "email": "updated@example.com" }`
  
- **Delete Author:**
  - `DELETE /api/author/:authorid`
  
- **Soft Delete Author:**
  - `DELETE /api/author/softdelete/:authorid`

#### Books

- **Get All Books:**
  - `GET /api/book`
  
- **Get Book by ID:**
  - `GET /api/book/:bookid`
  
- **Create Book:**
  - `POST /api/book`
  - Request body: `{ "title": "Book Title", "isbn": "1234567890", "publishedDate": "2023-01-01", "authorID": 1 }`
  
- **Update Book:**
  - `PUT /api/book/:bookid`
  - Request body: `{ "title": "Updated Title", "isbn": "0987654321", "publishedDate": "2023-01-01", "authorID": 1 }`
  
- **Delete Book:**
  - `DELETE /api/book/:bookid`
  
- **Soft Delete Book:**
  - `DELETE /api/book/softdelete/:bookid`
  
- **Search Books by Title:**
  - `GET /api/book/search/:title`

### Swagger Documentation

- Access Swagger documentation at `http://localhost:9090/swagger/index.html`.

### Configuration

- **Database Connection:**
  Configure the database connection in `config/config.go`:

    ```go
    dsn := "root:password@tcp(127.0.0.1:3306)/Library_Management_System_PyramakerzTask?charset=utf8mb4&parseTime=True&loc=Local"
    ```

### Running Tests

- Add unit and integration tests to ensure the correctness of your API. Use a testing framework compatible with Go to write and run your tests.

### Contributing

1. **Fork the repository**
2. **Create a new branch for your feature or bugfix**
3. **Commit your changes**
4. **Push to your branch**
5. **Create a pull request**
