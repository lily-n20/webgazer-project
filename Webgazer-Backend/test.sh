#!/bin/bash

# Test script for Readability Backend API

BASE_URL="http://localhost:8080"

echo "🧪 Testing Readability Backend API"
echo "=================================="
echo ""

# Test 1: Health Check
echo "1. Testing Health Endpoint..."
curl -s "$BASE_URL/api/health" | jq .
echo ""
echo ""

# Test 2: Create a Participant
echo "2. Creating a Participant..."
PARTICIPANT_RESPONSE=$(curl -s -X POST "$BASE_URL/api/participant" \
  -H "Content-Type: application/json" \
  -d '{
    "source": "test"
  }')
echo "$PARTICIPANT_RESPONSE" | jq .
PARTICIPANT_ID=$(echo "$PARTICIPANT_RESPONSE" | jq -r '.id')
echo "Participant ID: $PARTICIPANT_ID"
echo ""
echo ""

# Test 3: Create a Study Session
echo "3. Creating a Study Session..."
SESSION_RESPONSE=$(curl -s -X POST "$BASE_URL/api/session" \
  -H "Content-Type: application/json" \
  -d "{
    \"participant_id\": $PARTICIPANT_ID,
    \"calibration_points\": 25,
    \"font_left\": \"serif\",
    \"font_right\": \"sans\",
    \"time_left_ms\": 5000,
    \"time_right_ms\": 4500,
    \"time_a_ms\": 5000,
    \"time_b_ms\": 4500,
    \"font_preference\": \"A\",
    \"preferred_font_type\": \"serif\",
    \"user_agent\": \"test-agent\",
    \"screen_width\": 1920,
    \"screen_height\": 1080
  }")
echo "$SESSION_RESPONSE" | jq .
SESSION_ID=$(echo "$SESSION_RESPONSE" | jq -r '.id')
echo "Session ID: $SESSION_ID"
echo ""
echo ""

# Test 4: Add Calibration Data
echo "4. Adding Calibration Data..."
curl -s -X POST "$BASE_URL/api/calibration" \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": $SESSION_ID,
    \"point_index\": 0,
    \"click_number\": 1,
    \"x\": 100.5,
    \"y\": 200.3
  }" | jq .
echo ""
echo ""

# Test 5: Add Accuracy Measurement
echo "5. Adding Accuracy Measurement..."
curl -s -X POST "$BASE_URL/api/accuracy" \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": $SESSION_ID,
    \"accuracy\": 85.5,
    \"duration\": 5000,
    \"passed\": true
  }" | jq .
echo ""
echo ""

# Test 6: Add Quiz Response
echo "6. Adding Quiz Response..."
curl -s -X POST "$BASE_URL/api/quiz-response" \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": $SESSION_ID,
    \"question_id\": \"q1\",
    \"answer_index\": 1,
    \"is_correct\": true,
    \"response_time\": 3000
  }" | jq .
echo ""
echo ""

# Test 7: Add Gaze Point
echo "7. Adding Gaze Point..."
curl -s -X POST "$BASE_URL/api/gaze-point" \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": $SESSION_ID,
    \"x\": 500.2,
    \"y\": 300.8,
    \"panel\": \"A\",
    \"phase\": \"middle\"
  }" | jq .
echo ""
echo ""

# Test 8: Add Reading Event
echo "8. Adding Reading Event..."
curl -s -X POST "$BASE_URL/api/reading-event" \
  -H "Content-Type: application/json" \
  -d "{
    \"session_id\": $SESSION_ID,
    \"event_type\": \"start\",
    \"panel\": \"A\",
    \"duration\": 0
  }" | jq .
echo ""
echo ""

echo "✅ All tests completed!"
echo ""
echo "To view the data, check the SQLite database:"
echo "  sqlite3 readability.db"
echo ""
echo "Or use a SQLite browser to view the tables."

