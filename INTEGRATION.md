# Frontend-Backend Integration Guide

This guide explains how the frontend and backend are integrated for the readability study.

## Setup

### 1. Start the Backend

```bash
cd Webgazer-Backend
go run .
```

The backend will start on `http://localhost:8080` by default.

### 2. Configure Frontend (Optional)

If your backend is running on a different URL, create a `.env` file in `Webgazer-Frontend/`:

```env
VITE_API_URL=http://localhost:8080
```

### 3. Start the Frontend

```bash
cd Webgazer-Frontend
npm run dev
```

The frontend will start on `http://localhost:5173` by default.

## Data Flow

### 1. Participant Creation

- When a user first visits, a participant is automatically created (or reused from sessionStorage)
- Participant ID is stored in `sessionStorage` for the session

### 2. Calibration

- Calibration points are stored in `sessionStorage` as `calibration_points`
- This data is sent when the quiz is submitted

### 3. Reading Session

- Font preferences and reading times are stored in `sessionStorage`:
  - `font_left`, `font_right`
  - `time_left_ms`, `time_right_ms`
  - `timeA_ms`, `timeB_ms`
  - `font_preference`, `font_preferred_type`

### 4. Quiz Submission

- When the user submits the quiz, all data is collected from `sessionStorage`
- Data is sent to `/api/session` endpoint
- Quiz responses are included as JSON in `quiz_responses_json` field

## API Endpoints

### POST `/api/participant`

Creates a new participant.

**Request:**

```json
{
  "source": "web"
}
```

**Response:**

```json
{
  "success": true,
  "id": 1,
  "source": "web"
}
```

### POST `/api/session`

Creates a study session with all collected data.

**Request:**

```json
{
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
  "quiz_responses_json": "[{\"question_id\":\"q1\",\"answer\":1}]",
  "user_agent": "...",
  "screen_width": 1920,
  "screen_height": 1080
}
```

**Response:**

```json
{
  "success": true,
  "session_id": "abc123...",
  "id": 1
}
```

## Testing the Integration

1. Start both servers (backend and frontend)
2. Complete the study flow:
   - Calibration → Accuracy → Reading → Quiz
3. Submit the quiz
4. Check the backend database:
   ```bash
   sqlite3 Webgazer-Backend/readability.db
   SELECT * FROM study_sessions;
   SELECT * FROM participants;
   ```

## Troubleshooting

### CORS Errors

- Ensure the backend CORS settings include your frontend URL
- Check that both servers are running

### Data Not Saving

- Check browser console for errors
- Verify backend is running and accessible
- Check network tab in browser dev tools

### API Connection Issues

- Verify `VITE_API_URL` is set correctly (or defaults to `http://localhost:8080`)
- Check backend logs for errors
- Ensure backend port matches frontend configuration
