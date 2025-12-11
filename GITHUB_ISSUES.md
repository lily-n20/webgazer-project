# Retrospective GitHub Issues

This document contains GitHub issues that would have been created during development. Use these to create issues in your GitHub repository.

## Setup & Infrastructure

### Issue #1: Project Setup - Frontend

**Title:** Set up SvelteKit frontend project with TypeScript and Tailwind CSS

**Description:**

- Initialize SvelteKit project
- Configure TypeScript
- Set up Tailwind CSS
- Create basic routing structure
- Set up API client utilities

**Acceptance Criteria:**

- [ ] SvelteKit project initialized
- [ ] TypeScript configured
- [ ] Tailwind CSS working
- [ ] Basic routes created (home, setup, calibrate, read, quiz, admin)
- [ ] API client structure in place

---

### Issue #2: Project Setup - Backend

**Title:** Set up Go backend with GORM and SQLite

**Description:**

- Initialize Go project
- Set up GORM with SQLite
- Create database models
- Set up Echo framework for routing
- Configure CORS

**Acceptance Criteria:**

- [ ] Go project initialized
- [ ] Database models defined
- [ ] Auto-migration working
- [ ] Basic API routes set up
- [ ] CORS configured

---

## Core Features

### Issue #3: WebGazer Integration

**Title:** Integrate WebGazer.js for eye-tracking

**Description:**

- Add WebGazer.js library
- Create WebGazerManager component
- Implement initialization logic
- Handle video feed display
- Manage face overlay and prediction points

**Acceptance Criteria:**

- [ ] WebGazer.js integrated
- [ ] WebGazerManager component created
- [ ] Video feed can be shown/hidden
- [ ] Face overlay works
- [ ] Prediction points configurable

---

### Issue #4: Setup Page with Instructions

**Title:** Create setup page with positioning instructions and face overlay

**Description:**

- Create `/setup` route
- Add instruction modal for lighting and positioning
- Display WebGazer face overlay (centered, large)
- Add "Ready to Start" button
- Hide WebGazer UI on home page

**Acceptance Criteria:**

- [ ] Setup page created
- [ ] Instruction modal with proper formatting
- [ ] Face overlay centered and visible
- [ ] Navigation to calibration works
- [ ] WebGazer UI hidden on home page

---

### Issue #5: Calibration System

**Title:** Implement 25-point calibration grid

**Description:**

- Create calibration grid component
- Implement 25 calibration points
- Store calibration data
- Track click coordinates
- Submit calibration data to backend

**Acceptance Criteria:**

- [ ] 25-point grid displayed
- [ ] Click tracking works
- [ ] Data stored in sessionStorage
- [ ] Calibration data submitted to backend
- [ ] Progress indicator shows completion

---

### Issue #6: Accuracy Measurement

**Title:** Implement accuracy validation after calibration

**Description:**

- Create accuracy measurement page
- Display validation points
- Calculate accuracy percentage
- Check if accuracy meets threshold
- Allow retry if accuracy is too low

**Acceptance Criteria:**

- [ ] Accuracy measurement page created
- [ ] Validation points displayed
- [ ] Accuracy calculated correctly
- [ ] Pass/fail threshold enforced
- [ ] Retry functionality works

---

### Issue #7: Double-Elimination Font Comparison Logic

**Title:** Implement systematic pairwise comparison tournament system

**Description:**

- Create tournament logic in `fonts.ts`
- Implement winners/losers bracket system
- Handle font elimination (2 losses)
- Generate next comparisons automatically
- Track final preferred font

**Acceptance Criteria:**

- [ ] Tournament logic implemented
- [ ] 6 fonts compete correctly
- [ ] Brackets advance properly
- [ ] Elimination works (2 losses)
- [ ] Final preferred font determined
- [ ] Unit tests pass

---

### Issue #8: Reading Interface

**Title:** Create reading interface with side-by-side font comparison

**Description:**

- Create reading page with two panels
- Display passages in different fonts
- Track reading start/pause/resume/complete events
- Collect gaze points during reading
- Record reading times for each panel
- Implement font preference selection

**Acceptance Criteria:**

- [ ] Two-panel reading interface
- [ ] Fonts displayed correctly
- [ ] Reading events tracked
- [ ] Gaze points collected
- [ ] Reading times recorded
- [ ] Font preference selection works

---

### Issue #9: Quiz System

**Title:** Implement comprehension quiz with multiple questions

**Description:**

- Create quiz page
- Fetch quiz questions from backend
- Display questions with multiple choice answers
- Track user responses
- Submit quiz answers to backend
- Show completion screen

**Acceptance Criteria:**

- [ ] Quiz page created
- [ ] Questions fetched from API
- [ ] Multiple choice answers displayed
- [ ] Responses tracked
- [ ] Data submitted to backend
- [ ] Completion screen shows

---

### Issue #10: Multi-Passage Support

**Title:** Support multiple passages per study with separate quizzes

**Description:**

- Update backend to support passages
- Create passage management in admin
- Update frontend to handle multiple passages
- Navigate between passages
- Link quiz questions to specific passages

**Acceptance Criteria:**

- [ ] Backend supports passages
- [ ] Admin can create/manage passages
- [ ] Frontend handles multiple passages
- [ ] Navigation between passages works
- [ ] Passage-specific quizzes work

---

## Admin Interface

### Issue #11: Admin Interface - Content Management

**Title:** Create admin interface for managing study content

**Description:**

- Create admin page at `/admin`
- Add tabs for passages and quiz questions
- Implement CRUD operations for passages
- Implement CRUD operations for quiz questions
- Link content to study texts

**Acceptance Criteria:**

- [ ] Admin page accessible
- [ ] Passages can be created/edited/deleted
- [ ] Quiz questions can be created/edited/deleted
- [ ] Content linked to study texts
- [ ] Forms validate input

---

### Issue #12: Admin Interface - Data Visualization

**Title:** Add analytics dashboard with data visualizations

**Description:**

- Create statistics API endpoint
- Add analytics tab to admin interface
- Display key metrics (participants, sessions, etc.)
- Show font preference statistics
- Visualize quiz performance
- Display reading time comparisons

**Acceptance Criteria:**

- [ ] Statistics API endpoint created
- [ ] Analytics tab added
- [ ] Key metrics displayed
- [ ] Font preferences visualized
- [ ] Quiz performance shown
- [ ] Reading times compared

---

## User Experience

### Issue #13: Loading States and Transitions

**Title:** Add loading modals and smooth transitions

**Description:**

- Add loading modal during quiz question fetch
- Add loading modal during setup/calibration
- Improve transition between pages
- Hide content during loading
- Add fade-in animations

**Acceptance Criteria:**

- [ ] Loading modals appear when needed
- [ ] Smooth transitions between pages
- [ ] Content hidden during loading
- [ ] No flickering during navigation

---

### Issue #14: Modal System

**Title:** Create reusable modal component with proper formatting

**Description:**

- Create Modal component
- Support instruction modals
- Support loading modals
- Format text properly (preserve line breaks)
- Hide background content when modal is open

**Acceptance Criteria:**

- [ ] Modal component created
- [ ] Text formatting works correctly
- [ ] Background content hidden
- [ ] Modal can be dismissible or non-dismissible
- [ ] Proper styling and animations

---

### Issue #15: Fun Statistics for Users

**Title:** Add personalized fun facts on completion screen

**Description:**

- Fetch statistics on quiz completion
- Display user-specific fun fact
- Show personalized insights (preferred font, participant number, etc.)
- Include general statistics (most popular font, total data collected, etc.)
- Randomly select from available facts

**Acceptance Criteria:**

- [ ] Fun fact displayed on completion
- [ ] User-specific facts shown
- [ ] General statistics included
- [ ] Facts are randomized
- [ ] Styled attractively

---

## Data Collection

### Issue #16: Gaze Point Collection

**Title:** Implement real-time gaze point collection and batching

**Description:**

- Collect gaze points during reading
- Batch gaze points for efficient submission
- Submit gaze points to backend
- Track gaze by panel and phase
- Handle errors gracefully

**Acceptance Criteria:**

- [ ] Gaze points collected in real-time
- [ ] Batching works correctly
- [ ] Data submitted to backend
- [ ] Panel and phase tracked
- [ ] Error handling implemented

---

### Issue #17: Session Data Management

**Title:** Implement session storage and data submission

**Description:**

- Store session data in browser sessionStorage
- Collect all data before submission
- Submit complete session to backend
- Handle session persistence across navigation
- Clear session data on completion

**Acceptance Criteria:**

- [ ] Data stored in sessionStorage
- [ ] Data persists across navigation
- [ ] Complete session submitted
- [ ] Data cleared appropriately
- [ ] Error handling for submission

---

## Testing

### Issue #18: Frontend Unit Tests

**Title:** Write unit tests for font comparison tournament logic

**Description:**

- Set up Vitest testing framework
- Write tests for tournament initialization
- Test bracket progression
- Test elimination logic
- Test final preferred font determination

**Acceptance Criteria:**

- [ ] Vitest configured
- [ ] Tournament logic tests pass
- [ ] All edge cases covered
- [ ] Tests run in CI/CD

---

### Issue #19: Backend API Testing

**Title:** Create test scripts for all API endpoints

**Description:**

- Create test script for all endpoints
- Test CRUD operations
- Test data validation
- Test error handling
- Document test results

**Acceptance Criteria:**

- [ ] Test script created
- [ ] All endpoints tested
- [ ] Edge cases covered
- [ ] Tests documented

---

## Documentation

### Issue #20: README Documentation

**Title:** Create comprehensive README with setup and deployment instructions

**Description:**

- Write detailed overview
- Document local setup
- Document production deployment
- Include license information
- Add live demo link/video

**Acceptance Criteria:**

- [ ] README includes overview
- [ ] Local setup instructions clear
- [ ] Production deployment documented
- [ ] License specified
- [ ] Demo link included

---

## Bug Fixes & Improvements

### Issue #21: Fix WebGazer UI Visibility

**Title:** Ensure WebGazer elements are properly hidden/shown on different pages

**Description:**

- Hide WebGazer UI on home page
- Show video on setup page
- Hide video on calibration page
- Position face overlay correctly
- Ensure elements don't flicker

**Acceptance Criteria:**

- [ ] WebGazer UI hidden on home
- [ ] Video shows on setup
- [ ] Video hidden on calibration
- [ ] Face overlay positioned correctly
- [ ] No flickering

---

### Issue #22: Improve Modal Background Handling

**Title:** Fix modal background content visibility

**Description:**

- Hide content behind modals
- Ensure modal is always visible
- Fix tabindex warnings
- Improve backdrop styling

**Acceptance Criteria:**

- [ ] Content hidden behind modal
- [ ] Modal always visible
- [ ] No accessibility warnings
- [ ] Backdrop styled correctly

---

### Issue #23: Fix Statistics API Endpoint

**Title:** Fix statistics endpoint SQL queries and error handling

**Description:**

- Fix SQL query syntax
- Improve error handling
- Use GORM Scan instead of raw queries
- Handle edge cases (empty data)
- Add proper logging

**Acceptance Criteria:**

- [ ] Statistics endpoint works
- [ ] No SQL errors
- [ ] Error handling robust
- [ ] Edge cases handled
- [ ] Logging added

---

## Priority Order (Suggested)

1. **High Priority:** Issues #1, #2, #3, #4, #5, #6, #7, #8, #9
2. **Medium Priority:** Issues #10, #11, #12, #13, #14, #15, #16, #17
3. **Low Priority:** Issues #18, #19, #20, #21, #22, #23
