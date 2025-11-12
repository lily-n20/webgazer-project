#!/bin/bash

# Test a single endpoint - useful for debugging
# Usage: ./test-single-endpoint.sh <endpoint> <method> [data]
# Example: ./test-single-endpoint.sh /api/health GET
# Example: ./test-single-endpoint.sh /api/participant POST '{"source":"test"}'

BASE_URL="http://localhost:8080"
ENDPOINT=$1
METHOD=${2:-GET}
DATA=${3:-""}

if [ -z "$ENDPOINT" ]; then
    echo "Usage: $0 <endpoint> [method] [data]"
    echo "Example: $0 /api/health GET"
    echo "Example: $0 /api/participant POST '{\"source\":\"test\"}'"
    exit 1
fi

echo "Testing: ${METHOD} ${ENDPOINT}"
echo "URL: ${BASE_URL}${ENDPOINT}"
echo ""

if [ "$METHOD" = "GET" ]; then
    curl -v -X GET "${BASE_URL}${ENDPOINT}"
else
    echo "Data: ${DATA}"
    echo ""
    curl -v -X "${METHOD}" \
        -H "Content-Type: application/json" \
        -d "${DATA}" \
        "${BASE_URL}${ENDPOINT}"
fi

echo ""
echo ""
echo "💡 Tip: Use -v flag to see full request/response details"

