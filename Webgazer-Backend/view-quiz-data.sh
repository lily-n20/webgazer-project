#!/bin/bash

# View quiz responses from JSON field
DB_FILE="readability.db"

if [ ! -f "$DB_FILE" ]; then
    echo "Database file not found: $DB_FILE"
    exit 1
fi

echo "📋 Quiz Responses (from JSON in study_sessions)"
echo "================================================"
echo ""

sqlite3 -header -column "$DB_FILE" <<EOF
SELECT 
    id as session_id,
    session_id as session_uuid,
    quiz_responses_json
FROM study_sessions 
WHERE quiz_responses_json IS NOT NULL 
  AND quiz_responses_json != ''
ORDER BY created_at DESC;
EOF

echo ""
echo "💡 To see formatted JSON, use DB Browser for SQLite or run:"
echo "   sqlite3 readability.db \"SELECT quiz_responses_json FROM study_sessions WHERE id = 1;\" | python3 -m json.tool"

