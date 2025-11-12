# Manual Endpoint Testing Guide

Test each endpoint individually using curl commands.

## Prerequisites

Make sure the backend is running:

```bash
cd Webgazer-Backend
go run .
```

## 1. Health Check

```bash
curl http://localhost:8080/api/health
```

Expected: `{"status":"ok"}`

## 2. Fetch Study Text

```bash
curl http://localhost:8080/api/study-text
```

Expected:

```json
{
  "id": 1,
  "version": "default",
  "content": "Reading is a complex cognitive process...",
  "font_left": "serif",
  "font_right": "sans"
}
```

With version parameter:

```bash
curl http://localhost:8080/api/study-text?version=default
```

## 3. Fetch Quiz Questions

```bash
curl http://localhost:8080/api/quiz-questions
```

Expected: Array of quiz questions:

```json
[
  {
    "id": "q1",
    "prompt": "What is the purpose of this passage?",
    "choices": ["To teach advanced speed-reading", "To test font readability and comprehension", ...],
    "answer": 1
  },
  ...
]
```

With study_text_id parameter:

```bash
curl http://localhost:8080/api/quiz-questions?study_text_id=1
```

## 4. Create Participant

```bash
curl -X POST http://localhost:8080/api/participant \
  -H "Content-Type: application/json" \
  -d '{"source": "test"}'
```

Expected: `{"success":true,"id":1,"source":"test"}`

**Save the participant ID** for next steps (e.g., `PARTICIPANT_ID=1`)

## 5. Create Study Session

```bash
curl -X POST http://localhost:8080/api/session \
  -H "Content-Type: application/json" \
  -d '{
    "participant_id": 1,
    "calibration_points": 25,
    "font_left": "serif",
    "font_right": "sans",
    "time_left_ms": 5000,
    "time_right_ms": 4500,
    "time_a_ms": 5000,
    "time_b_ms": 4500,
    "font_preference": "A",
    "preferred_font_type": "serif",
    "user_agent": "test",
    "screen_width": 1920,
    "screen_height": 1080
  }'
```

Expected: `{"success":true,"session_id":"...","id":1}`

**Save the session ID** for next steps (e.g., `SESSION_ID=1`)

## 6. Submit Quiz Response

```bash
curl -X POST http://localhost:8080/api/quiz-response \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": 1,
    "question_id": "q1",
    "answer_index": 1,
    "is_correct": true,
    "response_time": 3000
  }'
```

Expected: `{"success":true,"id":1}`

## 7. Submit Calibration Data

```bash
curl -X POST http://localhost:8080/api/calibration \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": 1,
    "point_index": 0,
    "click_number": 1,
    "x": 100.5,
    "y": 200.3
  }'
```

Expected: `{"success":true,"id":1}`

## 8. Submit Accuracy Measurement

```bash
curl -X POST http://localhost:8080/api/accuracy \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": 1,
    "accuracy": 85.5,
    "duration": 5000,
    "passed": true
  }'
```

Expected: `{"success":true,"id":1}`

## 9. Submit Gaze Point

Submit a gaze point during reading session. The phase can be one of:

- `waiting` - Before reading starts
- `reading_A` - While reading Box A
- `reading_B` - While reading Box B
- `completed` - After both boxes are read

```bash
curl -X POST http://localhost:8080/api/gaze-point \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": 1,
    "x": 500.2,
    "y": 300.8,
    "panel": "A",
    "phase": "reading_A"
  }'
```

Expected: `{"success":true,"id":1}`

### Submit Multiple Gaze Points (Different Phases)

```bash
# Gaze point during reading Box A
curl -X POST http://localhost:8080/api/gaze-point \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": 1,
    "x": 500.2,
    "y": 300.8,
    "panel": "A",
    "phase": "reading_A"
  }'

# Gaze point during reading Box B
curl -X POST http://localhost:8080/api/gaze-point \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": 1,
    "x": 1200.5,
    "y": 400.2,
    "panel": "B",
    "phase": "reading_B"
  }'

# Gaze point while waiting
curl -X POST http://localhost:8080/api/gaze-point \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": 1,
    "x": 960.0,
    "y": 540.0,
    "panel": "A",
    "phase": "waiting"
  }'
```

## 10. Submit Reading Event

```bash
curl -X POST http://localhost:8080/api/reading-event \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": 1,
    "event_type": "start",
    "panel": "A",
    "duration": 0
  }'
```

Expected: `{"success":true,"id":1}`

## 11. Admin: List All Study Texts

```bash
curl http://localhost:8080/api/admin/study-text
```

Expected:

```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "version": "default",
      "content": "Reading is a complex cognitive process...",
      "font_left": "serif",
      "font_right": "sans",
      "active": true,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ]
}
```

## 12. Admin: Create Study Text

```bash
curl -X POST http://localhost:8080/api/admin/study-text \
  -H "Content-Type: application/json" \
  -d '{
    "version": "v2",
    "content": "Your new reading passage text here...",
    "font_left": "serif",
    "font_right": "sans",
    "active": false
  }'
```

Expected: `{"success":true,"id":2,"message":"Study text created successfully"}`

## 13. Admin: Update Study Text

```bash
curl -X PUT http://localhost:8080/api/admin/study-text \
  -H "Content-Type: application/json" \
  -d '{
    "id": 1,
    "content": "Updated reading passage text...",
    "font_left": "sans",
    "font_right": "serif",
    "active": true
  }'
```

Expected: `{"success":true,"id":1,"message":"Study text updated successfully"}`

## 14. Admin: Get Quiz Question

```bash
curl http://localhost:8080/api/admin/quiz-question?id=1
```

Expected:

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

## 15. Admin: Create Quiz Question

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

Expected: `{"success":true,"id":6,"message":"Quiz question created successfully"}`

## 16. Admin: Update Quiz Question

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

Expected: `{"success":true,"id":1,"message":"Quiz question updated successfully"}`

## 17. Admin: Delete Quiz Question

```bash
curl -X DELETE http://localhost:8080/api/admin/quiz-question?id=1
```

Expected: `{"success":true,"message":"Quiz question deleted successfully"}`

## Verify Data

After testing, check the database:

```bash
./view-db.sh
```

Or query directly:

```bash
sqlite3 readability.db "SELECT id, version, font_left, font_right, active, substr(content, 1, 50) as content_preview FROM study_texts;"
sqlite3 readability.db "SELECT * FROM quiz_questions;"
sqlite3 readability.db "SELECT * FROM quiz_responses;"
sqlite3 readability.db "SELECT * FROM calibration_data;"
sqlite3 readability.db "SELECT * FROM accuracy_measurements;"
sqlite3 readability.db "SELECT * FROM gaze_points;"
sqlite3 readability.db "SELECT * FROM reading_events;"
```

**Note:** The `study_texts` table now includes `font_left` and `font_right` columns.

## Using jq for Pretty Output

If you have `jq` installed, pipe responses through it:

```bash
curl -X POST http://localhost:8080/api/participant \
  -H "Content-Type: application/json" \
  -d '{"source": "test"}' | jq .
```

## Testing with Variables

You can use shell variables to chain tests:

```bash
# Create participant and save ID
PARTICIPANT_ID=$(curl -s -X POST http://localhost:8080/api/participant \
  -H "Content-Type: application/json" \
  -d '{"source": "test"}' | jq -r '.id')

# Create session using participant ID
SESSION_ID=$(curl -s -X POST http://localhost:8080/api/session \
  -H "Content-Type: application/json" \
  -d "{\"participant_id\": $PARTICIPANT_ID, \"font_preference\": \"A\"}" | jq -r '.id')

# Submit quiz response using session ID
curl -X POST http://localhost:8080/api/quiz-response \
  -H "Content-Type: application/json" \
  -d "{\"session_id\": $SESSION_ID, \"question_id\": \"q1\", \"answer_index\": 1}"
```
