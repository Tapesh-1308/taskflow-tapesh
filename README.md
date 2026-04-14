# TaskFlow

A modern full-stack task management application built with **Go**, **React**, and **PostgreSQL**. TaskFlow enables teams to organize projects, create tasks, assign work, and track progress with a clean, intuitive interface.

## Overview

TaskFlow is a project management tool that allows users to:
- **Create and manage projects** with descriptions and ownership
- **Create and organize tasks** within projects with status tracking (todo, in_progress, done)
- **Assign tasks** to team members
- **Filter and search** tasks by status, assignee, and project
- **User authentication** with JWT-based security

### Tech Stack

**Backend:**
- **Language:** Go 1.25.1
- **Framework:** Gin (HTTP web framework)
- **Database:** PostgreSQL 15
- **Authentication:** JWT (golang-jwt)
- **Password Hashing:** bcrypt
- **Database Driver:** pgx (PostgreSQL driver)
- **Logging:** log/slog (structured logging)
- **Environment:** godotenv

**Frontend:**
- **Framework:** React 18 with TypeScript
- **Build Tool:** Vite
- **UI Components:** Radix UI
- **State Management:** React Context + React Query
- **Styling:** Tailwind CSS
- **HTTP Client:** Axios
- **Form Handling:** React Hook Form

**Infrastructure:**
- **Containerization:** Docker & Docker Compose
- **Database Migrations:** SQL-based migrations
- **CORS:** Gin CORS middleware

---

## Architecture Decisions

### Backend Architecture: Clean Layered Architecture

The backend follows a **three-tier layered architecture** with clear separation of concerns:

```
cmd/main.go                    (Entry point & dependency injection)
├── handlers                   (HTTP request/response handling)
├── services                   (Business logic & validation)
├── repositories               (Data access layer)
├── models                     (Domain entities)
├── middleware                 (Auth, logging, CORS)
├── config                     (Configuration management)
├── db                         (Database connection)
└── logger                     (Structured logging - singleton pattern)
```

**Why this structure?**
- **Separation of Concerns:** Each layer has a single responsibility
- **Testability:** Repositories and services can be mocked independently
- **Dependency Injection:** Injected at application startup (main.go) for loose coupling
- **Scalability:** Easy to add new features without affecting existing layers

### Key Design Patterns

1. **Repository Pattern:**
   - Database operations are abstracted through interfaces
   - Makes switching databases or adding caching straightforward
   - Each entity (user, project, task) has its own repository

2. **Service Layer Pattern:**
   - Business logic lives in services, not handlers or repositories
   - Validation happens at the service level
   - Services compose multiple repositories for complex operations

3. **Dependency Injection:**
   - All dependencies (DB, logger, config) injected at construction
   - Main.go orchestrates dependency creation
   - Enables easy testing and configuration swapping

4. **Singleton Logger:**
   - Logger initialized once at application startup
   - Injected into all services, handlers, and repositories
   - Ensures consistent structured logging across the application
   - Uses Go's `log/slog` for JSON structured logging

### Frontend Architecture: Component-Based with Context

```
src/
├── components/ui/              (Reusable UI components)
├── features/
│   ├── auth/                   (Authentication feature)
│   ├── project/                (Project management feature)
│   └── task/                   (Task management feature)
├── layouts/                    (Layout components)
├── services/                   (API client & utilities)
└── app/
    ├── providers.tsx           (Context & provider setup)
    └── router.tsx              (React Router configuration)
```

**Why this structure?**
- **Feature-Based Organization:** Each feature is self-contained with its own API calls, types, and components
- **UI Component Library:** Radix UI provides accessible, unstyled components
- **Context API:** Lightweight state management for auth without Redux overhead
- **Separation:** API calls isolated in `api/` folders within features

### Architectural Tradeoffs & Decisions

| Decision | Rationale | Tradeoff |
|----------|-----------|----------|
| Layered monolith (not microservices) | Simpler deployment, easier debugging for a CRUD app | Not ideal if services need independent scaling |
| Context API instead of Redux | Lightweight, reduces boilerplate for small teams | Redux better for complex global state |
| Singleton logger injected everywhere | Ensures consistent logging, easier debugging | Tight coupling to logger (though interface-based) |
| JWT stateless auth | Scales without session storage, simpler infrastructure | No built-in revocation (would need token blacklist) |
| SQL migrations manually managed | Full control, no ORM limitations | More SQL to write, schema changes require migrations |
| Gin instead of net/http | Faster, built-in middleware, cleaner routing | Adds dependency, heavier than std library |
| Vite instead of Create React App | Faster builds, smaller bundle, better dev experience | Less ecosystem maturity than CRA |

### Intentionally Left Out (& Why)

1. **Real-time Updates (WebSockets):** 
   - Complexity not justified for MVP
   - Polling sufficient for most workflows
   - Can add later with Socket.io

2. **Advanced Authorization (RBAC):**
   - Only basic ownership checks implemented
   - Users can only manage their own projects
   - Role-based permissions could be added later
---

## Running Locally

### Prerequisites
- Docker
- Docker Compose

### Quick Start

```bash
# Clone the repository
git clone https://github.com/Tapesh-1308/taskflow-tapesh.git
cd taskflow

# Copy environment file
cp .env.example .env
cp backend/.env.example backend/.env
cp frontend/.env.example frontend/.env

# Start the application with Docker Compose
docker compose up --build

# Application will be available at:
# Frontend: http://localhost:5173
# Backend:  http://localhost:8080
```

### First Time Setup

After `docker compose up` completes:

1. **Database will automatically initialize** - PostgreSQL starts and accepts connections
2. **Migrations will run** - Schema is created from SQL migration files
3. **Frontend and backend services start** - Both available on the ports above

### Stopping the Application

```bash
# Stop all services
docker compose down
```
---

### Automatic Migration (Current Setup)

Migrations **do not run automatically** on startup. They must be run manually:


## Test Credentials

Use these credentials to log in immediately:

```
Email:    test@example.com
Password: password123
```

## API Reference

### Base URL
```
http://localhost:8080
```

### Authentication
Include JWT token in Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

### Public Endpoints

#### Register User
```http
POST /auth/register
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securePassword123"
}
```

**Response:** `201 Created`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123"
}
```

**Response:** `200 OK`
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### Authenticated Endpoints (Require JWT)

#### Get Current User Info
```http
GET /me
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Test User"
}
```

#### Get All Users
```http
GET /users?search=john
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe"
  }
]
```

---

### Project Endpoints

#### Create Project
```http
POST /projects
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Website Redesign",
  "description": "Redesign company website for 2024"
}
```

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Website Redesign",
  "description": "Redesign company website for 2024",
  "owner_id": "550e8400-e29b-41d4-a716-446655440001",
  "created_at": "2024-04-14T10:30:00Z"
}
```

#### List User's Projects
```http
GET /projects
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "projects": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Website Redesign",
      "description": "Redesign company website for 2024",
      "owner": {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "name": "Test User"
      },
      "created_at": "2024-04-14T10:30:00Z"
    }
  ]
}
```

#### Get Project with Tasks
```http
GET /projects/:id
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Website Redesign",
  "description": "Redesign company website for 2024",
  "owner": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Test User"
  },
  "tasks": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440002",
      "title": "Create wireframes",
      "status": "in_progress",
      "project_id": "550e8400-e29b-41d4-a716-446655440000",
      "assignee": {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "name": "Test User"
      },
      "created_at": "2024-04-14T10:35:00Z"
    }
  ],
  "created_at": "2024-04-14T10:30:00Z"
}
```

#### Update Project
```http
PATCH /projects/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "name": "Website Redesign v2",
  "description": "Updated description"
}
```

**Response:** `200 OK`
```json
{
  "message": "updated"
}
```

#### Delete Project
```http
DELETE /projects/:id
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "message": "deleted"
}
```

---

### Task Endpoints

#### Create Task
```http
POST /projects/:id/tasks
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Design homepage mockup"
}
```

**Response:** `201 Created`
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440002",
  "title": "Design homepage mockup",
  "status": "todo",
  "project_id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2024-04-14T10:40:00Z"
}
```

#### List Project Tasks
```http
GET /projects/:id/tasks?status=in_progress&assignee=550e8400-e29b-41d4-a716-446655440001
Authorization: Bearer <token>
```

**Query Parameters:**
- `status` (optional): `todo`, `in_progress`, `done`
- `assignee` (optional): User ID to filter by

**Response:** `200 OK`
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "title": "Design homepage mockup",
    "status": "in_progress",
    "project_id": "550e8400-e29b-41d4-a716-446655440000",
    "assignee": {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "Test User"
    },
    "created_at": "2024-04-14T10:40:00Z"
  }
]
```

#### Update Task
```http
PATCH /tasks/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Design homepage mockup - updated",
  "status": "done",
  "assignee_id": "550e8400-e29b-41d4-a716-446655440003"
}
```

**Response:** `200 OK`
```json
{
  "message": "updated"
}
```

#### Delete Task
```http
DELETE /tasks/:id
Authorization: Bearer <token>
```

**Response:** `200 OK`
```json
{
  "message": "deleted"
}
```

---

### Error Responses

All error responses follow this format:

```json
{
  "error": "error description"
}
```

**Common Status Codes:**
- `400` - Validation failed (invalid request body)
- `401` - Unauthorized (missing or invalid token)
- `403` - Forbidden (user not authorized to perform action)
- `404` - Not found (resource doesn't exist)
- `500` - Internal server error

**Example Error:**
```http
HTTP/1.1 403 Forbidden
Content-Type: application/json

{
  "error": "forbidden"
}
```

---

## What You'd Do With More Time
- Bonus Tasks
- Bug fixes
- SSE Implementation
- Code clean up
- Better UX 


### 12. **Shortcuts I Took (Technical Debt)**
   - **Used ChatGPT & Copilot**
   - **Input Validation:** Minimal validation (should use validator structs)
   - **Error Handling:** Generic error messages (should be more specific)
   - **Database:** No transaction handling (batch operations could fail partially)
   - **Authentication:** No refresh tokens (tokens valid forever)
   - **Logging:** Injected everywhere (could use global package or middleware)
   - **CORS:** Hardcoded localhost (should be configurable)
   - **Pagination:** Not implemented (could load massive datasets)

### Honest Assessment

**What Went Well:**
- Clean separation of concerns makes the codebase maintainable
- Structured logging from the start will save debugging time later
- Dependency injection enables easy testing
- TypeScript + React frontend catches bugs early
- Docker setup makes deployment trivial

**What Needs Work:**
- UI Bugs
- Responsiveness
- Zero tests - high risk for refactoring
- Performance not measured - could be hiding N+1 query problems
- No background jobs - can't scale notifications/batch operations
- Frontend state management could be more sophisticated for larger apps
- Database schema could use more constraints and indexes

**If I Shipped This Today:**
✅ Users can register, log in, create projects and tasks
✅ Can assign tasks to team members
✅ Can search and filter tasks
✅ API is RESTful and well-documented
✅ Code is organized and maintainable

❌ Won't handle 1000s of concurrent users well
❌ No way to notify users of updates
❌ No audit trail of changes
❌ Limited to single-team usage
