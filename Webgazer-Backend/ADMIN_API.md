# Admin API - Managing Study Text and Quiz Questions

This document explains how to create, update, and manage study text and quiz questions using the admin API endpoints.

## Quick Start: Terminal Admin Interface

The easiest way to manage study text and quiz questions is using the terminal-based admin interface.

### Prerequisites


1. **Start your backend server**:
   ```bash
   cd Webgazer-Backend
   go run .
   ```

### Using the Admin CLI

1. **Run the admin interface**:

   ```bash
   ./admin-cli.sh
   ```

3. **Navigate the menu**:
   - Enter a number (1-7) to select an option
   - Follow the prompts to enter data
   - Press Enter to continue after each operation

### Menu Options

**Study Text Management:**

- `1` - List all study texts (shows ID, version, active status, content preview)
- `2` - Create new study text (prompts for version, content, active status)
- `3` - Update study text (enter ID, then update fields)

**Quiz Question Management:**

- `4` - List all quiz questions (organized by study text)
- `5` - Create new quiz question (prompts for all fields)
- `6` - Update quiz question (enter ID, then update fields)
- `7` - Delete quiz question (with confirmation)

**Other:**

- `0` - Exit

### Troubleshooting Admin CLI

**"Cannot connect to API" error:**

- Make sure your backend is running: `cd Webgazer-Backend && go run .`
- Check that the backend is on port 8080 (or update API_URL)

**"404 page not found" error:**

- The backend needs to be restarted after adding admin endpoints
- Stop the backend (Ctrl+C) and restart: `go run .`

**"jq is required" error:**

- Install jq: `brew install jq`

---

## Manual API Usage (curl commands)

If you prefer to use curl commands directly or need to script operations, here are the manual API endpoints:

## Study Text Management

### List All Study Texts

```bash
curl http://localhost:8080/api/admin/study-text
```

Response:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "version": "default",
      "content": "Reading is a complex cognitive process...",
      "active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

### Create New Study Text

```bash
curl -X POST http://localhost:8080/api/admin/study-text \
  -H "Content-Type: application/json" \
  -d '{
    "version": "v2",
    "content": "Your new reading passage text here...",
    "active": false
  }'
```

**Note:** If `active` is set to `true`, all other study texts will be automatically deactivated.

### Update Existing Study Text

```bash
curl -X PUT http://localhost:8080/api/admin/study-text \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "content": "Updated reading passage text...",
    "active": true
  }'
```

You can update any combination of:

- `version` (string)
- `content` (string)
- `active` (boolean) - Setting to `true` will deactivate all others

### Example: Update Content Only

```bash
curl -X PUT http://localhost:8080/api/admin/study-text \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "content": "New content for the reading passage."
  }'
```

### Example: Activate a Study Text

```bash
curl -X PUT http://localhost:8080/api/admin/study-text \
  -H "Content-Type: application/json" \
  -d '{
    "id": 2,
    "active": true
  }'
```

## Quiz Question Management

### Get Single Quiz Question

```bash
curl http://localhost:8080/api/admin/quiz-question?id=1
```

Response:

```json
{
  "success": true,
  "data": {
    "id": 1,
    "study_text_id": 1,
    "question_id": "q1",
    "prompt": "What is the purpose of this passage?",
    "choices": ["Option 1", "Option 2", "Option 3", "Option 4"],
    "answer": 1,
    "order": 1
  }
}
```

### Create New Quiz Question

```bash
curl -X POST http://localhost:8080/api/admin/quiz-question \
  -H "Content-Type: application/json" \
  -d '{
    "study_text_id": 1,
    "question_id": "q6",
    "prompt": "What is the main theme?",
    "choices": [
      "Theme A",
      "Theme B",
      "Theme C",
      "Theme D"
    ],
    "answer": 0,
    "order": 6
  }'
```

**Required fields:**

- `study_text_id` (number) - ID of the study text this question belongs to
- `question_id` (string) - Unique identifier like "q1", "q2", etc.
- `prompt` (string) - The question text
- `choices` (array of strings) - Answer options
- `answer` (number) - Index of correct answer (0-based)
- `order` (number) - Display order

### Update Existing Quiz Question

```bash
curl -X PUT http://localhost:8080/api/admin/quiz-question \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "prompt": "Updated question prompt",
    "choices": [
      "New Option 1",
      "New Option 2",
      "New Option 3",
      "New Option 4"
    ],
    "answer": 2
  }'
```

You can update any combination of:

- `question_id` (string)
- `prompt` (string)
- `choices` (array of strings)
- `answer` (number)
- `order` (number)

### Example: Update Only the Prompt

```bash
curl -X PUT http://localhost:8080/api/admin/quiz-question \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "prompt": "What is the updated question?"
  }'
```

### Example: Update Only the Correct Answer

```bash
curl -X PUT http://localhost:8080/api/admin/quiz-question \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "answer": 2
  }'
```

### Delete Quiz Question

```bash
curl -X DELETE http://localhost:8080/api/admin/quiz-question?id=1
```

## Complete Workflow Example

### 1. List all study texts to find the one you want to update

```bash
curl http://localhost:8080/api/admin/study-text | jq
```

### 2. Update the study text content

```bash
curl -X PUT http://localhost:8080/api/admin/study-text \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "content": "Your updated reading passage goes here. Make sure it is engaging and appropriate for your study."
  }'
```

### 3. Get existing quiz questions for that study text

```bash
curl http://localhost:8080/api/quiz-questions?study_text_id=1 | jq
```

### 4. Update a specific quiz question

```bash
curl -X PUT http://localhost:8080/api/admin/quiz-question \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "prompt": "What is the main idea of the passage?",
    "choices": [
      "Idea A",
      "Idea B",
      "Idea C",
      "Idea D"
    ],
    "answer": 1,
    "order": 1
  }'
```

### 5. Add a new quiz question

```bash
curl -X POST http://localhost:8080/api/admin/quiz-question \
  -H "Content-Type: application/json" \
  -d '{
    "study_text_id": 1,
    "question_id": "q6",
    "prompt": "What is your opinion on the passage?",
    "choices": [
      "Very interesting",
      "Somewhat interesting",
      "Not interesting",
      "Boring"
    ],
    "answer": 0,
    "order": 6
  }'
```

## Using SQLite Directly (Alternative Method)

You can also update the database directly using SQLite:

### View Study Texts

```bash
sqlite3 readability.db "SELECT * FROM study_texts;"
```

### Update Study Text Content

```bash
sqlite3 readability.db "UPDATE study_texts SET content = 'New content here' WHERE id = 1;"
```

### View Quiz Questions

```bash
sqlite3 readability.db "SELECT id, question_id, prompt, choices, answer, \"order\" FROM quiz_questions WHERE study_text_id = 1;"
```

### Update Quiz Question

```bash
sqlite3 readability.db "UPDATE quiz_questions SET prompt = 'New question?' WHERE id = 1;"
```

**Note:** When updating `choices` directly in SQLite, make sure to use valid JSON:

```bash
sqlite3 readability.db "UPDATE quiz_questions SET choices = '[\"Option 1\",\"Option 2\",\"Option 3\",\"Option 4\"]' WHERE id = 1;"
```

## Tips

1. **Always check existing data first** - Use GET endpoints to see what's currently in the database
2. **Test with inactive study texts** - Create new study texts with `active: false` first, then activate when ready
3. **Use version numbers** - Use meaningful version strings like "v1", "v2", "experimental", etc.
4. **Order matters** - Quiz questions are displayed in `order` ascending, so set appropriate order values
5. **Answer index is 0-based** - The first choice is index 0, second is 1, etc.
