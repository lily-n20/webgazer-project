# Readability Study with WebGazer

**License:** [GNU General Public License v3.0](LICENSE)

**Live Demo:** [Add your deployment URL here]

## Overview

This application conducts a web-based readability study that uses eye-tracking technology (WebGazer.js) to measure how different fonts affect reading performance and comprehension. The study employs a systematic pairwise comparison method to evaluate popular reading fonts.

### Study Flow

Participants progress through the following stages:

1. **Setup Phase**: Users receive instructions on proper positioning and lighting conditions. A face overlay helps them center themselves in the camera frame before beginning.

2. **Calibration**: WebGazer eye-tracking is calibrated using 25 calibration points displayed across the screen. Users click on each point to establish gaze tracking accuracy.

3. **Accuracy Validation**: After calibration, users complete an accuracy check to verify the eye-tracking system is working correctly. This ensures data quality before proceeding.

4. **Font Comparison Process**: The core of the study uses a systematic pairwise comparison method:

   - Six fonts are evaluated: Georgia, Times New Roman, Merriweather (serif), Inter, Open Sans, Roboto (sans-serif)
   - Fonts are paired and displayed side-by-side
   - Users read passages in each font pair and select their preference
   - Preferred fonts advance in the primary comparison path, non-preferred fonts continue in the secondary comparison path
   - Fonts are excluded from further consideration after being non-preferred twice
   - The process continues until a final preferred font is determined through the final comparison

5. **Reading Sessions**: During font comparisons, the system tracks:

   - Eye gaze coordinates in real-time
   - Reading start/pause/resume/complete events
   - Time spent reading each font panel
   - Font preferences at each comparison stage

6. **Comprehension Quiz**: After reading passages, users answer comprehension questions to verify understanding and measure reading effectiveness.

7. **Data Submission**: All collected data (gaze points, timing, preferences, quiz responses) is submitted to the backend for analysis.

### Technical Architecture

**Frontend (SvelteKit + TypeScript)**:

- Manages user interface and study flow
- Integrates WebGazer.js for eye-tracking
- Implements systematic pairwise comparison logic
- Handles real-time gaze data collection
- Stores session data in browser sessionStorage
- Communicates with backend via REST API

**Backend (Go + GORM + SQLite)**:

- RESTful API for data storage and retrieval
- SQLite database for persistent storage
- Auto-migration of database schema
- Admin endpoints for managing study content
- CORS configuration for frontend communication

**Data Collection**:

- Calibration data: 25 point clicks with coordinates
- Gaze tracking: Continuous eye position coordinates during reading
- Reading events: Start, pause, resume, complete timestamps
- Timing metrics: Reading duration for each font panel
- Font preferences: User selections at each comparison stage
- Quiz responses: Comprehension question answers with response times
- Accuracy measurements: Calibration validation results

### Key Features

- **Systematic Pairwise Comparison**: Methodical evaluation of 6 fonts with structured comparison rules
- **Real-Time Eye Tracking**: Continuous gaze data collection during reading sessions
- **Multi-Passage Support**: Study texts can contain multiple passages with separate quizzes
- **Admin Interface**: Web-based admin panel for managing study content
- **Session Management**: Browser-based session storage with backend persistence
- **Comprehensive Data Collection**: All user interactions and eye movements recorded

## Running Locally

### Prerequisites

- **Node.js** 18+ and npm
- **Go** 1.21+
- Modern web browser with camera access

### Setup

**1. Start Backend:**

```bash
cd Webgazer-Backend
go mod tidy
go run .
```

Backend runs on `http://localhost:8080` by default.

**2. Start Frontend:**

```bash
cd Webgazer-Frontend
npm install
npm run dev
```

Frontend runs on `http://localhost:5173` by default.

**3. Access Application:**

Open `http://localhost:5173` in your browser and allow camera access when prompted.

## Production Deployment

Production deployment will be handled during lab sessions. The application is ready for deployment with the following build commands:

**Backend Build:**

```bash
cd Webgazer-Backend
go build -o readability-backend
```

**Frontend Build:**

```bash
cd Webgazer-Frontend
npm install
npm run build
```

The production build outputs static files in `Webgazer-Frontend/build/` for frontend deployment and a Go binary for backend deployment.

## Project Structure

```
Webgazer-Frontend/
  src/routes/          # Pages (setup, calibrate, read, quiz, admin)
  src/lib/
    components/        # UI components
    api.ts            # API client
    fonts.ts          # Comparison logic

Webgazer-Backend/
  main.go             # Server & routes
  models.go           # Database models
  seed.go             # Initial data
  readability.db      # SQLite database
```

## Testing

**Frontend:**

```bash
cd Webgazer-Frontend
npm test
npm run test:run
```

**Backend:**

```bash
cd Webgazer-Backend
./test-endpoints.sh
```

## Database

SQLite database: `Webgazer-Backend/readability.db`

**Tables:** participants, study_sessions, calibration_data, gaze_points, quiz_responses, reading_events, accuracy_measurements, study_texts, passages, quiz_questions

**View:**

```bash
cd Webgazer-Backend
sqlite3 readability.db
```

## Configuration

**Environment Variables:**

- Backend: `PORT` (default: 8080)
- Frontend: `VITE_API_URL` (default: http://localhost:8080)

Create `.env` in `Webgazer-Frontend/`:

```
VITE_API_URL=http://localhost:8080
```

**CORS:** Configured for localhost:5173, localhost:4173, localhost:3000. Update `Webgazer-Backend/main.go` if needed.

## Troubleshooting

**WebGazer not working:**

- Grant camera permissions
- Check browser console
- Ensure proper lighting
- Refresh page

**Connection issues:**

- Verify backend on port 8080
- Check `VITE_API_URL`
- Review CORS settings

**Database reset:**

- Delete `readability.db` (auto-recreated on start)

**Session storage:**

- Data in browser `sessionStorage`
- Persists across navigation
- Clears on browser close
- Submitted when quiz completes

## Development

**Frontend:**

```bash
npm run dev      # Development
npm run build    # Production
npm run check    # Type check
npm run lint     # Lint
```

**Backend:**

```bash
go run .         # Run server
go mod tidy      # Update deps
```

## Admin Interface

Access at `/admin` to manage study texts, passages, quiz questions, and view data.

## License

This project is licensed under the GNU General Public License v3.0. See the [LICENSE](LICENSE) file for details.

## Documentation

- Integration Guide: `INTEGRATION.md`
- Backend API: `Webgazer-Backend/README.md`
- Admin API: `Webgazer-Backend/ADMIN_API.md`
