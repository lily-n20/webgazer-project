package main

import (
	"log"
)

// seedInitialData populates the database with initial study text and quiz questions
func seedInitialData() {
	// Check if study text already exists
	var count int64
	db.Model(&StudyText{}).Count(&count)
	if count > 0 {
		return // Data already seeded
	}

	// Create study text
	studyText := StudyText{
		Version:   "default",
		FontLeft:  "serif",
		FontRight: "sans",
		Active:    true,
	}

	if err := db.Create(&studyText).Error; err != nil {
		log.Printf("Error creating study text: %v", err)
		return
	}

	// Create multiple passages for the study text with different font combinations
	passages := []Passage{
		{
			StudyTextID: studyText.ID,
			Order:       0,
			Title:       "Passage 1: Introduction to Reading",
			Content:     `Reading is a complex cognitive process that involves decoding symbols to derive meaning. This process requires the coordination of multiple brain regions working together to transform written text into comprehensible information. The human brain processes visual information through the eyes, sending signals to various neural networks that interpret and understand the text.`,
			FontLeft:    "serif",
			FontRight:   "sans",
		},
		{
			StudyTextID: studyText.ID,
			Order:       1,
			Title:       "Passage 2: Typography and Readability",
			Content:     `Typography plays a crucial role in how we perceive and understand written content. Different font styles can significantly impact reading speed, comprehension, and overall user experience. Serif fonts, with their decorative strokes, are often associated with traditional print media, while sans-serif fonts offer a cleaner, more modern appearance.`,
			FontLeft:    "sans",
			FontRight:   "serif",
		},
		{
			StudyTextID: studyText.ID,
			Order:       2,
			Title:       "Passage 3: Reading Research",
			Content:     `Researchers have conducted extensive studies to understand how different typographic choices affect reading performance. These studies examine factors such as font size, line spacing, letter spacing, and font style. The goal is to identify optimal typography settings that maximize readability and comprehension for various audiences and contexts.`,
			FontLeft:    "serif",
			FontRight:   "sans",
		},
		{
			StudyTextID: studyText.ID,
			Order:       3,
			Title:       "Passage 4: Digital Reading",
			Content:     `The shift from print to digital media has introduced new challenges and opportunities in typography. Screen readability differs from print, requiring careful consideration of font rendering, display resolution, and viewing conditions. Designers must balance aesthetic appeal with functional readability to create effective digital reading experiences.`,
			FontLeft:    "sans",
			FontRight:   "serif",
		},
		{
			StudyTextID: studyText.ID,
			Order:       4,
			Title:       "Passage 5: Accessibility in Design",
			Content:     `Accessibility is a fundamental principle in modern design, ensuring that content is readable and understandable for people with diverse abilities and needs. This includes considerations for visual impairments, cognitive differences, and various reading contexts. Good typography choices can make content more accessible to a wider audience.`,
			FontLeft:    "serif",
			FontRight:   "sans",
		},
		{
			StudyTextID: studyText.ID,
			Order:       5,
			Title:       "Passage 6: The Future of Reading",
			Content:     `As technology continues to evolve, so too will our understanding of reading and typography. Emerging technologies like e-ink displays, variable fonts, and adaptive interfaces offer new possibilities for optimizing reading experiences. The future of typography lies in creating flexible, responsive designs that adapt to individual preferences and reading contexts.`,
			FontLeft:    "sans",
			FontRight:   "serif",
		},
	}

	for _, passage := range passages {
		if err := db.Create(&passage).Error; err != nil {
			log.Printf("Error creating passage: %v", err)
		}
	}

	// Create quiz questions
	questions := []QuizQuestion{
		{
			StudyTextID: studyText.ID,
			QuestionID:  "q1",
			Prompt:      "What is the purpose of this passage?",
			Choices:     `["To teach advanced speed-reading","To test font readability and comprehension","To explain eye-tracking algorithms","To measure typing accuracy"]`,
			Answer:      1,
			Order:       1,
		},
		{
			StudyTextID: studyText.ID,
			QuestionID:  "q2",
			Prompt:      "How should you read the passage?",
			Choices:     `["As quickly as possible without understanding","Only the first sentence","At a natural pace focusing on understanding","Backwards to test attention"]`,
			Answer:      2,
			Order:       2,
		},
		{
			StudyTextID: studyText.ID,
			QuestionID:  "q3",
			Prompt:      "According to the passage, what does reading involve?",
			Choices:     `["Only recognizing letters","Decoding symbols to derive meaning","Memorizing text word-for-word","Counting words per minute"]`,
			Answer:      1,
			Order:       3,
		},
		{
			StudyTextID: studyText.ID,
			QuestionID:  "q4",
			Prompt:      "What should you avoid when reading this passage?",
			Choices:     `["Reading at a natural pace","Focusing on understanding","Skimming through the content","Decoding the symbols"]`,
			Answer:      2,
			Order:       4,
		},
		{
			StudyTextID: studyText.ID,
			QuestionID:  "q5",
			Prompt:      "What is described as a \"complex cognitive process\"?",
			Choices:     `["Writing","Reading","Speaking","Listening"]`,
			Answer:      1,
			Order:       5,
		},
	}

	for _, q := range questions {
		if err := db.Create(&q).Error; err != nil {
			log.Printf("Error creating quiz question %s: %v", q.QuestionID, err)
		}
	}

	log.Println("Initial study data seeded successfully")
}

