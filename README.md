# Payslip Generation System

A scalable Golang-based Payslip Generation System using Clean Architecture, PostgreSQL, JWT authentication, RESTful APIs, structured logging (logrus), and audit trails.

## Features
- Employee and Admin authentication (JWT)
- Payslip generation and retrieval
- Attendance, overtime, and reimbursement management
- Modular, testable codebase (Clean Architecture)
- Structured logging and audit trails

---

## Project Structure
```
internal/
  config/           # Configuration and DB connection
  delivery/
    dto/            # Data Transfer Objects
    http/
      handler/      # Gin HTTP Handlers
      middleware/   # Gin Middleware (JWT, logging, etc)
  domain/           # Domain models/entities
  error_const/      # Error constants
  mocks/            # Mock implementations for testing
  repository/
    postgres/       # PostgreSQL repository implementations
  service/          # Business logic (use cases)
  utils/            # Utilities (JWT, logger, etc)
cmd/
  server/           # Main server entrypoint
  migration/        # DB migration tool
  seeder/           # DB seeder tool
database/
  migrations/       # SQL migration scripts
  seeder/           # SQL seed scripts
```

---

## Prerequisites
- Go 1.20+
- Docker & Docker Compose
- PostgreSQL (if not using Docker)

---

## Setup & Run

### 1. Clone the repository
```bash
git clone <repo-url>
cd payroll-system
```

### 2. Environment Variables
Copy `.env.example` to `.env` and adjust as needed (DB credentials, JWT secret, etc).

### 3. Start PostgreSQL (with Docker Compose)
```bash
docker-compose up -d
```

### 4. Run Database Migration
```bash
go run ./cmd/migration/migration.go
```

### 5. Seed Database
```bash
go run ./cmd/seeder/seeder.go
```

### 6. Build the Service
```bash
go build -o payroll-server ./cmd/server
```

### 7. Run the Service
```bash
go run ./cmd/server/main.go
# or
./payroll-server
```

### 8. Run Tests
```bash
go test ./...
```

---

## API Contract

### Health Check
#### GET /healthz
- **Response:**
  ```json
  { "status": "ok" }
  ```

### Authentication
#### POST /api/v1/login/admin
- **Request:**
  ```json
  { "email": "admin@example.com", "password": "string" }
  ```
- **Response:**
  ```json
  {
    "message": "Admin login successful",
    "data": {
      "token": "jwt-token"
    }
  }
  ```

#### POST /api/v1/login/employee
- **Request:**
  ```json
  { "email": "employee@example.com", "password": "string" }
  ```
- **Response:**
  ```json
  {
    "message": "Employee login successful",
    "data": {
      "token": "jwt-token"
    }
  }
  ```

---

## Admin Endpoints (require JWT, admin role)

#### POST /api/v1/admin/payroll-period
- **Body:**
  ```json
  { "start_date": "YYYY-MM-DD", "end_date": "YYYY-MM-DD" }
  ```
- **Response:**
  ```json
  { "message": "Payroll period created", "data": { "period_id": 1 } }
  ```

#### POST /api/v1/admin/payroll-period/run
- **Body:**
  ```json
  { "period_id": 1 }
  ```
- **Response:**
  ```json
  { "message": "Payroll run started", "data": null }
  ```

#### POST /api/v1/admin/payroll-period/lock
- **Body:**
  ```json
  { "period_id": 1 }
  ```
- **Response:**
  ```json
  { "message": "Payroll period locked", "data": null }
  ```

#### GET /api/v1/admin/payroll-summary/:period_id
- **Response:**
  ```json
  { "message": "Payroll summary retrieved successfully", "data": { /* summary object */ } }
  ```

---

## Employee Endpoints (require JWT, employee role)

#### POST /api/v1/employee/attendance
- **Body:**
  ```json
  { "date": "YYYY-MM-DD", "status": "present|absent|leave" }
  ```
- **Response:**
  ```json
  { "message": "Attendance recorded successfully", "data": null }
  ```

#### POST /api/v1/employee/overtime
- **Body:**
  ```json
  { "date": "YYYY-MM-DD", "hours": 2 }
  ```
- **Response:**
  ```json
  { "message": "Overtime submitted successfully", "data": null }
  ```

#### POST /api/v1/employee/reimbursement
- **Body:**
  ```json
  { "date": "YYYY-MM-DD", "amount": 100000, "description": "Medical" }
  ```
- **Response:**
  ```json
  { "message": "Reimbursement submitted successfully", "data": null }
  ```

#### GET /api/v1/employee/payslip/:period_id
- **Response:**
  ```json
  { "message": "Payslip retrieved successfully", "data": { /* payslip object */ } }
  ```

---

### Error Response (all endpoints)
- **Format:**
  ```json
  { "message": "error message", "error": "error details" }
  ```

---

## Low Level Design

### Clean Architecture Layers
- **Delivery Layer:** Gin HTTP handlers, DTOs, middleware
- **Service Layer:** Business logic, use cases
- **Repository Layer:** PostgreSQL data access (using database/sql)
- **Domain Layer:** Entities, interfaces

### Logging & Audit
- All requests and errors are logged using logrus (JSON format)
- Audit trails are stored for sensitive actions (login, payslip access)

### JWT Auth
- JWT tokens are issued on login
- Middleware validates JWT for protected routes

### Example Sequence: Get Payslip
1. **Client** → `POST /api/v1/login` → **Server** (returns JWT)
2. **Client** → `GET /api/v1/employee/payslip` (with JWT) → **Server**
3. **Handler** parses request, validates JWT
4. **Service** fetches payslip from repository
5. **Repository** queries PostgreSQL
6. **Service** returns payslip DTO
7. **Handler** returns JSON response

---

## License
MIT
