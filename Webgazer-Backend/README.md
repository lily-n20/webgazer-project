# Readability Study Backend

Simple Go/GORM backend for storing readability study data.

## Setup

1. **Install dependencies:**

   ```bash
   cd Webgazer-Backend
   go mod tidy
   ```

2. **Run the server:**

   ```bash
   go run .
   ```

   The server will start on port 8080 (or the PORT environment variable if set).

3. **Database:**
   - SQLite database file `readability.db` will be created automatically
   - Tables are auto-migrated on first run

## Database Models

### Participant

- `id` - Primary key
- `source` - Source of participant (e.g., "mturk", "prolific", "internal")
- `created_at` - Timestamp

### StudySession

- Main session record linking all study data
- Links to Participant via `participant_id`
- Contains reading session metadata (fonts, timing, preferences)
- Has relationships to: CalibrationData, AccuracyMeasurement, QuizResponse, GazePoint, ReadingEvent

### CalibrationData

- Individual calibration point clicks
- Fields: `point_index`, `click_number`, `x`, `y`, `timestamp`
- Links to StudySession via `session_id`

### AccuracyMeasurement

- Accuracy check results from calibration validation
- Fields: `accuracy` (%), `duration` (ms), `passed` (bool), `timestamp`
- Links to StudySession via `session_id`

### QuizResponse

- Individual quiz answers
- Fields: `question_id`, `answer_index`, `is_correct`, `response_time`, `timestamp`
- Links to StudySession via `session_id`

### GazePoint

- Eye-tracking data points during reading
- Fields: `x`, `y`, `panel` (A/B/left/right), `phase` (start/middle/end), `timestamp`
- Links to StudySession via `session_id`

### ReadingEvent

- Reading session milestones
- Fields: `event_type` (start/pause/resume/complete), `panel`, `duration`, `timestamp`
- Links to StudySession via `session_id`

## API Endpoints

### POST `/api/session`

Save a study session. Expects JSON body with:

```json
{
  "session_id": "optional-custom-id",
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
  "quiz_responses_json": "[{\"question_id\":\"q1\",\"answer\":1},...]",
  "user_agent": "optional",
  "screen_width": 1920,
  "screen_height": 1080
}
```

### POST `/api/quiz-response`

Save an individual quiz answer.

**Request:**

```json
{
  "session_id": 1,
  "question_id": "q1",
  "answer_index": 1,
  "is_correct": true,
  "response_time": 3000
}
```

### POST `/api/calibration`

Save a calibration point click.

**Request:**

```json
{
  "session_id": 1,
  "point_index": 0,
  "click_number": 1,
  "x": 100.5,
  "y": 200.3
}
```

### POST `/api/accuracy`

Save an accuracy measurement.

**Request:**

```json
{
  "session_id": 1,
  "accuracy": 85.5,
  "duration": 5000,
  "passed": true
}
```

### POST `/api/gaze-point`

Save a gaze tracking data point.

**Request:**

```json
{
  "session_id": 1,
  "x": 500.2,
  "y": 300.8,
  "panel": "A",
  "phase": "middle"
}
```

### POST `/api/reading-event`

Save a reading session milestone.

**Request:**

```json
{
  "session_id": 1,
  "event_type": "start",
  "panel": "A",
  "duration": 0
}
```

### GET `/api/health`

Health check endpoint.

## Testing

### Quick Test

1. **Start the server:**

   ```bash
   cd Webgazer-Backend
   go run .
   ```

2. **In another terminal, test the health endpoint:**

   ```bash
   curl http://localhost:8080/api/health
   ```

   Expected response: `{"status":"ok"}`

3. **Test creating a session:**

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

4. **Test all endpoints automatically:**

   ```bash
   ./test-endpoints.sh
   ```

   This script tests all endpoints and shows success/failure for each.

5. **Test endpoints manually:**
   See `test-endpoints-manual.md` for individual curl commands to test each endpoint.

6. **Run the Go test script:**

   ```bash
   go run scripts/test.go
   ```

7. **Or use the bash test script (requires jq):**
   ```bash
   ./test.sh
   ```

### View Database

**Option 1: Quick View Script**

```bash
./view-db.sh
```

This shows a summary of all data in the database.

**Option 2: Interactive SQLite CLI**

```bash
sqlite3 readability.db
```

Then run SQL queries:

```sql
.tables                    -- List all tables
SELECT * FROM participants;
SELECT * FROM study_sessions;
.mode column              -- Better formatting
.headers on               -- Show column headers
```

**Option 3: Run Pre-written Queries**

```bash
sqlite3 readability.db < view-db.sql
```

**Option 4: GUI Tool (Recommended)**

- [DB Browser for SQLite](https://sqlitebrowser.org/) - Free, cross-platform
- [TablePlus](https://tableplus.com/) - Modern database client (paid/free tier)
- [DBeaver](https://dbeaver.io/) - Universal database tool (free)

Just open `readability.db` with any of these tools.

**Option 5: VS Code Extension**

- Install "SQLite Viewer" extension
- Right-click `readability.db` → "Open Database"

## Note

All endpoints are now implemented! The backend supports:

- ✅ Participant creation
- ✅ Study session creation
- ✅ Individual quiz responses
- ✅ Calibration data points
- ✅ Accuracy measurements
- ✅ Gaze tracking points
- ✅ Reading events

When quiz responses are submitted, they are saved both as JSON in `study_sessions.quiz_responses_json` (for backward compatibility) and as individual records in the `quiz_responses` table.

## CORS

Configured to allow requests from:

- http://localhost:5173 (Vite dev server)
- http://localhost:4173 (Vite preview)
- http://localhost:3000
