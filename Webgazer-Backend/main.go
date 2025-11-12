package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
		&StudyText{},
		&Passage{},
		&QuizQuestion{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database initialized successfully")

	// Setup Gin router
	router := gin.Default()

	// Configure trusted proxies (security best practice)
	// For local development, we don't trust any proxies
	// In production, configure this based on your infrastructure
	router.SetTrustedProxies(nil)

	// Configure CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:4173", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type"}
	router.Use(cors.New(config))

	// API routes
	api := router.Group("/api")
	{
		api.POST("/participant", handleParticipant)
		api.POST("/session", handleSession)
		api.POST("/quiz-response", handleQuizResponse)
		api.POST("/calibration", handleCalibration)
		api.POST("/gaze-point", handleGazePoint)
		api.POST("/reading-event", handleReadingEvent)
		api.POST("/accuracy", handleAccuracy)
		api.GET("/study-text", handleStudyText)
		api.GET("/quiz-questions", handleQuizQuestions)
		api.GET("/health", handleHealth)

		// Admin routes
		admin := api.Group("/admin")
		{
			admin.POST("/study-text", handleAdminStudyText)
			admin.PUT("/study-text", handleAdminStudyText)
			admin.GET("/study-text", handleAdminStudyText)
			admin.POST("/passage", handleAdminPassage)
			admin.PUT("/passage", handleAdminPassage)
			admin.DELETE("/passage", handleAdminPassage)
			admin.GET("/passage", handleAdminPassage)
			admin.POST("/quiz-question", handleAdminQuizQuestion)
			admin.PUT("/quiz-question", handleAdminQuizQuestion)
			admin.DELETE("/quiz-question", handleAdminQuizQuestion)
			admin.GET("/quiz-question", handleAdminQuizQuestion)
		}
	}

	// Seed initial data if database is empty
	seedInitialData()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s\n", port)
	log.Fatal(router.Run(":" + port))
}

func handleHealth(c *gin.Context) {
	c.JSON(200, gin.H{"status": "ok"})
}

func handleParticipant(c *gin.Context) {
	var participant Participant
	if err := c.ShouldBindJSON(&participant); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// Set default source if not provided
	if participant.Source == "" {
		participant.Source = "web"
	}

	// Create participant in database
	if err := db.Create(&participant).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save participant: " + err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"id":      participant.ID,
		"source":  participant.Source,
	})
}

func handleSession(c *gin.Context) {
	var session StudySession
	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// Create session in database
	if err := db.Create(&session).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save session: " + err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"success":   true,
		"session_id": session.SessionID,
		"id":        session.ID,
	})
}

func handleQuizResponse(c *gin.Context) {
	var quizResponse QuizResponse
	if err := c.ShouldBindJSON(&quizResponse); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// Set timestamp if not provided
	if quizResponse.Timestamp.IsZero() {
		quizResponse.Timestamp = time.Now()
	}

	// Create quiz response in database
	if err := db.Create(&quizResponse).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save quiz response: " + err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"id":      quizResponse.ID,
	})
}

func handleCalibration(c *gin.Context) {
	var calibration CalibrationData
	if err := c.ShouldBindJSON(&calibration); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// Set timestamp if not provided
	if calibration.Timestamp.IsZero() {
		calibration.Timestamp = time.Now()
	}

	// Create calibration data in database
	if err := db.Create(&calibration).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save calibration data: " + err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"id":      calibration.ID,
	})
}

func handleGazePoint(c *gin.Context) {
	var gazePoint GazePoint
	if err := c.ShouldBindJSON(&gazePoint); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// Set timestamp if not provided
	if gazePoint.Timestamp.IsZero() {
		gazePoint.Timestamp = time.Now()
	}

	// Create gaze point in database
	if err := db.Create(&gazePoint).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save gaze point: " + err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"id":      gazePoint.ID,
	})
}

func handleReadingEvent(c *gin.Context) {
	var readingEvent ReadingEvent
	if err := c.ShouldBindJSON(&readingEvent); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// Set timestamp if not provided
	if readingEvent.Timestamp.IsZero() {
		readingEvent.Timestamp = time.Now()
	}

	// Create reading event in database
	if err := db.Create(&readingEvent).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save reading event: " + err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"id":      readingEvent.ID,
	})
}

func handleAccuracy(c *gin.Context) {
	var accuracy AccuracyMeasurement
	if err := c.ShouldBindJSON(&accuracy); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
		return
	}

	// Set timestamp if not provided
	if accuracy.Timestamp.IsZero() {
		accuracy.Timestamp = time.Now()
	}

	// Create accuracy measurement in database
	if err := db.Create(&accuracy).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to save accuracy measurement: " + err.Error()})
		return
	}

	c.JSON(201, gin.H{
		"success": true,
		"id":      accuracy.ID,
	})
}

func handleStudyText(c *gin.Context) {
	// Get version from query parameter, default to "default"
	version := c.DefaultQuery("version", "default")

	var studyText StudyText
	if err := db.Preload("Passages").Where("version = ? AND active = ?", version, true).First(&studyText).Error; err != nil {
		// If not found, try to get any active study text
		if err := db.Preload("Passages").Where("active = ?", true).First(&studyText).Error; err != nil {
			c.JSON(404, gin.H{"error": "No study text found"})
			return
		}
	}

	// Build response - include passages if they exist, otherwise use legacy content
	response := gin.H{
		"id":        studyText.ID,
		"version":   studyText.Version,
		"font_left": studyText.FontLeft,
		"font_right": studyText.FontRight,
	}

	// If passages exist, return them; otherwise return legacy content for backward compatibility
	if len(studyText.Passages) > 0 {
		response["passages"] = studyText.Passages
	} else {
		response["content"] = studyText.Content
	}

	c.JSON(200, response)
}

func handleQuizQuestions(c *gin.Context) {
	// Get study_text_id from query parameter
	studyTextID := c.Query("study_text_id")

	var questions []QuizQuestion
	query := db.Order("`order` ASC")

	if studyTextID != "" {
		query = query.Where("study_text_id = ?", studyTextID)
	} else {
		// If no study_text_id provided, get questions for active study text
		var studyText StudyText
		if err := db.Where("active = ?", true).First(&studyText).Error; err != nil {
			c.JSON(404, gin.H{"error": "No active study text found"})
			return
		}
		query = query.Where("study_text_id = ?", studyText.ID)
	}

	if err := query.Find(&questions).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch quiz questions: " + err.Error()})
		return
	}

	// Format response to match frontend expectations
	type QuizQResponse struct {
		ID      string   `json:"id"`
		Prompt  string   `json:"prompt"`
		Choices []string `json:"choices"`
		Answer  int      `json:"answer"`
	}

	response := make([]QuizQResponse, len(questions))
	for i, q := range questions {
		var choices []string
		if err := json.Unmarshal([]byte(q.Choices), &choices); err != nil {
			log.Printf("Error unmarshaling choices for question %s: %v", q.QuestionID, err)
			continue
		}

		response[i] = QuizQResponse{
			ID:      q.QuestionID,
			Prompt:  q.Prompt,
			Choices: choices,
			Answer:  q.Answer,
		}
	}

	c.JSON(200, response)
}

// Admin endpoints for managing study text, passages, and quiz questions

func handleAdminPassage(c *gin.Context) {
	switch c.Request.Method {
	case "POST":
		// Create new passage
		var passage Passage
		if err := c.ShouldBindJSON(&passage); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		// Validate required fields
		if passage.StudyTextID == 0 || passage.Content == "" {
			c.JSON(400, gin.H{"error": "study_text_id and content are required"})
			return
		}

		// Verify study text exists
		var studyText StudyText
		if err := db.First(&studyText, passage.StudyTextID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Study text not found"})
			return
		}

		// If order not specified, set it to the next available order
		if passage.Order == 0 {
			var maxOrder int
			db.Model(&Passage{}).Where("study_text_id = ?", passage.StudyTextID).Select("COALESCE(MAX(`order`), -1)").Scan(&maxOrder)
			passage.Order = maxOrder + 1
		}

		if err := db.Create(&passage).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create passage: " + err.Error()})
			return
		}

		c.JSON(201, gin.H{
			"success": true,
			"id":      passage.ID,
			"message": "Passage created successfully",
		})

	case "PUT":
		// Update existing passage
		var updateData struct {
			ID        uint   `json:"id"`
			Order     *int   `json:"order,omitempty"`
			Content   string `json:"content,omitempty"`
			Title     string `json:"title,omitempty"`
			FontLeft  string `json:"font_left,omitempty"`
			FontRight string `json:"font_right,omitempty"`
		}

		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		if updateData.ID == 0 {
			c.JSON(400, gin.H{"error": "ID is required"})
			return
		}

		var passage Passage
		if err := db.First(&passage, updateData.ID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Passage not found"})
			return
		}

		// Update fields
		if updateData.Content != "" {
			passage.Content = updateData.Content
		}
		if updateData.Title != "" {
			passage.Title = updateData.Title
		}
		if updateData.Order != nil {
			passage.Order = *updateData.Order
		}
		if updateData.FontLeft != "" {
			passage.FontLeft = updateData.FontLeft
		}
		if updateData.FontRight != "" {
			passage.FontRight = updateData.FontRight
		}

		if err := db.Save(&passage).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update passage: " + err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"id":      passage.ID,
			"message": "Passage updated successfully",
		})

	case "DELETE":
		// Delete passage
		id := c.Query("id")
		if id == "" {
			c.JSON(400, gin.H{"error": "ID parameter is required"})
			return
		}

		if err := db.Delete(&Passage{}, id).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete passage: " + err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "Passage deleted successfully",
		})

	case "GET":
		// Get passages - either by study_text_id or by id
		studyTextID := c.Query("study_text_id")
		id := c.Query("id")

		if id != "" {
			// Get single passage by ID
			var passage Passage
			if err := db.First(&passage, id).Error; err != nil {
				c.JSON(404, gin.H{"error": "Passage not found"})
				return
			}

			c.JSON(200, gin.H{
				"success": true,
				"data":    passage,
			})
		} else if studyTextID != "" {
			// Get all passages for a study text
			var passages []Passage
			if err := db.Where("study_text_id = ?", studyTextID).Order("`order` ASC").Find(&passages).Error; err != nil {
				c.JSON(500, gin.H{"error": "Failed to fetch passages: " + err.Error()})
				return
			}

			c.JSON(200, gin.H{
				"success": true,
				"data":    passages,
			})
		} else {
			c.JSON(400, gin.H{"error": "Either id or study_text_id parameter is required"})
		}

	default:
		c.JSON(405, gin.H{"error": "Method not allowed"})
	}
}

func handleAdminStudyText(c *gin.Context) {
	switch c.Request.Method {
	case "POST":
		// Create new study text
		var studyText StudyText
		if err := c.ShouldBindJSON(&studyText); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		// Set defaults
		if studyText.Version == "" {
			studyText.Version = "default"
		}
		if studyText.FontLeft == "" {
			studyText.FontLeft = "serif"
		}
		if studyText.FontRight == "" {
			studyText.FontRight = "sans"
		}

		// Check if version already exists (idempotent behavior)
		var existingStudyText StudyText
		if err := db.Where("version = ?", studyText.Version).First(&existingStudyText).Error; err == nil {
			// Version exists, return existing study text
			c.JSON(200, gin.H{
				"success": true,
				"id":      existingStudyText.ID,
				"message": "Study text with this version already exists",
			})
			return
		}

		// If this is set to active, deactivate all others
		if studyText.Active {
			db.Model(&StudyText{}).Where("active = ?", true).Update("active", false)
		}

		if err := db.Create(&studyText).Error; err != nil {
			// Check for unique constraint violation (fallback check)
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				c.JSON(409, gin.H{
					"error": fmt.Sprintf("Study text with version '%s' already exists", studyText.Version),
				})
				return
			}
			c.JSON(500, gin.H{"error": "Failed to create study text: " + err.Error()})
			return
		}

		c.JSON(201, gin.H{
			"success": true,
			"id":      studyText.ID,
			"message": "Study text created successfully",
		})

	case "PUT":
		// Update existing study text
		var updateData struct {
			ID        uint   `json:"id"`
			Version   string `json:"version,omitempty"`
			Content   string `json:"content,omitempty"`
			FontLeft  string `json:"font_left,omitempty"`
			FontRight string `json:"font_right,omitempty"`
			Active    *bool  `json:"active,omitempty"`
		}

		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		if updateData.ID == 0 {
			c.JSON(400, gin.H{"error": "ID is required"})
			return
		}

		var studyText StudyText
		if err := db.First(&studyText, updateData.ID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Study text not found"})
			return
		}

		// Update fields
		if updateData.Version != "" {
			studyText.Version = updateData.Version
		}
		if updateData.Content != "" {
			studyText.Content = updateData.Content
		}
		if updateData.FontLeft != "" {
			studyText.FontLeft = updateData.FontLeft
		}
		if updateData.FontRight != "" {
			studyText.FontRight = updateData.FontRight
		}
		if updateData.Active != nil {
			// If setting to active, deactivate all others first
			if *updateData.Active {
				db.Model(&StudyText{}).Where("active = ? AND id != ?", true, updateData.ID).Update("active", false)
			}
			studyText.Active = *updateData.Active
		}

		if err := db.Save(&studyText).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update study text: " + err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"id":      studyText.ID,
			"message": "Study text updated successfully",
		})

	case "GET":
		// List all study texts
		var studyTexts []StudyText
		if err := db.Order("created_at DESC").Find(&studyTexts).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch study texts: " + err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"data":    studyTexts,
		})

	default:
		c.JSON(405, gin.H{"error": "Method not allowed"})
	}
}

func handleAdminQuizQuestion(c *gin.Context) {
	switch c.Request.Method {
	case "POST":
		// Create new quiz question
		var questionData struct {
			StudyTextID uint     `json:"study_text_id"`
			QuestionID  string   `json:"question_id"`
			Prompt      string   `json:"prompt"`
			Choices     []string `json:"choices"`
			Answer      int      `json:"answer"`
			Order       int      `json:"order"`
		}

		if err := c.ShouldBindJSON(&questionData); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		// Validate required fields
		if questionData.StudyTextID == 0 || questionData.QuestionID == "" || questionData.Prompt == "" {
			c.JSON(400, gin.H{"error": "study_text_id, question_id, and prompt are required"})
			return
		}

		// Convert choices to JSON string
		choicesJSON, err := json.Marshal(questionData.Choices)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid choices format: " + err.Error()})
			return
		}

		question := QuizQuestion{
			StudyTextID: questionData.StudyTextID,
			QuestionID:  questionData.QuestionID,
			Prompt:      questionData.Prompt,
			Choices:     string(choicesJSON),
			Answer:      questionData.Answer,
			Order:       questionData.Order,
		}

		if err := db.Create(&question).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to create quiz question: " + err.Error()})
			return
		}

		c.JSON(201, gin.H{
			"success": true,
			"id":      question.ID,
			"message": "Quiz question created successfully",
		})

	case "PUT":
		// Update existing quiz question
		var updateData struct {
			ID         uint      `json:"id"`
			QuestionID string    `json:"question_id,omitempty"`
			Prompt     string    `json:"prompt,omitempty"`
			Choices    []string  `json:"choices,omitempty"`
			Answer     *int      `json:"answer,omitempty"`
			Order      *int      `json:"order,omitempty"`
		}

		if err := c.ShouldBindJSON(&updateData); err != nil {
			c.JSON(400, gin.H{"error": "Invalid JSON: " + err.Error()})
			return
		}

		if updateData.ID == 0 {
			c.JSON(400, gin.H{"error": "ID is required"})
			return
		}

		var question QuizQuestion
		if err := db.First(&question, updateData.ID).Error; err != nil {
			c.JSON(404, gin.H{"error": "Quiz question not found"})
			return
		}

		// Update fields
		if updateData.QuestionID != "" {
			question.QuestionID = updateData.QuestionID
		}
		if updateData.Prompt != "" {
			question.Prompt = updateData.Prompt
		}
		if updateData.Choices != nil {
			choicesJSON, err := json.Marshal(updateData.Choices)
			if err != nil {
				c.JSON(400, gin.H{"error": "Invalid choices format: " + err.Error()})
				return
			}
			question.Choices = string(choicesJSON)
		}
		if updateData.Answer != nil {
			question.Answer = *updateData.Answer
		}
		if updateData.Order != nil {
			question.Order = *updateData.Order
		}

		if err := db.Save(&question).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to update quiz question: " + err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"id":      question.ID,
			"message": "Quiz question updated successfully",
		})

	case "DELETE":
		// Delete quiz question
		id := c.Query("id")
		if id == "" {
			c.JSON(400, gin.H{"error": "ID parameter is required"})
			return
		}

		if err := db.Delete(&QuizQuestion{}, id).Error; err != nil {
			c.JSON(500, gin.H{"error": "Failed to delete quiz question: " + err.Error()})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
			"message": "Quiz question deleted successfully",
		})

	case "GET":
		// Get single quiz question by ID
		id := c.Query("id")
		if id == "" {
			c.JSON(400, gin.H{"error": "ID parameter is required"})
			return
		}

		var question QuizQuestion
		if err := db.First(&question, id).Error; err != nil {
			c.JSON(404, gin.H{"error": "Quiz question not found"})
			return
		}

		// Parse choices JSON
		var choices []string
		json.Unmarshal([]byte(question.Choices), &choices)

		c.JSON(200, gin.H{
			"success": true,
			"data": gin.H{
				"id":           question.ID,
				"study_text_id": question.StudyTextID,
				"question_id":  question.QuestionID,
				"prompt":       question.Prompt,
				"choices":      choices,
				"answer":       question.Answer,
				"order":        question.Order,
			},
		})

	default:
		c.JSON(405, gin.H{"error": "Method not allowed"})
	}
}

