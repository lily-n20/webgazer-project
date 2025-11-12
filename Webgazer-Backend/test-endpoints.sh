#!/bin/bash

# Test script for individual API endpoints
BASE_URL="http://localhost:8080"

echo "🧪 Testing Readability Backend API Endpoints"
echo "=============================================="
echo ""

# Check if server is running
if ! curl -s "${BASE_URL}/api/health" > /dev/null 2>&1; then
    echo -e "${RED}❌ Error: Backend server is not running on ${BASE_URL}${NC}"
    echo "Please start the server first:"
    echo "  cd Webgazer-Backend"
    echo "  go run ."
    exit 1
fi
echo -e "${GREEN}✓ Backend server is running${NC}"
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test function
test_endpoint() {
    local name=$1
    local method=$2
    local endpoint=$3
    local data=$4
    
    echo -e "${YELLOW}Testing: ${name}${NC}"
    echo "  Endpoint: ${method} ${endpoint}"
    
    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "${BASE_URL}${endpoint}")
    else
        response=$(curl -s -w "\n%{http_code}" -X "${method}" \
            -H "Content-Type: application/json" \
            -d "${data}" \
            "${BASE_URL}${endpoint}")
    fi
    
    http_code=$(echo "$response" | tail -n1)
    # Use sed to remove last line (works on both Linux and macOS)
    body=$(echo "$response" | sed '$d')
    
    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        echo -e "  ${GREEN}✓ Success (${http_code})${NC}"
        echo "  Response: $body" | jq . 2>/dev/null || echo "  Response: $body"
    else
        echo -e "  ${RED}✗ Failed (${http_code})${NC}"
        echo "  Response: $body"
    fi
    echo ""
}

# Test 1: Health Check
test_endpoint "Health Check" "GET" "/api/health" ""

# Test 2: Create Participant
PARTICIPANT_DATA='{"source": "test"}'
PARTICIPANT_RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "$PARTICIPANT_DATA" \
    "${BASE_URL}/api/participant")
PARTICIPANT_ID=$(echo "$PARTICIPANT_RESPONSE" | jq -r '.id' 2>/dev/null)
test_endpoint "Create Participant" "POST" "/api/participant" "$PARTICIPANT_DATA"

if [ -z "$PARTICIPANT_ID" ] || [ "$PARTICIPANT_ID" = "null" ]; then
    echo -e "${RED}⚠ Warning: Could not get participant ID. Using ID 1 for remaining tests.${NC}"
    PARTICIPANT_ID=1
fi
echo "  Using Participant ID: $PARTICIPANT_ID"
echo ""

# Test 3: Create Study Session
SESSION_DATA="{
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
    \"user_agent\": \"test-script\",
    \"screen_width\": 1920,
    \"screen_height\": 1080
}"
SESSION_RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "$SESSION_DATA" \
    "${BASE_URL}/api/session")
SESSION_ID=$(echo "$SESSION_RESPONSE" | jq -r '.id' 2>/dev/null)
test_endpoint "Create Study Session" "POST" "/api/session" "$SESSION_DATA"

if [ -z "$SESSION_ID" ] || [ "$SESSION_ID" = "null" ]; then
    echo -e "${RED}⚠ Warning: Could not get session ID. Using ID 1 for remaining tests.${NC}"
    SESSION_ID=1
fi
echo "  Using Session ID: $SESSION_ID"
echo ""

# Test 4: Submit Quiz Response
QUIZ_RESPONSE_DATA="{
    \"session_id\": $SESSION_ID,
    \"question_id\": \"q1\",
    \"answer_index\": 1,
    \"is_correct\": true,
    \"response_time\": 3000
}"
test_endpoint "Submit Quiz Response" "POST" "/api/quiz-response" "$QUIZ_RESPONSE_DATA"

# Test 5: Submit Calibration Data
CALIBRATION_DATA="{
    \"session_id\": $SESSION_ID,
    \"point_index\": 0,
    \"click_number\": 1,
    \"x\": 100.5,
    \"y\": 200.3
}"
test_endpoint "Submit Calibration Data" "POST" "/api/calibration" "$CALIBRATION_DATA"

# Test 6: Submit Accuracy Measurement
ACCURACY_DATA="{
    \"session_id\": $SESSION_ID,
    \"accuracy\": 85.5,
    \"duration\": 5000,
    \"passed\": true
}"
test_endpoint "Submit Accuracy Measurement" "POST" "/api/accuracy" "$ACCURACY_DATA"

# Test 7: Submit Gaze Point
GAZE_DATA="{
    \"session_id\": $SESSION_ID,
    \"x\": 500.2,
    \"y\": 300.8,
    \"panel\": \"A\",
    \"phase\": \"middle\"
}"
test_endpoint "Submit Gaze Point" "POST" "/api/gaze-point" "$GAZE_DATA"

# Test 8: Submit Reading Event
READING_EVENT_DATA="{
    \"session_id\": $SESSION_ID,
    \"event_type\": \"start\",
    \"panel\": \"A\",
    \"duration\": 0
}"
test_endpoint "Submit Reading Event" "POST" "/api/reading-event" "$READING_EVENT_DATA"

echo "=============================================="
echo -e "${GREEN}All endpoint tests completed!${NC}"
echo ""
echo "To verify data was saved, run: ./view-db.sh"

