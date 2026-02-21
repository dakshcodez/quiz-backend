# Quiz App Backend — Phase 3

A REST API for a quiz application, built with Go and the Gin framework. Phase 3 introduces **real scoring logic**, **submission tracking**, and clearer **separation of business logic** from HTTP handling, while staying in-memory with no authentication.

---

## Short Overview

This phase is the third step in the workshop: **Phase 1** gave you student endpoints and in-memory storage; **Phase 2** added full teacher CRUD. **Phase 3** focuses on making the quiz submission flow real: compare student answers to correct answers, compute an actual score, store submissions in memory, and optionally expose an endpoint to view results. You also introduce a small **scoring service** (or similar) so business rules live outside handlers, making the code easier to test and extend.

---

## What Phase 2 Achieved

Phase 2 expanded the API with **teacher endpoints**:

- Create a quiz (`POST /teacher/create_quiz`)
- Add a question to a quiz (`POST /teacher/add_question/:quiz_id`)
- Update a question (`PUT /teacher/update_question/:question_id`)
- Delete a question (`DELETE /teacher/delete_question/:question_id`)
- View a full quiz with correct answers (`GET /teacher/view_quiz/:quiz_id`)
- List all quizzes (`GET /teacher/all_quizzes`)

All of this used in-memory storage with `sync.RWMutex` and safe slice updates. Student endpoints remained: view quiz (no answers) and submit quiz (mock scoring only).

---

## What Phase 3 Adds

- **Real scoring logic** — `POST /student/give_quiz` compares each submitted answer to the correct answer for that question and computes the actual score (e.g. 1 point per correct answer).
- **Submission tracking** — Each quiz submission is stored in memory (e.g. quiz ID, score, total, timestamp). Data resets on server restart.
- **Business logic separation** — Scoring and answer comparison live in a dedicated layer (e.g. a scoring service or package), so handlers stay thin and logic can be tested independently.
- **Optional results endpoint** — An endpoint such as `GET /student/results` or `GET /student/submissions` can return a list of stored submissions (or a single result by ID) for teaching and debugging.
- **Structured response payloads** — The give_quiz response includes `score`, `total`, and optionally `submission_id` or similar, with clear JSON shapes.

---

## Learning Objectives (Phase 3)

By the end of this phase, you will have practiced:

- **Writing backend business logic** — Implementing rules (e.g. “one point per correct answer”) in a dedicated place instead of inside the handler.
- **Working with maps and slices** — Storing submissions (e.g. slice of submissions or map by ID), iterating over questions and answers, and keeping updates thread-safe with the existing RWMutex.
- **Comparing submitted answers** — Looking up the correct answer per question (e.g. by question ID), normalizing option format (e.g. "A" vs "a"), and counting correct answers.
- **Structuring scoring logic** — Inputs (quiz, submitted answers) → outputs (score, total, pass/fail if desired) in a clear, testable function or service.
- **Maintaining thread safety** — Ensuring submission storage uses the same mutex (or a consistent locking strategy) so concurrent submissions do not race.
- **Designing response payloads** — Returning JSON that includes score, total, and optionally submission metadata (e.g. ID, timestamp) for clients and for a possible results endpoint.

---

## Tech Stack

| Layer          | Choice              |
| -------------- | ------------------- |
| Language       | Go                  |
| HTTP framework | Gin                 |
| Storage        | In-memory (no DB)   |
| Concurrency    | sync.RWMutex        |
| Business logic | Service layer (e.g. internal/services or internal/scoring) |

---

## Folder Structure

Phase 3 may introduce a **services** (or **scoring**) package for business logic. A typical layout:

```
quiz-backend/
├── cmd/
│   └── server/
│       └── main.go                  # Entry point: wires store, services, routes
├── internal/
│   ├── handlers/
│   │   ├── student_handler.go       # Student: view quiz, give quiz (uses scoring)
│   │   └── teacher_handler.go       # Teacher: CRUD quizzes and questions
│   ├── models/
│   │   └── quiz.go                  # Quiz, Question, Submission, request/response structs
│   ├── routes/
│   │   └── routes.go                # Student and teacher routes
│   ├── services/
│   │   └── scoring.go               # Scoring logic: compare answers, compute score
│   └── store/
│       └── memory_store.go          # Quizzes + submissions (or separate submission store)
├── .env.example
├── go.mod
└── README.md
```

If the scoring logic is small, it might live in a single file under `internal/services` (e.g. `scoring.go`) or inside the store; the important idea is that **handlers call into a clear “scoring” abstraction** rather than inlining all logic.

---

## How to Run the Project

**Prerequisites:** Go 1.21+ installed.

```bash
cd quiz-backend
go run ./cmd/server
```

Server listens on port **8080** by default. Use `PORT=3000 go run ./cmd/server` to change it.

---

## Endpoints Overview

### Student Endpoints

| Method | Path                          | Description                                      |
| ------ | ----------------------------- | ------------------------------------------------ |
| GET    | `/student/view_quiz/:quiz_id` | View quiz and questions (no correct answers)     |
| POST   | `/student/give_quiz`          | Submit answers; **real scoring**; returns score and total |
| GET    | `/student/results`            | *(Optional)* List stored quiz submissions       |

### Teacher Endpoints

| Method | Path                                    | Description                              |
| ------ | --------------------------------------- | ---------------------------------------- |
| POST   | `/teacher/create_quiz`                   | Create a new quiz                         |
| POST   | `/teacher/add_question/:quiz_id`        | Add a question to a quiz                  |
| PUT    | `/teacher/update_question/:question_id` | Update a question by ID                   |
| DELETE | `/teacher/delete_question/:question_id` | Delete a question by ID                   |
| GET    | `/teacher/view_quiz/:quiz_id`           | View full quiz (includes correct answers) |
| GET    | `/teacher/all_quizzes`                  | List all quizzes                          |

---

## Example Requests and Responses (Phase 3)

### POST /student/give_quiz (real scoring)

**Request**

| Field   | Value                                      |
| ------- | ------------------------------------------ |
| Method  | POST                                       |
| URL     | `http://localhost:8080/student/give_quiz`  |
| Headers | `Content-Type: application/json`           |
| Body    | Raw → JSON                                 |

**Request body:**
```json
{
  "quiz_id": "quiz1",
  "answers": {
    "question1": "B",
    "question2": "C"
  }
}
```

- `quiz_id` — ID of the quiz being submitted.
- `answers` — Object mapping **question ID** to the chosen **option letter** (e.g. `"A"`, `"B"`, `"C"`, `"D"`). Only provided question IDs are scored; missing questions are typically counted as wrong.

**cURL:**
```bash
curl -X POST http://localhost:8080/student/give_quiz \
  -H "Content-Type: application/json" \
  -d '{"quiz_id":"quiz1","answers":{"question1":"B","question2":"C"}}'
```

**Example response (200 OK)**

For quiz1 with two questions, answering both correctly (e.g. question1 → "B", question2 → "C"):

```json
{
  "quiz_id": "quiz1",
  "score": 2,
  "total": 2,
  "message": "Submission recorded"
}
```

If one answer is wrong (e.g. question1 → "A", question2 → "C"):

```json
{
  "quiz_id": "quiz1",
  "score": 1,
  "total": 2,
  "message": "Submission recorded"
}
```

If submission storage is implemented and an ID is returned:

```json
{
  "quiz_id": "quiz1",
  "submission_id": "sub_abc123",
  "score": 2,
  "total": 2,
  "message": "Submission recorded"
}
```

- **404** — Quiz not found.
- **400** — Invalid JSON or missing `quiz_id` / `answers`.

---

### GET /student/results (optional)

If an endpoint to view stored submissions is added:

| Field  | Value                                    |
| ------ | ---------------------------------------- |
| Method | GET                                      |
| URL    | `http://localhost:8080/student/results`  |

**Example response (200 OK)**

```json
[
  {
    "id": "sub_abc123",
    "quiz_id": "quiz1",
    "score": 2,
    "total": 2,
    "submitted_at": "2025-02-21T12:00:00Z"
  },
  {
    "id": "sub_def456",
    "quiz_id": "quiz1",
    "score": 1,
    "total": 2,
    "submitted_at": "2025-02-21T12:05:00Z"
  }
]
```

Exact field names (e.g. `id` vs `submission_id`, `submitted_at` vs `created_at`) may follow your models.

---

## Scoring Logic (High Level)

1. **Validate input** — Ensure `quiz_id` and `answers` are present; optionally validate that the quiz exists before scoring.
2. **Load quiz** — Fetch the quiz (with questions and correct answers) from the store.
3. **Compare answers** — For each question in the quiz, get the student’s answer from `answers[question_id]`. Compare it to the question’s `correct_answer` (e.g. normalize to uppercase so "a" and "A" both match "A"). If they match, add one point.
4. **Compute totals** — **Score** = number of correct answers. **Total** = number of questions in the quiz (or number of questions that were submitted, depending on your rule).
5. **Persist submission (optional)** — Store quiz_id, score, total, and optionally an ID and timestamp in memory (e.g. slice or map), protected by the same mutex as the rest of the store.
6. **Respond** — Return JSON with `score`, `total`, and optionally `submission_id` and a short message.

All of steps 2–5 can live inside a **scoring service** (or similar) that receives the quiz and the answers map and returns score, total, and any submission record; the handler then writes the submission to the store (if applicable) and returns the HTTP response.

---

## Known Limitations

- **Still in-memory** — Quizzes and submissions are lost when the server restarts. No persistence.
- **No authentication** — Anyone can call student and teacher endpoints. No concept of “current user” or “teacher vs student” identity.
- **No RBAC** — No role-based access control; any client can create quizzes, add questions, or submit answers.
- **Data resets on restart** — All created quizzes, questions, and submissions disappear after a restart. The only data that returns is the seeded quiz from startup.

---

## What’s Coming in Phase 4

Phase 4 is planned to move toward a production-style design:

- **Authentication** — Identify who is calling the API (e.g. login, tokens or sessions).
- **Middleware** — Cross-cutting concerns (logging, auth checks, request ID) in HTTP middleware.
- **RBAC** — Restrict teacher endpoints to teachers and student endpoints to students (or appropriate roles).
- **Production-ready architecture** — Clear separation of config, logging, and error handling; possibly health checks and graceful shutdown.
- **Possibly persistent storage** — Database or file-based persistence so quizzes and results survive restarts.

---

## Summary

Phase 3 replaces mock scoring with **real scoring**: compare submitted answers to correct answers, compute score and total, and optionally store submissions in memory and expose a results endpoint. By moving scoring into a dedicated service (or package), you keep handlers thin and business logic testable. The API remains in-memory with no auth or RBAC, setting the stage for Phase 4’s authentication and production concerns.
