#!/bin/bash

# Quick script to view database contents
DB_FILE="readability.db"

if [ ! -f "$DB_FILE" ]; then
    echo "Database file not found: $DB_FILE"
    echo "Make sure you've run the backend at least once to create the database."
    exit 1
fi

echo "📊 Readability Study Database Viewer"
echo "======================================"
echo ""

# Check if sqlite3 is installed
if ! command -v sqlite3 &> /dev/null; then
    echo "❌ sqlite3 is not installed."
    echo "Install it with: brew install sqlite3 (macOS) or apt-get install sqlite3 (Linux)"
    exit 1
fi

echo "Tables in database:"
sqlite3 "$DB_FILE" ".tables"
echo ""

echo "📋 Participants:"
sqlite3 "$DB_FILE" "SELECT id, source, created_at FROM participants LIMIT 10;"
echo ""

echo "📋 Study Sessions (last 5):"
sqlite3 -header -column "$DB_FILE" "SELECT id, session_id, participant_id, font_preference, preferred_font_type, time_ams, time_bms, created_at FROM study_sessions ORDER BY created_at DESC LIMIT 5;"
echo ""

echo "📋 Quiz Responses (from study_sessions JSON field):"
sqlite3 -header -column "$DB_FILE" "SELECT id, session_id, quiz_responses_json FROM study_sessions WHERE quiz_responses_json IS NOT NULL AND quiz_responses_json != '' ORDER BY created_at DESC LIMIT 5;"
echo ""

echo "📋 Individual Quiz Responses Table (currently empty - not implemented yet):"
sqlite3 -header -column "$DB_FILE" "SELECT COUNT(*) as count FROM quiz_responses;"
echo ""

echo "📊 Statistics:"
echo "Total Participants: $(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM participants;")"
echo "Total Sessions: $(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM study_sessions;")"
echo "Sessions with Quiz Data (JSON): $(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM study_sessions WHERE quiz_responses_json IS NOT NULL AND quiz_responses_json != '';")"
echo "Individual Quiz Responses Table: $(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM quiz_responses;")"
echo "Total Calibration Data: $(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM calibration_data;")"
echo "Total Gaze Points: $(sqlite3 "$DB_FILE" "SELECT COUNT(*) FROM gaze_points;")"
echo ""

echo "💡 Note: Quiz responses are stored as JSON in study_sessions.quiz_responses_json"
echo "💡 Individual quiz responses should also be in the quiz_responses table"
echo "   (if the frontend successfully submitted them)"
echo ""
echo "💡 To explore interactively: sqlite3 $DB_FILE"
echo "💡 Or use DB Browser for SQLite: https://sqlitebrowser.org/"

