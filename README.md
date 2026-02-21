# Quiz App Backend — Phase 2

This repository contains Phase 2 of the Quiz Backend workshop, built with Go and the Gin web framework. This phase expands the system from a student-only viewer to a full teacher-managed quiz platform with CRUD capabilities.

## What Phase 1 Covered
Phase 1 introduced a minimal student-facing API that used in-memory storage. Students could view a hardcoded quiz and submit answers to get a mock score. No administration or quiz creation features were available.

## What Phase 2 Adds
In Phase 2, we introduce **Teacher Endpoints** to manage the quiz content dynamically:
- **Full Quiz Management:** Create new quizzes with titles and descriptions.
- **Question CRUD:** Add, update, and delete questions within specific quizzes.
- **Teacher Views:** Specialized endpoints to view quizzes including correct answers and list all available quizzes.
- **Concurrency Safety:** Integration of `sync.RWMutex` to handle concurrent read/write operations on our in-memory data safely.

## Learning Objectives (Phase 2)
- **Expanding REST APIs:** Building out multiple resource-related endpoints.
- **URL Parameters:** Extracting variables like `:quiz_id` and `:question_id` from routes.
- **Nested Data Handling:** Managing complex structures where questions are segments of a quiz slice.
- **Slice Manipulation:** Implementing efficient logic to update and delete elements from Go slices.
- **Safe Concurrency:** Learning why and how to use mutexes to prevent data races.
- **Status Codes:** Returning semantic HTTP response codes (e.g., `201 Created`, `404 Not Found`, `400 Bad Request`).

## Folder Structure
```text
.
├── cmd/
│   └── server/          # Application entry point (main.go)
├── internal/
│   ├── handlers/       # HTTP request handlers for student/teacher
│   ├── models/         # Data structures (Quiz, Question, etc.)
│   ├── routes/         # Route registrations and grouping
│   └── store/          # In-memory data store logic
├── go.mod/go.sum       # Dependency management
└── README.md           # Documentation (You are here)
```

## Complete Endpoints List

### Student Routes
- `GET /student/view_quiz/:quiz_id` - View a quiz (questions only)
- `POST /student/give_quiz` - Submit answers (mock scoring)

### Teacher Routes
- `POST /teacher/create_quiz` - Create a new empty quiz
- `POST /teacher/add_question/:quiz_id` - Add a question to a quiz
- `PUT /teacher/update_question/:question_id` - Update an existing question
- `DELETE /teacher/delete_question/:question_id` - Remove a question
- `GET /teacher/view_quiz/:quiz_id` - View quiz including correct answers
- `GET /teacher/all_quizzes` - List all created quizzes

## Example Requests and Responses

### 1. Create a Quiz
**POST** `/teacher/create_quiz`

**Request Body:**
```json
{
  "id": "chem-101",
  "title": "Chemistry Basics",
  "description": "Introductory quiz on periodic table elements."
}
```

**Success Response (200 OK):**
```json
{
  "id": "chem-101",
  "title": "Chemistry Basics",
  "description": "Introductory quiz on periodic table elements.",
  "questions": []
}
```

**Error Response (400 Bad Request - ID exists):**
```json
{
  "error": "quiz id already exists"
}
```

### 2. Add a Question
**POST** `/teacher/add_question/chem-101`

**Request Body:**
```json
{
  "id": "q1",
  "question_text": "What is the atomic symbol for Gold?",
  "option_a": "Ag",
  "option_b": "Au",
  "option_c": "Gd",
  "option_d": "Gl",
  "correct_answer": "B"
}
```

**Success Response (201 Created):**
```json
{
  "id": "q1",
  "question_text": "What is the atomic symbol for Gold?",
  "option_a": "Ag",
  "option_b": "Au",
  "option_c": "Gd",
  "option_d": "Gl",
  "correct_answer": "B"
}
```

### 3. Update a Question
**PUT** `/teacher/update_question/q1`

**Request Body:**
```json
{
  "id": "q1",
  "question_text": "What is the atomic symbol for Gold? (Updated)",
  "option_a": "Ag",
  "option_b": "Au",
  "option_c": "Gd",
  "option_d": "Gl",
  "correct_answer": "B"
}
```

**Success Response (200 OK):**
```json
{
  "id": "q1",
  "question_text": "What is the atomic symbol for Gold? (Updated)",
  "option_a": "Ag",
  "option_b": "Au",
  "option_c": "Gd",
  "option_d": "Gl",
  "correct_answer": "B"
}
```

### 4. Delete a Question
**DELETE** `/teacher/delete_question/q1`

**Success Response (200 OK):**
```json
{
  "message": "question deleted"
}
```

**Error Response (404 Not Found):**
```json
{
  "error": "question not found"
}
```

### 5. View Quiz (Teacher)
**GET** `/teacher/view_quiz/chem-101`

**Success Response (200 OK):**
```json
{
  "id": "chem-101",
  "title": "Chemistry Basics",
  "description": "Introductory quiz on periodic table elements.",
  "questions": [
    {
      "id": "q1",
      "question_text": "What is the atomic symbol for Gold?",
      "option_a": "Ag",
      "option_b": "Au",
      "option_c": "Gd",
      "option_d": "Gl",
      "correct_answer": "B"
    }
  ]
}
```

### 6. List All Quizzes
**GET** `/teacher/all_quizzes`

**Success Response (200 OK):**
```json
[
  {
    "id": "chem-101",
    "title": "Chemistry Basics",
    "description": "Introductory quiz on periodic table elements.",
    "questions": []
  }
]
```

## In-Memory Storage
Data is currently stored in application memory using Go maps and slices.
> [!IMPORTANT]
> Because storage is volatile, **all data is lost when the server restarts**. Persistent database integration is planned for a future phase.

## Known Limitations
- **No Authentication:** All endpoints are public and do not require login.
- **No Role Enforcement:** Any user can call teacher endpoints.
- **No Persistence:** Data resets on server restart.
- **Mock Scoring:** Student submissions return a hardcoded score regardless of answers.
- **Single-process:** The in-memory store is local to a single running instance.

## What’s Coming in Phase 3
Phase 3 will bridge the gap between "dumb" storage and actual application logic:
- **Real Scoring Logic:** Comparing student answers against correct ones.
- **Submission Tracking:** Recording student performance.
- **Business Logic Implementation:** Enforcing rules during the quiz-taking process.
