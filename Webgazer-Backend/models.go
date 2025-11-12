package main

import (
	"time"

	"gorm.io/gorm"
)

// Participant represents a study participant
type Participant struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Source    string    `gorm:"index" json:"source"` // e.g., "mturk", "prolific", "internal", etc.
	CreatedAt time.Time `json:"created_at"`
	
	// Relationships
	StudySessions []StudySession `gorm:"foreignKey:ParticipantID;references:ID" json:"study_sessions,omitempty"`
}

// StudySession represents a complete study session
type StudySession struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	SessionID         string    `gorm:"uniqueIndex;not null" json:"session_id"`
	ParticipantID     uint      `gorm:"index" json:"participant_id"`
	CreatedAt         time.Time `json:"created_at"`
	
	// Relationships
	Participant        Participant        `gorm:"foreignKey:ParticipantID;references:ID" json:"participant,omitempty"`
	CalibrationData    []CalibrationData  `gorm:"foreignKey:SessionID;references:ID" json:"calibration_data,omitempty"`
	AccuracyMeasurements []AccuracyMeasurement `gorm:"foreignKey:SessionID;references:ID" json:"accuracy_measurements,omitempty"`
	QuizResponses      []QuizResponse     `gorm:"foreignKey:SessionID;references:ID" json:"quiz_responses,omitempty"`
	GazePoints         []GazePoint        `gorm:"foreignKey:SessionID;references:ID" json:"gaze_points,omitempty"`
	ReadingEvents      []ReadingEvent     `gorm:"foreignKey:SessionID;references:ID" json:"reading_events,omitempty"`
	
	// Calibration data (legacy - kept for backward compatibility)
	CalibrationPoints int `json:"calibration_points"`
	
	// Reading session data
	FontLeft          string  `json:"font_left"`           // "serif" or "sans"
	FontRight         string  `json:"font_right"`          // "serif" or "sans"
	TimeLeftMS        int     `json:"time_left_ms"`        // reading time for left side
	TimeRightMS       int     `json:"time_right_ms"`       // reading time for right side
	TimeAMS           int     `json:"time_a_ms"`           // reading time for box A
	TimeBMS           int     `json:"time_b_ms"`           // reading time for box B
	FontPreference    string  `json:"font_preference"`     // "A" or "B"
	PreferredFontType string  `json:"preferred_font_type"` // "serif" or "sans"
	
	// Quiz responses (legacy - kept for backward compatibility)
	QuizResponsesJSON string  `json:"quiz_responses_json"` // JSON array of {question_id, answer_index}
	
	// Additional metadata
	UserAgent         string  `json:"user_agent,omitempty"`
	ScreenWidth       int     `json:"screen_width,omitempty"`
	ScreenHeight      int     `json:"screen_height,omitempty"`
}

// BeforeCreate hook to generate session ID if not provided
func (s *StudySession) BeforeCreate(tx *gorm.DB) error {
	if s.SessionID == "" {
		s.SessionID = generateSessionID()
	}
	return nil
}

// CalibrationData represents individual calibration point clicks
type CalibrationData struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID uint      `gorm:"index;not null" json:"session_id"`
	PointIndex int      `gorm:"not null" json:"point_index"` // Which calibration point (0-based)
	ClickNumber int     `gorm:"not null" json:"click_number"` // Which click on this point (1-5)
	X          float64  `gorm:"not null" json:"x"`           // X coordinate of calibration point
	Y          float64  `gorm:"not null" json:"y"`           // Y coordinate of calibration point
	Timestamp  time.Time `gorm:"not null" json:"timestamp"`
	
	// Relationship
	Session StudySession `gorm:"foreignKey:SessionID;references:ID" json:"session,omitempty"`
}

// AccuracyMeasurement represents accuracy check results
type AccuracyMeasurement struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID uint      `gorm:"index;not null" json:"session_id"`
	Accuracy  float64   `gorm:"not null" json:"accuracy"`    // Accuracy percentage
	Duration  int       `gorm:"not null" json:"duration"`    // Measurement duration in milliseconds
	Passed    bool      `gorm:"not null" json:"passed"`      // Whether it passed the threshold
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
	
	// Relationship
	Session StudySession `gorm:"foreignKey:SessionID;references:ID" json:"session,omitempty"`
}

// QuizResponse represents an individual quiz answer
type QuizResponse struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SessionID   uint      `gorm:"index;not null" json:"session_id"`
	QuestionID  string    `gorm:"not null" json:"question_id"`  // e.g., "q1", "q2"
	AnswerIndex int       `gorm:"not null" json:"answer_index"`  // Selected answer index (0-based)
	IsCorrect   *bool     `json:"is_correct,omitempty"`          // Whether answer is correct (nullable)
	ResponseTime int      `json:"response_time,omitempty"`       // Time to answer in milliseconds (optional)
	Timestamp   time.Time `gorm:"not null" json:"timestamp"`
	
	// Relationship
	Session StudySession `gorm:"foreignKey:SessionID;references:ID" json:"session,omitempty"`
}

// GazePoint represents a single gaze tracking data point
type GazePoint struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID uint      `gorm:"index;not null" json:"session_id"`
	X         float64   `gorm:"not null" json:"x"`              // X coordinate
	Y         float64   `gorm:"not null" json:"y"`              // Y coordinate
	Panel     string    `json:"panel,omitempty"`                 // "A", "B", "left", "right", or empty
	Phase     string    `json:"phase,omitempty"`                 // "start", "middle", "end", or empty
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
	
	// Relationship
	Session StudySession `gorm:"foreignKey:SessionID;references:ID" json:"session,omitempty"`
}

// ReadingEvent represents reading session milestones
type ReadingEvent struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID uint      `gorm:"index;not null" json:"session_id"`
	EventType string    `gorm:"not null" json:"event_type"`     // "start", "pause", "resume", "complete"
	Panel     string    `gorm:"not null" json:"panel"`            // "A", "B", "left", "right"
	Duration  int       `json:"duration,omitempty"`               // Duration in milliseconds (for complete events)
	Timestamp time.Time `gorm:"not null" json:"timestamp"`
	
	// Relationship
	Session StudySession `gorm:"foreignKey:SessionID;references:ID" json:"session,omitempty"`
}

// StudyText represents a reading passage for the study
type StudyText struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Version   string    `gorm:"uniqueIndex;not null" json:"version"` // e.g., "v1", "default"
	Content   string    `gorm:"type:text" json:"content,omitempty"`  // Legacy: single passage (deprecated, use Passages instead)
	FontLeft  string    `gorm:"default:serif" json:"font_left"`      // Font for left panel: "serif" or "sans"
	FontRight string    `gorm:"default:sans" json:"font_right"`      // Font for right panel: "serif" or "sans"
	Active    bool      `gorm:"default:true" json:"active"`          // Whether this is the active version
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// Relationships
	QuizQuestions []QuizQuestion `gorm:"foreignKey:StudyTextID;references:ID" json:"quiz_questions,omitempty"`
	Passages      []Passage      `gorm:"foreignKey:StudyTextID;references:ID;order:order ASC" json:"passages,omitempty"`
}

// Passage represents a single reading passage within a study text
type Passage struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	StudyTextID uint     `gorm:"index;not null" json:"study_text_id"`
	Order      int       `gorm:"not null" json:"order"`              // Display order (0, 1, 2, ...)
	Content    string    `gorm:"type:text;not null" json:"content"`   // The passage text
	Title      string    `json:"title,omitempty"`                     // Optional title for the passage
	FontLeft   string    `gorm:"default:serif" json:"font_left,omitempty"`      // Font for left panel: "serif" or "sans" (optional, falls back to StudyText)
	FontRight  string    `gorm:"default:sans" json:"font_right,omitempty"`      // Font for right panel: "serif" or "sans" (optional, falls back to StudyText)
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	// Relationship
	StudyText StudyText `gorm:"foreignKey:StudyTextID;references:ID" json:"study_text,omitempty"`
}

// QuizQuestion represents a quiz question for a study text
type QuizQuestion struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	StudyTextID uint     `gorm:"index;not null" json:"study_text_id"`
	QuestionID string    `gorm:"not null" json:"question_id"`  // e.g., "q1", "q2"
	Prompt     string    `gorm:"type:text;not null" json:"prompt"`
	Choices    string    `gorm:"type:text;not null" json:"choices"` // JSON array of choices
	Answer     int       `gorm:"not null" json:"answer"`             // Index of correct answer (0-based)
	Order      int       `gorm:"default:0" json:"order"`             // Display order
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	// Relationship
	StudyText StudyText `gorm:"foreignKey:StudyTextID;references:ID" json:"study_text,omitempty"`
}

