package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func main() {
	// Initialize database
	var err error
	db, err = gorm.Open(sqlite.Open("readability.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto-migrate the schema
	err = db.AutoMigrate(
		&Participant{},
		&StudySession{},
		&CalibrationData{},
		&AccuracyMeasurement{},
		&QuizResponse{},
		&GazePoint{},
		&ReadingEvent{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database initialized successfully")

	// Setup routes
	http.HandleFunc("/api/participant", handleParticipant)
	http.HandleFunc("/api/session", handleSession)
	http.HandleFunc("/api/quiz-response", handleQuizResponse)
	http.HandleFunc("/api/calibration", handleCalibration)
	http.HandleFunc("/api/gaze-point", handleGazePoint)
	http.HandleFunc("/api/reading-event", handleReadingEvent)
	http.HandleFunc("/api/accuracy", handleAccuracy)
	http.HandleFunc("/api/health", handleHealth)

	// CORS middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:5173", "http://localhost:4173", "http://localhost:3000"}),
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(http.DefaultServeMux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleParticipant(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var participant Participant
	if err := json.NewDecoder(r.Body).Decode(&participant); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set default source if not provided
	if participant.Source == "" {
		participant.Source = "web"
	}

	// Create participant in database
	if err := db.Create(&participant).Error; err != nil {
		http.Error(w, "Failed to save participant: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      participant.ID,
		"source":  participant.Source,
	})
}

func handleSession(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var session StudySession
	if err := json.NewDecoder(r.Body).Decode(&session); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Create session in database
	if err := db.Create(&session).Error; err != nil {
		http.Error(w, "Failed to save session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"session_id": session.SessionID,
		"id": session.ID,
	})
}

func handleQuizResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var quizResponse QuizResponse
	if err := json.NewDecoder(r.Body).Decode(&quizResponse); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set timestamp if not provided
	if quizResponse.Timestamp.IsZero() {
		quizResponse.Timestamp = time.Now()
	}

	// Create quiz response in database
	if err := db.Create(&quizResponse).Error; err != nil {
		http.Error(w, "Failed to save quiz response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      quizResponse.ID,
	})
}

func handleCalibration(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var calibration CalibrationData
	if err := json.NewDecoder(r.Body).Decode(&calibration); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set timestamp if not provided
	if calibration.Timestamp.IsZero() {
		calibration.Timestamp = time.Now()
	}

	// Create calibration data in database
	if err := db.Create(&calibration).Error; err != nil {
		http.Error(w, "Failed to save calibration data: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      calibration.ID,
	})
}

func handleGazePoint(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var gazePoint GazePoint
	if err := json.NewDecoder(r.Body).Decode(&gazePoint); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set timestamp if not provided
	if gazePoint.Timestamp.IsZero() {
		gazePoint.Timestamp = time.Now()
	}

	// Create gaze point in database
	if err := db.Create(&gazePoint).Error; err != nil {
		http.Error(w, "Failed to save gaze point: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      gazePoint.ID,
	})
}

func handleReadingEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var readingEvent ReadingEvent
	if err := json.NewDecoder(r.Body).Decode(&readingEvent); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set timestamp if not provided
	if readingEvent.Timestamp.IsZero() {
		readingEvent.Timestamp = time.Now()
	}

	// Create reading event in database
	if err := db.Create(&readingEvent).Error; err != nil {
		http.Error(w, "Failed to save reading event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      readingEvent.ID,
	})
}

func handleAccuracy(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var accuracy AccuracyMeasurement
	if err := json.NewDecoder(r.Body).Decode(&accuracy); err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Set timestamp if not provided
	if accuracy.Timestamp.IsZero() {
		accuracy.Timestamp = time.Now()
	}

	// Create accuracy measurement in database
	if err := db.Create(&accuracy).Error; err != nil {
		http.Error(w, "Failed to save accuracy measurement: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"id":      accuracy.ID,
	})
}

