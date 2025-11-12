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

## 2. Create Participant

```bash
curl -X POST http://localhost:8080/api/participant \
  -H "Content-Type: application/json" \
  -d '{"source": "test"}'
```

Expected: `{"success":true,"id":1,"source":"test"}`

**Save the participant ID** for next steps (e.g., `PARTICIPANT_ID=1`)

## 3. Create Study Session

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

## 4. Submit Quiz Response

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

## 5. Submit Calibration Data

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

## 6. Submit Accuracy Measurement

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

## 7. Submit Gaze Point

```bash
curl -X POST http://localhost:8080/api/gaze-point \
  -H "Content-Type: application/json" \
  -d '{
    "session_id": 1,
    "x": 500.2,
    "y": 300.8,
    "panel": "A",
    "phase": "middle"
  }'
```

Expected: `{"success":true,"id":1}`

## 8. Submit Reading Event

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

## Verify Data

After testing, check the database:

```bash
./view-db.sh
```

Or query directly:

```bash
sqlite3 readability.db "SELECT * FROM quiz_responses;"
sqlite3 readability.db "SELECT * FROM calibration_data;"
sqlite3 readability.db "SELECT * FROM accuracy_measurements;"
sqlite3 readability.db "SELECT * FROM gaze_points;"
sqlite3 readability.db "SELECT * FROM reading_events;"
```

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
