# System Overview: Web-Based Readability Study with Eye-Tracking

## **Purpose**

This is a web-based readability study application that uses eye-tracking technology (WebGazer.js) to measure how different fonts affect reading performance and comprehension. The study employs a systematic pairwise comparison method to evaluate popular reading fonts through a tournament-style bracket system.

---

## **Core Functionality**

### **Study Flow**

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

### **Font Comparison System**

- **6 Fonts Evaluated**:

  - **Serif**: Georgia, Times New Roman, Merriweather
  - **Sans-serif**: Inter, Open Sans, Roboto

- **Tournament-Style Bracket System**:
  - **Winners Bracket**: Primary comparison path for preferred fonts
  - **Losers Bracket**: Secondary path for non-preferred fonts
  - **Elimination Logic**: Fonts that are non-preferred twice are eliminated
  - **Grand Final**: Winners bracket champion vs. Losers bracket champion
  - **Bracket Reset**: If losers champion wins, a second grand final is required

---

## **Technology Stack**

### **Frontend**

- **Framework**: SvelteKit 2.x with TypeScript
- **Build Tool**: Vite 7.x
- **Styling**: Tailwind CSS 4.x with Typography plugin
- **Eye-Tracking**: WebGazer.js 3.4.0
- **State Management**: Svelte stores (writable)
- **Testing**: Vitest
- **Linting**: ESLint + Prettier

### **Backend**

- **Language**: Go 1.21+
- **Web Framework**: Echo v4
- **ORM**: GORM
- **Database**: SQLite
- **Port**: 8080 (configurable via `PORT` environment variable)

---

## **Architecture**

### **Frontend Structure**

```
Webgazer-Frontend/
├── src/
│   ├── routes/              # SvelteKit pages
│   │   ├── +page.svelte     # Home/landing page
│   │   ├── setup/           # Setup instructions
│   │   ├── calibrate/       # Calibration page
│   │   ├── accuracy/        # Accuracy validation
│   │   ├── read/            # Reading/comparison page
│   │   ├── quiz/            # Quiz page
│   │   └── admin/           # Admin interface
│   ├── lib/
│   │   ├── components/      # Reusable components
│   │   │   ├── calibration/ # CalibrationGrid, ProgressBar
│   │   │   ├── accuracy/    # AccuracyMeasurer, GazeOverlay
│   │   │   ├── reading/     # ReadingPanel
│   │   │   ├── quiz/        # QuizQuestion, QuizResults
│   │   │   └── shared/      # Modal, WebGazerManager
│   │   ├── stores/          # State management
│   │   │   └── webgazer.ts  # WebGazer state store
│   │   ├── api.ts           # Backend API client
│   │   └── fonts.ts         # Font comparison logic
│   └── static/
│       └── webgazer.js      # Local WebGazer library
```

### **Backend Structure**

```
Webgazer-Backend/
├── main.go          # Server, routes, handlers
├── models.go        # Database models
├── seed.go          # Initial data seeding
├── utils.go         # Utility functions
└── readability.db   # SQLite database (auto-created)
```

---

## **Database Schema**

### **Core Tables**

1. **`participants`**: Participant records (source tracking)
2. **`study_sessions`**: Main session records linking all data
3. **`calibration_data`**: 25 calibration point clicks per session
4. **`accuracy_measurements`**: Post-calibration accuracy checks
5. **`gaze_points`**: Eye-tracking coordinates during reading
6. **`reading_events`**: Start/pause/resume/complete events
7. **`quiz_responses`**: Individual quiz answers with timing
8. **`study_texts`**: Reading passages (versioned)
9. **`passages`**: Multiple passages per study text
10. **`quiz_questions`**: Questions linked to study texts or passages

### **Relationships**

- Participant → StudySessions (one-to-many)
- StudySession → CalibrationData, AccuracyMeasurements, QuizResponses, GazePoints, ReadingEvents (one-to-many)
- StudyText → Passages, QuizQuestions (one-to-many)
- Passage → QuizQuestions (one-to-many)

---

## **Data Collection**

### **Calibration Data**

- 25 calibration points
- 5 clicks per point
- X/Y coordinates and timestamps

### **Gaze Tracking**

- Continuous collection during reading (every 100ms)
- X/Y coordinates
- Panel identification (A/B, left/right)
- Phase tracking (start/middle/end)
- Batched submission (10 points per batch)

### **Reading Metrics**

- Reading duration per panel
- Start/pause/resume/complete events
- Font preferences at each comparison
- Screen dimensions and user agent

### **Quiz Data**

- Question responses with answer indices
- Correctness flags
- Response times
- Timestamps

---

## **API Endpoints**

### **Public Endpoints**

- `POST /api/participant` - Create participant
- `POST /api/session` - Save study session
- `POST /api/calibration` - Save calibration point
- `POST /api/accuracy` - Save accuracy measurement
- `POST /api/gaze-point` - Save gaze data point
- `POST /api/reading-event` - Save reading event
- `POST /api/quiz-response` - Save quiz answer
- `GET /api/study-text` - Get active study text
- `GET /api/quiz-questions` - Get quiz questions
- `GET /api/health` - Health check

### **Admin Endpoints** (`/api/admin/`)

- `GET/POST/PUT /study-text` - Manage study texts
- `GET/POST/PUT/DELETE /passage` - Manage passages
- `GET/POST/PUT/DELETE /quiz-question` - Manage quiz questions
- `GET /statistics` - Study statistics dashboard

---

## **Key Features**

### **1. Multi-Passage Support**

- Study texts can contain multiple passages
- Each passage can have its own font configuration
- Separate quizzes per passage
- Progress tracking across passages

### **2. Session Management**

- Browser `sessionStorage` for client-side state
- Persistent session IDs
- Resume capability
- Automatic cleanup on completion

### **3. Real-Time Eye Tracking**

- WebGazer.js integration
- Continuous gaze collection
- Visual feedback (optional red dot indicator)
- Automatic calibration validation

### **4. Admin Interface**

- Web-based content management
- Create/edit study texts, passages, questions
- Statistics dashboard
- Version management for study texts

### **5. Font Comparison Algorithm**

- Systematic tournament bracket
- Automatic elimination logic
- Bracket reset handling
- Progress tracking

---

## **Configuration**

### **Environment Variables**

- **Backend**: `PORT` (default: 8080)
- **Frontend**: `VITE_API_URL` (default: http://localhost:8080)

### **CORS Configuration**

- Allowed origins: localhost:5173, localhost:4173, localhost:3000
- Configurable in `main.go`

---

## **Development Workflow**

### **Frontend**

```bash
npm install          # Install dependencies
npm run dev          # Development server (localhost:5173)
npm run build        # Production build
npm run check        # Type checking
npm run lint         # Linting
npm test             # Run tests
```

### **Backend**

```bash
go mod tidy          # Update dependencies
go run .             # Run server
go build -o readability-backend  # Build binary
./test-endpoints.sh  # Test all endpoints
```

---

## **Data Flow**

1. **Participant starts** → Creates participant record
2. **Setup** → Instructions displayed
3. **Calibration** → 25 points clicked, data saved
4. **Accuracy check** → Validation measurement saved
5. **Reading session**:
   - Study text fetched from backend
   - Font comparison initialized
   - Gaze data collected continuously
   - Reading events tracked
   - Preferences recorded
6. **Quiz** → Questions fetched, responses saved
7. **Completion** → All data submitted to backend

---

## **Security & Privacy**

- No authentication (research study)
- CORS configured for specific origins
- SQLite database (local file)
- Session data in browser storage
- Camera access required (WebGazer)

---

## **Testing**

- **Frontend**: Vitest unit tests
- **Backend**: Bash scripts (`test-endpoints.sh`, `test.sh`)
- **Manual testing**: `test-endpoints-manual.md`

---

## **Deployment**

- **Frontend**: Static build (`npm run build`) → `build/` directory
- **Backend**: Go binary (`go build`) → `readability-backend`
- **Database**: SQLite file (auto-created on first run)
- **Production**: Deploy frontend to static host, backend to server

---

## **Statistics & Analytics**

Admin statistics endpoint provides:

- Participant counts by source
- Session totals
- Font preference distribution
- Quiz performance metrics
- Reading time averages
- Accuracy measurement statistics
- Gaze point distributions
- Calibration data counts

---

## **Project Structure**

```
webgazer-project/
├── Webgazer-Frontend/       # SvelteKit frontend
│   ├── src/
│   │   ├── routes/          # Pages
│   │   ├── lib/             # Components, stores, utilities
│   │   └── static/          # Static assets
│   ├── package.json
│   └── vite.config.ts
├── Webgazer-Backend/        # Go backend
│   ├── main.go             # Server & routes
│   ├── models.go           # Database models
│   ├── seed.go             # Initial data
│   ├── go.mod
│   └── readability.db      # SQLite database
├── README.md                # Main documentation
├── LICENSE                  # GNU GPL v3.0
└── SYSTEM_OVERVIEW.md       # This file
```

---

## **License**

This project is licensed under the **GNU General Public License v3.0**. See the [LICENSE](LICENSE) file for details.

---

## **Summary**

This system is designed for academic research on font readability, using web-based eye-tracking to collect quantitative data on reading behavior and font preferences. The tournament-style comparison method ensures systematic evaluation of fonts while maintaining user engagement through interactive reading sessions and comprehension quizzes.
