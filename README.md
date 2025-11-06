# Blog REST API

A RESTful API for a simple blog system built with Go, featuring user authentication, CRUD operations for posts and comments, input validation, and database transactions.

## Features

- ✅ User authentication and authorization (JWT-based)
- ✅ CRUD operations for posts
- ✅ CRUD operations for comments
- ✅ Input validation and error responses
- ✅ Database integration with transactions
- ✅ Security best practices (password hashing, JWT tokens)
- ✅ Proper error handling

## Tech Stack

- **Go 1.21+**
- **Gin** - HTTP web framework
- **PostgreSQL** - Database
- **JWT** - Authentication tokens
- **bcrypt** - Password hashing

## Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Git

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd majoo-case1-rest-api
```

2. Install dependencies:
```bash
go mod download
```

3. Set up environment variables:
```bash
cp .env.example .env
```

Edit `.env` file with your database credentials:
```
DATABASE_URL=postgres://username:password@localhost/blogdb?sslmode=disable
JWT_SECRET=your-secret-key-change-in-production
PORT=8080
```

4. Create PostgreSQL database:
```sql
CREATE DATABASE blogdb;
```

5. Run the application:
```bash
go run main.go
```

The server will start on `http://localhost:8080` (or the port specified in `.env`).

## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication

All protected endpoints require a JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

### Endpoints

#### Authentication

##### Register User
```http
POST /api/v1/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (201 Created):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

##### Login
```http
POST /api/v1/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### Posts

##### Get All Posts
```http
GET /api/v1/posts?page=1&limit=10
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "posts": [
    {
      "id": 1,
      "user_id": 1,
      "title": "My First Post",
      "content": "This is the content...",
      "author": "johndoe",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "page": 1,
  "limit": 10
}
```

##### Get Post by ID
```http
GET /api/v1/posts/:id
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": 1,
  "user_id": 1,
  "title": "My First Post",
  "content": "This is the content...",
  "author": "johndoe",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

##### Create Post
```http
POST /api/v1/posts
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "My New Post",
  "content": "This is the post content..."
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "user_id": 1,
  "title": "My New Post",
  "content": "This is the post content...",
  "author": "johndoe",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

##### Update Post
```http
PUT /api/v1/posts/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated Title",
  "content": "Updated content..."
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "user_id": 1,
  "title": "Updated Title",
  "content": "Updated content...",
  "author": "johndoe",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T01:00:00Z"
}
```

##### Delete Post
```http
DELETE /api/v1/posts/:id
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "message": "Post deleted successfully"
}
```

#### Comments

##### Get Comments by Post
```http
GET /api/v1/posts/:postId/comments
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "post_id": 1,
  "comments": [
    {
      "id": 1,
      "post_id": 1,
      "user_id": 2,
      "content": "Great post!",
      "author": "janedoe",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

##### Get Comment by ID
```http
GET /api/v1/comments/:id
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "id": 1,
  "post_id": 1,
  "user_id": 2,
  "content": "Great post!",
  "author": "janedoe",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

##### Create Comment
```http
POST /api/v1/posts/:postId/comments
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "This is a comment..."
}
```

**Response (201 Created):**
```json
{
  "id": 1,
  "post_id": 1,
  "user_id": 2,
  "content": "This is a comment...",
  "author": "janedoe",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

##### Update Comment
```http
PUT /api/v1/comments/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "content": "Updated comment..."
}
```

**Response (200 OK):**
```json
{
  "id": 1,
  "post_id": 1,
  "user_id": 2,
  "content": "Updated comment...",
  "author": "janedoe",
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T01:00:00Z"
}
```

##### Delete Comment
```http
DELETE /api/v1/comments/:id
Authorization: Bearer <token>
```

**Response (200 OK):**
```json
{
  "message": "Comment deleted successfully"
}
```

### Error Responses

All errors follow a consistent format:

```json
{
  "error": "Bad Request",
  "message": "Invalid input: username is required"
}
```

**Common HTTP Status Codes:**
- `200 OK` - Success
- `201 Created` - Resource created successfully
- `400 Bad Request` - Invalid input or validation error
- `401 Unauthorized` - Missing or invalid authentication token
- `403 Forbidden` - User doesn't have permission (e.g., trying to update/delete another user's post)
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists (e.g., duplicate email/username)
- `500 Internal Server Error` - Server error

## Security Features

1. **Password Hashing**: Passwords are hashed using bcrypt before storage
2. **JWT Authentication**: Secure token-based authentication
3. **Authorization**: Users can only modify their own posts and comments
4. **Input Validation**: All inputs are validated using struct tags
5. **SQL Injection Prevention**: Using parameterized queries
6. **CORS Support**: Configured for cross-origin requests

## Database Schema

### Users Table
- `id` (SERIAL PRIMARY KEY)
- `username` (VARCHAR(50) UNIQUE)
- `email` (VARCHAR(100) UNIQUE)
- `password_hash` (VARCHAR(255))
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Posts Table
- `id` (SERIAL PRIMARY KEY)
- `user_id` (INTEGER, FOREIGN KEY)
- `title` (VARCHAR(255))
- `content` (TEXT)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Comments Table
- `id` (SERIAL PRIMARY KEY)
- `post_id` (INTEGER, FOREIGN KEY)
- `user_id` (INTEGER, FOREIGN KEY)
- `content` (TEXT)
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

## Project Structure

```
majoo-case1-rest-api/
├── main.go                 # Application entry point
├── database/
│   └── database.go         # Database connection and migrations
├── models/
│   ├── user.go            # User model and request/response structs
│   ├── post.go            # Post model and request/response structs
│   └── comment.go         # Comment model and request/response structs
├── handlers/
│   ├── auth.go            # Authentication handlers
│   ├── post.go            # Post CRUD handlers
│   └── comment.go         # Comment CRUD handlers
├── middleware/
│   └── auth.go            # JWT authentication middleware
├── utils/
│   ├── errors.go          # Error response utilities
│   ├── jwt.go             # JWT token generation and validation
│   └── password.go        # Password hashing utilities
├── .env.example           # Environment variables template
├── go.mod                 # Go module dependencies
└── README.md             # This file
```

## Testing the API

You can use tools like Postman, cURL, or any HTTP client to test the API.

### Example: Register and Create a Post

```bash
# 1. Register a user
curl -X POST http://localhost:8080/api/v1/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'

# 2. Login (if needed)
curl -X POST http://localhost:8080/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'

# 3. Create a post (use token from step 1 or 2)
curl -X POST http://localhost:8080/api/v1/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "title": "My First Post",
    "content": "This is my first blog post!"
  }'
```

## Development

### Running in Development Mode

```bash
go run main.go
```

### Building for Production

```bash
go build -o blog-api main.go
./blog-api
```

## License

This project is created for testing purposes.

