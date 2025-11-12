//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "http://localhost:8080"

func main() {
	fmt.Println("🧪 Testing Readability Backend API")
	fmt.Println("==================================")
	fmt.Println()

	// Test 1: Health Check
	fmt.Println("1. Testing Health Endpoint...")
	resp, err := http.Get(baseURL + "/api/health")
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n\n", string(body))

	// Test 2: Create a Study Session
	fmt.Println("2. Creating a Study Session...")
	sessionData := map[string]interface{}{
		"participant_id":     1,
		"calibration_points": 25,
		"font_left":          "serif",
		"font_right":         "sans",
		"time_left_ms":       5000,
		"time_right_ms":      4500,
		"time_a_ms":          5000,
		"time_b_ms":          4500,
		"font_preference":    "A",
		"preferred_font_type": "serif",
		"user_agent":         "test-agent",
		"screen_width":        1920,
		"screen_height":       1080,
	}

	jsonData, _ := json.Marshal(sessionData)
	resp, err = http.Post(baseURL+"/api/session", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}
	defer resp.Body.Close()
	body, _ = io.ReadAll(resp.Body)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Response: %s\n\n", string(body))

	// Parse session ID from response
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	sessionID, ok := result["id"].(float64)
	if !ok {
		fmt.Println("⚠️  Could not parse session ID from response")
		return
	}

	fmt.Printf("✅ Session created with ID: %.0f\n", sessionID)
	fmt.Println("\n✅ All basic tests passed!")
	fmt.Println("\nTo test with curl, use:")
	fmt.Println("  curl http://localhost:8080/api/health")
	fmt.Println("  curl -X POST http://localhost:8080/api/session -H 'Content-Type: application/json' -d '{...}'")
}

