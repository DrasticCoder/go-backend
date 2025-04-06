# Go Backend API Documentation

## 🔧 Core Components

### Router Setup

The application uses Gorilla Mux for routing with global middleware:

```go
router := mux.NewRouter()
router.Use(middleware.ErrorHandler)
router.Use(middleware.RateLimiter)
```

## 🛣️ API Endpoints

### Public Routes

- **Swagger Documentation**

  - `GET /swagger/` - API documentation interface

- **Health Check**

  - `GET /api/v1/health` - Server health status

- **Authentication**
  - `POST /api/v1/auth/register` - User registration
  - `POST /api/v1/auth/login` - User login

### Protected Routes

All protected routes require JWT authentication and specific role-based access:

- `GET /api/v1/users/profile` - User profile (Requires `admin` or `premium` role)
- `GET /api/v1/protected` - Generic protected route

## 🔒 Security Features

### JWT Authentication

```go
// JWT verification middleware
token := r.Header.Get("Authorization")
// Format: Bearer <token>
```

### Role-Based Access Control (RBAC)

- Implemented through middleware
- Supports multiple roles (e.g., "admin", "premium")

## 🛠️ Development Setup

### Database

- PostgreSQL database required
- Database name: `crmdb`

### Running the Server

1. Install dependencies:

```bash
go mod tidy
```

2. Install Swagger CLI:

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

3. Generate Swagger docs:

```bash
swag init -g cmd/main.go
```

4. Run the server:

```bash
go run cmd/main.go
```

Server runs at `http://localhost:8080`

### Hot Reloading

Using `air` for development:

```bash
go install github.com/cosmtrek/air@latest
air init
air
```

## 📚 API Documentation

- Swagger UI available at: `http://localhost:8080/swagger/index.html`
- Authentication format: `Authorization: Bearer <your_token_here>`

## 🔐 Security Notes

- Environment variables are properly secured (`.env` in `.gitignore`)
- JWT tokens required for protected routes
- Role-based access control implemented
- Rate limiting enabled
