/**
 * API client for Readability Study Backend
 */

const API_BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080';

export interface StudySessionData {
	participant_id?: number;
	session_id?: string;
	calibration_points?: number;
	font_left?: string;
	font_right?: string;
	time_left_ms?: number;
	time_right_ms?: number;
	time_a_ms?: number;
	time_b_ms?: number;
	font_preference?: string;
	preferred_font_type?: string;
	quiz_responses_json?: string;
	user_agent?: string;
	screen_width?: number;
	screen_height?: number;
}

export interface QuizResponseData {
	session_id: number;
	question_id: string;
	answer_index: number;
	is_correct?: boolean;
	response_time?: number;
}

export interface ApiResponse {
	success: boolean;
	session_id?: string;
	id?: number;
	error?: string;
}

/**
 * Create or get a participant
 */
export async function createParticipant(source: string = 'web'): Promise<number> {
	// Check if participant ID already exists in sessionStorage
	const existingId = sessionStorage.getItem('participant_id');
	if (existingId) {
		return parseInt(existingId, 10);
	}

	try {
		const response = await fetch(`${API_BASE_URL}/api/participant`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({ source })
		});

		if (!response.ok) {
			throw new Error(`Failed to create participant: ${response.statusText}`);
		}

		const data = await response.json();
		const participantId = data.id;

		// Store in sessionStorage for reuse
		sessionStorage.setItem('participant_id', String(participantId));

		return participantId;
	} catch (error) {
		console.error('Error creating participant:', error);
		// Return a temporary ID if API fails
		const tempId = Date.now();
		sessionStorage.setItem('participant_id', String(tempId));
		return tempId;
	}
}

/**
 * Submit a study session to the backend
 */
export async function submitStudySession(data: StudySessionData): Promise<ApiResponse> {
	try {
		// Ensure participant_id exists
		if (!data.participant_id) {
			data.participant_id = await createParticipant();
		}

		// Add metadata if not provided
		if (!data.user_agent) {
			data.user_agent = navigator.userAgent;
		}
		if (!data.screen_width) {
			data.screen_width = window.screen.width;
		}
		if (!data.screen_height) {
			data.screen_height = window.screen.height;
		}

		const response = await fetch(`${API_BASE_URL}/api/session`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});

		if (!response.ok) {
			const errorText = await response.text();
			throw new Error(`Failed to submit session: ${errorText}`);
		}

		const result = await response.json();

		// Store session ID for later use
		if (result.session_id) {
			sessionStorage.setItem('session_id', result.session_id);
		}
		if (result.id) {
			sessionStorage.setItem('session_db_id', String(result.id));
		}

		return result;
	} catch (error) {
		console.error('Error submitting study session:', error);
		return {
			success: false,
			error: error instanceof Error ? error.message : 'Unknown error'
		};
	}
}

/**
 * Submit individual quiz responses
 */
export async function submitQuizResponse(data: QuizResponseData): Promise<boolean> {
	try {
		const response = await fetch(`${API_BASE_URL}/api/quiz-response`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});

		if (!response.ok) {
			const errorText = await response.text();
			console.error(`Failed to submit quiz response: ${response.status} ${errorText}`, data);
			return false;
		}

		return true;
	} catch (error) {
		console.error('Error submitting quiz response:', error, data);
		return false;
	}
}

/**
 * Submit calibration data point
 */
export async function submitCalibrationData(data: {
	session_id: number;
	point_index: number;
	click_number: number;
	x: number;
	y: number;
}): Promise<boolean> {
	try {
		const response = await fetch(`${API_BASE_URL}/api/calibration`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});

		return response.ok;
	} catch (error) {
		console.error('Error submitting calibration data:', error);
		return false;
	}
}

/**
 * Submit accuracy measurement
 */
export async function submitAccuracyMeasurement(data: {
	session_id: number;
	accuracy: number;
	duration: number;
	passed: boolean;
}): Promise<boolean> {
	try {
		const response = await fetch(`${API_BASE_URL}/api/accuracy`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});

		return response.ok;
	} catch (error) {
		console.error('Error submitting accuracy measurement:', error);
		return false;
	}
}

/**
 * Submit gaze point
 */
export async function submitGazePoint(data: {
	session_id: number;
	x: number;
	y: number;
	panel?: string;
	phase?: string;
}): Promise<boolean> {
	try {
		const response = await fetch(`${API_BASE_URL}/api/gaze-point`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});

		return response.ok;
	} catch (error) {
		console.error('Error submitting gaze point:', error);
		return false;
	}
}

/**
 * Submit reading event
 */
export async function submitReadingEvent(data: {
	session_id: number;
	event_type: string;
	panel: string;
	duration?: number;
}): Promise<boolean> {
	try {
		const response = await fetch(`${API_BASE_URL}/api/reading-event`, {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify(data)
		});

		return response.ok;
	} catch (error) {
		console.error('Error submitting reading event:', error);
		return false;
	}
}

/**
 * Collect all session data from sessionStorage and submit
 */
export async function submitCompleteSession(quizAnswers: Record<string, number>): Promise<boolean> {
	try {
		// Collect data from sessionStorage
		const sessionData: StudySessionData = {
			participant_id: parseInt(sessionStorage.getItem('participant_id') || '0', 10) || undefined,
			calibration_points:
				parseInt(sessionStorage.getItem('calibration_points') || '0', 10) || undefined,
			font_left: sessionStorage.getItem('font_left') || undefined,
			font_right: sessionStorage.getItem('font_right') || undefined,
			time_left_ms: parseInt(sessionStorage.getItem('time_left_ms') || '0', 10) || undefined,
			time_right_ms: parseInt(sessionStorage.getItem('time_right_ms') || '0', 10) || undefined,
			time_a_ms: parseInt(sessionStorage.getItem('timeA_ms') || '0', 10) || undefined,
			time_b_ms: parseInt(sessionStorage.getItem('timeB_ms') || '0', 10) || undefined,
			font_preference: sessionStorage.getItem('font_preference') || undefined,
			preferred_font_type: sessionStorage.getItem('font_preferred_type') || undefined
		};

		// Convert quiz answers to JSON
		const quizResponses = Object.entries(quizAnswers).map(([questionId, answerIndex]) => ({
			question_id: questionId,
			answer: answerIndex
		}));
		sessionData.quiz_responses_json = JSON.stringify(quizResponses);

		// Submit to backend
		const result = await submitStudySession(sessionData);

		if (result.success && result.id) {
			// Also submit individual quiz responses to the quiz_responses table
			const sessionDbId = result.id;
			const quizQuestions = await import('$lib/studyText').then((m) => m.QUIZ);

			// Submit each quiz response individually
			const quizSubmissionPromises = [];
			for (const [questionId, answerIndex] of Object.entries(quizAnswers)) {
				const question = quizQuestions.find((q) => q.id === questionId);
				const isCorrect = question ? answerIndex === question.answer : undefined;

				quizSubmissionPromises.push(
					submitQuizResponse({
						session_id: sessionDbId,
						question_id: questionId,
						answer_index: answerIndex,
						is_correct: isCorrect
					}).catch((error) => {
						console.error(`Failed to submit quiz response for ${questionId}:`, error);
						return false;
					})
				);
			}

			// Wait for all quiz responses to be submitted (but don't fail if some fail)
			const results = await Promise.all(quizSubmissionPromises);
			const successCount = results.filter((r) => r === true).length;
			console.log(
				`Submitted ${successCount}/${quizSubmissionPromises.length} quiz responses individually`
			);

			return true;
		}

		return false;
	} catch (error) {
		console.error('Error submitting complete session:', error);
		return false;
	}
}
