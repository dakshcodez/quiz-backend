# Quiz App Backend — Phase 1

A simple REST API for a quiz application, built with Go and the Gin framework. This is Phase 1 of a backend workshop and focuses on basic API design, routing, and in-memory data handling.

---

## Short Description

This project exposes two **student** endpoints: one to view a quiz (and its questions) and one to submit answers. Data lives only in memory, so it’s easy to run and experiment with. No persistence, no auth — just HTTP and JSON.

---

## Learning Objectives (Phase 1)

By the end of this phase, you will have seen:

- How to structure a small Go HTTP server with **Gin**
- How to define **routes** and **handlers** and return **JSON**
- How to keep data in memory using a **map** and protect it with **sync.RWMutex**
- How to validate request bodies and return sensible **HTTP status codes** (e.g. 200, 404, 400)
- A clear split between **models**, **handlers**, **store**, and **routes**

---

## Tech Stack

| Layer        | Choice              |
| ------------ | ------------------- |
| Language     | Go                  |
| HTTP framework | Gin             |
| Storage      | In-memory (no DB)   |
| Concurrency  | sync.RWMutex        |

---

## Folder Structure

```
quiz-backend/
├── cmd/
│   └── server/
│       └── main.go          # Entry point: starts server and wires routes
├── internal/
│   ├── handlers/
│   │   └── student_handler.go   # View quiz & submit quiz (mock scoring)
│   ├── models/
│   │   └── quiz.go              # Quiz, Question, request/response structs
│   ├── store/
│   │   └── memory_store.go      # In-memory map + RWMutex, seed data
│   └── routes/
│       └── routes.go            # Registers GET/POST routes
├── .env.example
├── go.mod
└── README.md
```

---

## How to Run the Project

**Prerequisites:** Go 1.21+ installed.

1. Clone the repo and go into the project folder:
   ```bash
   cd quiz-backend
   ```

2. Start the server:
   ```bash
   go run ./cmd/server
   ```
   You should see the server listening (default port **8080**).

3. (Optional) Use a different port via environment variable:
   ```bash
   PORT=3000 go run ./cmd/server
   ```

The server runs until you stop it with `Ctrl+C`.

---

## Example Requests (Postman / cURL)

Base URL: `http://localhost:8080`

### 1. View a quiz

**GET** `/student/view_quiz/:quiz_id`

| Field   | Value                          |
| ------- | ------------------------------ |
| Method  | GET                            |
| URL     | `http://localhost:8080/student/view_quiz/quiz1` |

- **quiz1** is the ID of the pre-seeded quiz. Use it as in the URL above.
- For a non-existent quiz (e.g. `quiz99`), the API returns **404**.

**cURL:**
```bash
curl http://localhost:8080/student/view_quiz/quiz1
```

---

### 2. Submit quiz answers (mock scoring)

**POST** `/student/give_quiz`

| Field   | Value                                      |
| ------- | ------------------------------------------ |
| Method  | POST                                       |
| URL     | `http://localhost:8080/student/give_quiz`   |
| Headers | `Content-Type: application/json`           |
| Body    | Raw → JSON (see below)                     |

**Request body (JSON):**
```json
{
  "quiz_id": "quiz1",
  "answers": {
    "question1": "A",
    "question2": "C"
  }
}
```

- `quiz_id` and `answers` are required. Keys in `answers` are question IDs; values are option letters (e.g. `"A"`, `"B"`, `"C"`, `"D"`).
- In Phase 1 the server does **not** compute a real score; it always returns a mock response.

**cURL:**
```bash
curl -X POST http://localhost:8080/student/give_quiz \
  -H "Content-Type: application/json" \
  -d '{"quiz_id":"quiz1","answers":{"question1":"A","question2":"C"}}'
```

---

## Example JSON Responses

### GET /student/view_quiz/quiz1 (200 OK)

```json
{
  "id": "quiz1",
  "title": "Sample Quiz",
  "description": "A short quiz for the workshop",
  "questions": [
    {
      "id": "question1",
      "question_text": "What is 2 + 2?",
      "option_a": "3",
      "option_b": "4",
      "option_c": "5",
      "option_d": "6"
    },
    {
      "id": "question2",
      "question_text": "What is the capital of France?",
      "option_a": "London",
      "option_b": "Berlin",
      "option_c": "Paris",
      "option_d": "Madrid"
    }
  ]
}
```

Note: The API **does not** return the correct answers — only the question text and options.

---

### POST /student/give_quiz (200 OK) — Phase 1 mock

```json
{
  "quiz_id": "quiz1",
  "score": 5,
  "message": "Phase 1 mock scoring"
}
```

Real scoring will be added in a later phase.

---

### 404 — Quiz not found (GET /student/view_quiz/quiz99)

```json
{
  "error": "quiz not found"
}
```

---

### 400 — Invalid request (e.g. missing body or fields)

```json
{
  "error": "Key: 'GiveQuizRequest.QuizID' Error:Field validation for 'QuizID' failed on the 'required' tag"
}
```

---

## In-Memory Storage

All quiz data is kept **in memory** in a single struct:

- A **map** from quiz ID (string) to quiz (title, description, questions).
- The map is only modified at startup when we **seed** one quiz with two questions.
- Every read (e.g. when you call “view quiz”) uses **sync.RWMutex**: multiple readers can run at once; any future writer would get an exclusive lock. This keeps the map safe when many requests are handled at the same time.

So: no files, no external storage — just a map in RAM. Simple and good for learning.

---

## Known Limitations (Phase 1)

- **Data resets on restart** — Stopping and restarting the server wipes all in-memory data. The only quiz you see is the one that’s seeded when the server starts.
- **No real scoring** — `POST /student/give_quiz` only validates the JSON and returns a fixed mock response. It does not check answers or compute a score.
- **Single pre-seeded quiz** — Only one quiz (`quiz1`) with two questions exists. You cannot create new quizzes yet.
- **No teacher endpoints** — Only student-facing endpoints are implemented.

---

## What’s Coming in Phase 2

In the next phase we plan to add:

- **Teacher endpoints** — e.g. create a quiz, add questions.
- **Real scoring** — `POST /student/give_quiz` will compare answers to the correct ones and return an actual score.
- Possibly **persistent storage** so that quizzes and results survive server restarts.

---

## Summary

Phase 1 gives you a minimal but clean Quiz App backend: view one quiz, submit answers (mock response), in-memory storage with a map and RWMutex, and clear separation of routes, handlers, models, and store. Use this as the base for Phase 2 and beyond.
