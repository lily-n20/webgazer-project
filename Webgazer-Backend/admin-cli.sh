#!/bin/bash

# Simple terminal-based admin interface for Readability Study
# Usage: ./admin-cli.sh

API_URL="${API_URL:-http://localhost:8080}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Helper functions
show_message() {
    local color=$1
    local message=$2
    echo -e "${color}${message}${NC}"
}

show_success() {
    show_message "$GREEN" "✓ $1"
}

show_error() {
    show_message "$RED" "✗ $1"
}

show_info() {
    show_message "$BLUE" "ℹ $1"
}

# API functions
api_get() {
    local response=$(curl -s -w "\n%{http_code}" "$API_URL$1")
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | sed '$d')
    
    if [[ "$http_code" -ge 200 && "$http_code" -lt 300 ]]; then
        echo "$body"
        return 0
    else
        show_error "HTTP $http_code: $body"
        echo "$body"
        return 1
    fi
}

api_post() {
    local response=$(curl -s -w "\n%{http_code}" -X POST "$API_URL$1" \
        -H "Content-Type: application/json" \
        -d "$2")
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | sed '$d')
    
    if [[ "$http_code" -ge 200 && "$http_code" -lt 300 ]]; then
        echo "$body"
        return 0
    else
        show_error "HTTP $http_code: $body"
        echo "$body"
        return 1
    fi
}

api_put() {
    local response=$(curl -s -w "\n%{http_code}" -X PUT "$API_URL$1" \
        -H "Content-Type: application/json" \
        -d "$2")
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | sed '$d')
    
    if [[ "$http_code" -ge 200 && "$http_code" -lt 300 ]]; then
        echo "$body"
        return 0
    else
        show_error "HTTP $http_code: $body"
        echo "$body"
        return 1
    fi
}

api_delete() {
    local response=$(curl -s -w "\n%{http_code}" -X DELETE "$API_URL$1")
    local http_code=$(echo "$response" | tail -n1)
    local body=$(echo "$response" | sed '$d')
    
    if [[ "$http_code" -ge 200 && "$http_code" -lt 300 ]]; then
        echo "$body"
        return 0
    else
        show_error "HTTP $http_code: $body"
        echo "$body"
        return 1
    fi
}

# Check API connection
check_api_connection() {
    local health=$(curl -s "$API_URL/api/health" 2>&1)
    if echo "$health" | grep -q "ok"; then
        return 0
    else
        show_error "Cannot connect to API at $API_URL"
        show_error "Response: $health"
        show_info "Make sure the backend is running: cd Webgazer-Backend && go run ."
        return 1
    fi
}

# Study Text Functions
list_study_texts() {
    echo ""
    show_info "Loading study texts..."
    
    if ! check_api_connection; then
        return 1
    fi
    
    response=$(api_get "/api/admin/study-text")
    local exit_code=$?
    
    if [[ $exit_code -eq 0 ]] && echo "$response" | jq -e '.success' > /dev/null 2>&1; then
        echo "$response" | jq -r '.data[] | "ID: \(.id) | Version: \(.version) | Active: \(.active) | Content: \(.content[0:50])..."'
    else
        show_error "Failed to load study texts"
        if echo "$response" | grep -q "404"; then
            show_error "Admin endpoints not found. Make sure you've restarted the backend after adding admin routes."
            show_info "The backend needs to be restarted to load the admin endpoints."
        else
            echo "$response"
        fi
    fi
}

create_study_text() {
    echo ""
    read -p "Version (default: default): " version
    version=${version:-default}
    
    echo "Enter content (press Enter, then type content, end with Ctrl+D):"
    content=$(cat)
    
    read -p "Set as active? (y/n): " active_choice
    active=false
    if [[ "$active_choice" == "y" || "$active_choice" == "Y" ]]; then
        active=true
    fi
    
    data=$(jq -n \
        --arg version "$version" \
        --arg content "$content" \
        --argjson active "$active" \
        '{version: $version, content: $content, active: $active}')
    
    response=$(api_post "/api/admin/study-text" "$data")
    
    if echo "$response" | jq -e '.success' > /dev/null 2>&1; then
        show_success "Study text created!"
        echo "$response" | jq '.'
    else
        show_error "Failed to create study text"
        echo "$response"
    fi
}

update_study_text() {
    echo ""
    read -p "Enter Study Text ID to update: " id
    
    if [[ -z "$id" ]]; then
        show_error "ID is required"
        return
    fi
    
    read -p "New version (press Enter to skip): " version
    echo "New content (press Enter, then type content, end with Ctrl+D, or just Enter to skip):"
    content=$(cat)
    
    read -p "Set as active? (y/n/skip): " active_choice
    active_json="null"
    if [[ "$active_choice" == "y" || "$active_choice" == "Y" ]]; then
        active_json="true"
    elif [[ "$active_choice" == "n" || "$active_choice" == "N" ]]; then
        active_json="false"
    fi
    
    data=$(jq -n \
        --argjson id "$id" \
        --arg version "$version" \
        --arg content "$content" \
        --argjson active "$active_json" \
        '{id: $id} + 
         (if $version != "" then {version: $version} else {} end) +
         (if $content != "" then {content: $content} else {} end) +
         (if $active != null then {active: $active} else {} end)')
    
    response=$(api_put "/api/admin/study-text" "$data")
    
    if echo "$response" | jq -e '.success' > /dev/null 2>&1; then
        show_success "Study text updated!"
        echo "$response" | jq '.'
    else
        show_error "Failed to update study text"
        echo "$response"
    fi
}

# Quiz Question Functions
list_quiz_questions() {
    echo ""
    show_info "Loading quiz questions..."
    
    # Get study texts first
    study_texts=$(api_get "/api/admin/study-text")
    
    echo "$study_texts" | jq -r '.data[] | @json' | while read -r text_json; do
        text_id=$(echo "$text_json" | jq -r '.id')
        text_version=$(echo "$text_json" | jq -r '.version')
        
        questions=$(api_get "/api/quiz-questions?study_text_id=$text_id")
        
        if echo "$questions" | jq -e '. | length > 0' > /dev/null 2>&1; then
            echo ""
            show_info "Study Text: $text_version (ID: $text_id)"
            echo "$questions" | jq -r '.[] | "  ID: \(.id) | Q: \(.id) | Prompt: \(.prompt[0:40])... | Answer: \(.choices[.answer])"'
        fi
    done
}

create_quiz_question() {
    echo ""
    read -p "Study Text ID: " study_text_id
    read -p "Question ID (e.g., q1, q6): " question_id
    read -p "Prompt: " prompt
    
    echo "Enter choices (one per line, press Enter after each, empty line to finish):"
    choices=()
    while true; do
        read -p "Choice: " choice
        if [[ -z "$choice" ]]; then
            break
        fi
        choices+=("$choice")
    done
    
    if [[ ${#choices[@]} -eq 0 ]]; then
        show_error "At least one choice is required"
        return
    fi
    
    echo "Choices entered:"
    for i in "${!choices[@]}"; do
        echo "  [$i] ${choices[$i]}"
    done
    
    read -p "Correct answer index (0-$((${#choices[@]} - 1))): " answer
    read -p "Display order: " order
    
    choices_json=$(printf '%s\n' "${choices[@]}" | jq -R . | jq -s .)
    
    data=$(jq -n \
        --argjson study_text_id "$study_text_id" \
        --arg question_id "$question_id" \
        --arg prompt "$prompt" \
        --argjson choices "$choices_json" \
        --argjson answer "$answer" \
        --argjson order "${order:-0}" \
        '{study_text_id: $study_text_id, question_id: $question_id, prompt: $prompt, choices: $choices, answer: $answer, order: $order}')
    
    response=$(api_post "/api/admin/quiz-question" "$data")
    
    if echo "$response" | jq -e '.success' > /dev/null 2>&1; then
        show_success "Quiz question created!"
        echo "$response" | jq '.'
    else
        show_error "Failed to create quiz question"
        echo "$response"
    fi
}

update_quiz_question() {
    echo ""
    read -p "Enter Quiz Question ID to update: " id
    
    if [[ -z "$id" ]]; then
        show_error "ID is required"
        return
    fi
    
    # Get current question
    current=$(api_get "/api/admin/quiz-question?id=$id")
    
    if ! echo "$current" | jq -e '.success' > /dev/null 2>&1; then
        show_error "Question not found"
        return
    fi
    
    echo "Current question:"
    echo "$current" | jq '.data'
    
    read -p "New prompt (press Enter to skip): " prompt
    read -p "Update choices? (y/n): " update_choices
    
    choices_json="null"
    if [[ "$update_choices" == "y" || "$update_choices" == "Y" ]]; then
        echo "Enter choices (one per line, empty line to finish):"
        choices=()
        while true; do
            read -p "Choice: " choice
            if [[ -z "$choice" ]]; then
                break
            fi
            choices+=("$choice")
        done
        choices_json=$(printf '%s\n' "${choices[@]}" | jq -R . | jq -s .)
    fi
    
    read -p "New answer index (press Enter to skip): " answer
    read -p "New order (press Enter to skip): " order
    
    data=$(jq -n \
        --argjson id "$id" \
        --arg prompt "$prompt" \
        --argjson choices "$choices_json" \
        --arg answer "$answer" \
        --arg order "$order" \
        '{id: $id} +
         (if $prompt != "" then {prompt: $prompt} else {} end) +
         (if $choices != "null" then {choices: $choices} else {} end) +
         (if $answer != "" then {answer: ($answer | tonumber)} else {} end) +
         (if $order != "" then {order: ($order | tonumber)} else {} end)')
    
    response=$(api_put "/api/admin/quiz-question" "$data")
    
    if echo "$response" | jq -e '.success' > /dev/null 2>&1; then
        show_success "Quiz question updated!"
        echo "$response" | jq '.'
    else
        show_error "Failed to update quiz question"
        echo "$response"
    fi
}

delete_quiz_question() {
    echo ""
    read -p "Enter Quiz Question ID to delete: " id
    
    if [[ -z "$id" ]]; then
        show_error "ID is required"
        return
    fi
    
    read -p "Are you sure? (y/n): " confirm
    if [[ "$confirm" != "y" && "$confirm" != "Y" ]]; then
        show_info "Cancelled"
        return
    fi
    
    response=$(api_delete "/api/admin/quiz-question?id=$id")
    
    if echo "$response" | jq -e '.success' > /dev/null 2>&1; then
        show_success "Quiz question deleted!"
    else
        show_error "Failed to delete quiz question"
        echo "$response"
    fi
}

# Main menu
show_menu() {
    clear
    echo "=========================================="
    echo "  Readability Study Admin (Terminal)"
    echo "=========================================="
    echo ""
    echo "API URL: $API_URL"
    echo ""
    echo "Study Text Management:"
    echo "  1) List all study texts"
    echo "  2) Create new study text"
    echo "  3) Update study text"
    echo ""
    echo "Quiz Question Management:"
    echo "  4) List all quiz questions"
    echo "  5) Create new quiz question"
    echo "  6) Update quiz question"
    echo "  7) Delete quiz question"
    echo ""
    echo "  0) Exit"
    echo ""
    read -p "Select option: " choice
    
    case $choice in
        1) list_study_texts; read -p "Press Enter to continue..."; show_menu ;;
        2) create_study_text; read -p "Press Enter to continue..."; show_menu ;;
        3) update_study_text; read -p "Press Enter to continue..."; show_menu ;;
        4) list_quiz_questions; read -p "Press Enter to continue..."; show_menu ;;
        5) create_quiz_question; read -p "Press Enter to continue..."; show_menu ;;
        6) update_quiz_question; read -p "Press Enter to continue..."; show_menu ;;
        7) delete_quiz_question; read -p "Press Enter to continue..."; show_menu ;;
        0) show_info "Goodbye!"; exit 0 ;;
        *) show_error "Invalid option"; sleep 1; show_menu ;;
    esac
}

# Check dependencies
if ! command -v jq &> /dev/null; then
    show_error "jq is required but not installed."
    echo "Install with: brew install jq"
    exit 1
fi

if ! command -v curl &> /dev/null; then
    show_error "curl is required but not installed."
    exit 1
fi

# Check dependencies and connection before starting
if ! check_api_connection; then
    echo ""
    read -p "Continue anyway? (y/n): " continue_anyway
    if [[ "$continue_anyway" != "y" && "$continue_anyway" != "Y" ]]; then
        exit 1
    fi
fi

# Start
show_menu

