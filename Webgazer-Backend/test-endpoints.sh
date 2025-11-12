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

# Test 2: Fetch Study Text
test_endpoint "Fetch Study Text" "GET" "/api/study-text" ""
test_endpoint "Fetch Study Text (with version)" "GET" "/api/study-text?version=default" ""

# Test 2a: Admin - List Study Texts
test_endpoint "Admin: List Study Texts" "GET" "/api/admin/study-text" ""

# Test 2b: Admin - Create Study Text
ADMIN_STUDY_TEXT_DATA='{"version":"test","content":"Test passage for endpoint testing","font_left":"serif","font_right":"sans","active":false}'
test_endpoint "Admin: Create Study Text" "POST" "/api/admin/study-text" "$ADMIN_STUDY_TEXT_DATA"

# Test 3: Fetch Quiz Questions
test_endpoint "Fetch Quiz Questions" "GET" "/api/quiz-questions" ""

# Test 4: Create Participant
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

# Test 5: Create Study Session
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

# Test 6: Submit Quiz Response
QUIZ_RESPONSE_DATA="{
    \"session_id\": $SESSION_ID,
    \"question_id\": \"q1\",
    \"answer_index\": 1,
    \"is_correct\": true,
    \"response_time\": 3000
}"
test_endpoint "Submit Quiz Response" "POST" "/api/quiz-response" "$QUIZ_RESPONSE_DATA"

# Test 7: Submit Calibration Data
CALIBRATION_DATA="{
    \"session_id\": $SESSION_ID,
    \"point_index\": 0,
    \"click_number\": 1,
    \"x\": 100.5,
    \"y\": 200.3
}"
test_endpoint "Submit Calibration Data" "POST" "/api/calibration" "$CALIBRATION_DATA"

# Test 8: Submit Accuracy Measurement
ACCURACY_DATA="{
    \"session_id\": $SESSION_ID,
    \"accuracy\": 85.5,
    \"duration\": 5000,
    \"passed\": true
}"
test_endpoint "Submit Accuracy Measurement" "POST" "/api/accuracy" "$ACCURACY_DATA"

# Test 9: Submit Gaze Point
GAZE_DATA="{
    \"session_id\": $SESSION_ID,
    \"x\": 500.2,
    \"y\": 300.8,
    \"panel\": \"A\",
    \"phase\": \"middle\"
}"
test_endpoint "Submit Gaze Point" "POST" "/api/gaze-point" "$GAZE_DATA"

# Test 10: Submit Reading Event
READING_EVENT_DATA="{
    \"session_id\": $SESSION_ID,
    \"event_type\": \"start\",
    \"panel\": \"A\",
    \"duration\": 0
}"
test_endpoint "Submit Reading Event" "POST" "/api/reading-event" "$READING_EVENT_DATA"

# Test 11: Admin - Get Quiz Question
test_endpoint "Admin: Get Quiz Question" "GET" "/api/admin/quiz-question?id=1" ""

# Test 12: Admin - Create Quiz Question
ADMIN_QUIZ_DATA='{
    "study_text_id": 1,
    "question_id": "q_test",
    "prompt": "Test question?",
    "choices": ["A", "B", "C", "D"],
    "answer": 0,
    "order": 99
}'
test_endpoint "Admin: Create Quiz Question" "POST" "/api/admin/quiz-question" "$ADMIN_QUIZ_DATA"

# Test 13: Admin - Create Passage
# First, get a study text ID (use the default one, which should be ID 1)
STUDY_TEXT_ID=1
ADMIN_PASSAGE_DATA="{
    \"study_text_id\": $STUDY_TEXT_ID,
    \"order\": 0,
    \"title\": \"Test Passage\",
    \"content\": \"This is a test passage for endpoint testing. It contains sample text to verify the passage creation endpoint works correctly.\",
    \"font_left\": \"serif\",
    \"font_right\": \"sans\"
}"
PASSAGE_RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "$ADMIN_PASSAGE_DATA" \
    "${BASE_URL}/api/admin/passage")
PASSAGE_ID=$(echo "$PASSAGE_RESPONSE" | jq -r '.id' 2>/dev/null)
test_endpoint "Admin: Create Passage" "POST" "/api/admin/passage" "$ADMIN_PASSAGE_DATA"

if [ -z "$PASSAGE_ID" ] || [ "$PASSAGE_ID" = "null" ]; then
    echo -e "${YELLOW}⚠ Warning: Could not get passage ID. Using ID 1 for remaining tests.${NC}"
    PASSAGE_ID=1
fi
echo "  Using Passage ID: $PASSAGE_ID"
echo ""

# Test 14: Admin - Get Passages by Study Text ID
test_endpoint "Admin: Get Passages by Study Text ID" "GET" "/api/admin/passage?study_text_id=$STUDY_TEXT_ID" ""

# Test 15: Admin - Get Single Passage by ID
test_endpoint "Admin: Get Single Passage" "GET" "/api/admin/passage?id=$PASSAGE_ID" ""

# Test 16: Admin - Update Passage
ADMIN_PASSAGE_UPDATE_DATA="{
    \"id\": $PASSAGE_ID,
    \"title\": \"Updated Test Passage\",
    \"content\": \"This is an updated test passage. The content has been modified to test the update endpoint.\",
    \"font_left\": \"sans\",
    \"font_right\": \"serif\"
}"
test_endpoint "Admin: Update Passage" "PUT" "/api/admin/passage" "$ADMIN_PASSAGE_UPDATE_DATA"

# Test 17: Admin - Delete Passage (create a new one first to delete)
ADMIN_PASSAGE_DELETE_DATA="{
    \"study_text_id\": $STUDY_TEXT_ID,
    \"order\": 999,
    \"content\": \"This passage will be deleted.\",
    \"font_left\": \"serif\",
    \"font_right\": \"sans\"
}"
DELETE_PASSAGE_RESPONSE=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "$ADMIN_PASSAGE_DELETE_DATA" \
    "${BASE_URL}/api/admin/passage")
DELETE_PASSAGE_ID=$(echo "$DELETE_PASSAGE_RESPONSE" | jq -r '.id' 2>/dev/null)
if [ -n "$DELETE_PASSAGE_ID" ] && [ "$DELETE_PASSAGE_ID" != "null" ]; then
    test_endpoint "Admin: Delete Passage" "DELETE" "/api/admin/passage?id=$DELETE_PASSAGE_ID" ""
else
    echo -e "${YELLOW}⚠ Warning: Could not create passage for deletion test. Skipping delete test.${NC}"
    echo ""
fi

echo "=============================================="
echo -e "${GREEN}All endpoint tests completed!${NC}"
echo ""
echo "To verify data was saved, run: ./view-db.sh"

