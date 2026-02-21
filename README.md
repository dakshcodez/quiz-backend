# Quiz App Backend (Production)

Production-ready Quiz API built with **Go**, **Gin**, **Supabase (PostgreSQL)**, **pgx v5**, and **JWT** with role-based access (teacher | student).

---

## How to Run

### Prerequisites

- Go 1.21+
- Supabase project (PostgreSQL)
- Environment variables (see [.env.example](.env.example))

### 1. Database

- Create a Supabase project and get the **connection string** (URI). Prefer the **connection pooler** (port 6543) for server apps.
- Run the schema if needed: in Supabase SQL Editor, run [migrations/001_schema.sql](migrations/001_schema.sql).

### 2. Environment

```bash
cp .env.example .env
# Edit .env: set DATABASE_URL and JWT_SECRET.
```

Required:

- `DATABASE_URL` — PostgreSQL URI (e.g. from Supabase).
- `JWT_SECRET` — Long random string for signing JWTs.

Optional: `PORT` (default 8080), `JWT_EXPIRY_HOURS`, `DB_TIMEOUT_SEC`.

### 3. Start server

```bash
go run ./cmd/server
```

Server listens on `:8080` (or `PORT`). Health check: `GET /health`.

---

## Architecture

- **Handler → Service → Repository → Database.** No DB or JWT logic in handlers.
- **Dependency injection:** repos and services are constructed in `main.go` and passed into handlers.
- **Versioned API:** `/api/v1/...` with RBAC (teacher vs student).

---

## Example cURL Commands

Base URL: `http://localhost:8080`. Replace `TOKEN` with the JWT from login.

### Auth (public)

**Register (teacher)**

```bash
curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"teacher@example.com","password":"secret123","role":"teacher"}'
```

**Register (student)**

```bash
curl -s -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"student@example.com","password":"secret123","role":"student"}'
```

**Login**

```bash
curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"teacher@example.com","password":"secret123"}'
```

Response includes `token` and `user` (id, email, role). Use `token` in `Authorization: Bearer <token>`.

### Teacher (Bearer token required)

**Create quiz**

```bash
curl -s -X POST http://localhost:8080/api/v1/teacher/quizzes \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{"title":"Math Quiz","description":"Basic math"}'
```

**Add question**

```bash
# Replace QUIZ_ID with id from create quiz response
curl -s -X POST "http://localhost:8080/api/v1/teacher/quizzes/QUIZ_ID/questions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{
    "question_text":"What is 2+2?",
    "option_a":"3","option_b":"4","option_c":"5","option_d":"6",
    "correct_answer":"B"
  }'
```

**List my quizzes**

```bash
curl -s http://localhost:8080/api/v1/teacher/quizzes \
  -H "Authorization: Bearer TOKEN"
```

**Get quiz submissions**

```bash
curl -s "http://localhost:8080/api/v1/teacher/quizzes/QUIZ_ID/submissions" \
  -H "Authorization: Bearer TOKEN"
```

### Student (Bearer token required)

**List quizzes**

```bash
curl -s http://localhost:8080/api/v1/quizzes \
  -H "Authorization: Bearer TOKEN"
```

**Get quiz (no correct answers)**

```bash
curl -s "http://localhost:8080/api/v1/quizzes/QUIZ_ID" \
  -H "Authorization: Bearer TOKEN"
```

**Submit quiz**

```bash
curl -s -X POST "http://localhost:8080/api/v1/quizzes/QUIZ_ID/submit" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer TOKEN" \
  -d '{"answers":{"QUESTION_ID_1":"B","QUESTION_ID_2":"C"}}'
```

**My submissions**

```bash
curl -s http://localhost:8080/api/v1/student/submissions \
  -H "Authorization: Bearer TOKEN"
```

---

## Example JSON Responses

### POST /auth/login (200)

```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "email": "teacher@example.com",
    "role": "teacher"
  }
}
```

### POST /api/v1/quizzes/:quizId/submit (200)

```json
{
  "quiz_id": "uuid",
  "score": 3,
  "total": 5,
  "details": [
    { "question_id": "uuid-1", "correct": true },
    { "question_id": "uuid-2", "correct": false },
    { "question_id": "uuid-3", "correct": true }
  ]
}
```

Correct answers are never exposed in the student API.

### Error response (4xx / 5xx)

```json
{
  "code": 401,
  "message": "invalid or expired token"
}
```

---

## Project layout

```
cmd/server/main.go
internal/
  config/       # Env config
  database/     # pgx pool
  models/       # Domain and request/response structs
  repository/   # DB access (user, quiz, submission)
  services/     # Auth, quiz + scoring logic
  handlers/     # HTTP handlers
  middleware/   # JWT auth, RequireRole
  routes/       # Route setup
  utils/        # Response helpers, context keys
migrations/     # SQL schema
```

---

## License

MIT.
