# Frontend Integration Example

To send data from your Svelte frontend to the backend, you can add this function:

```typescript
// In your quiz page or wherever you want to submit data
async function submitStudyData() {
  const sessionData = {
    session_id: sessionStorage.getItem("session_id") || undefined,
    calibration_points: 25, // or get from calibration page
    font_left: sessionStorage.getItem("font_left"),
    font_right: sessionStorage.getItem("font_right"),
    time_left_ms: parseInt(sessionStorage.getItem("time_left_ms") || "0"),
    time_right_ms: parseInt(sessionStorage.getItem("time_right_ms") || "0"),
    time_a_ms: parseInt(sessionStorage.getItem("timeA_ms") || "0"),
    time_b_ms: parseInt(sessionStorage.getItem("timeB_ms") || "0"),
    font_preference: sessionStorage.getItem("font_preference"),
    preferred_font_type: sessionStorage.getItem("font_preferred_type"),
    quiz_responses: JSON.stringify(
      Object.entries(answers).map(([questionId, answerIndex]) => ({
        question_id: questionId,
        answer: answerIndex,
      }))
    ),
    user_agent: navigator.userAgent,
    screen_width: window.screen.width,
    screen_height: window.screen.height,
  };

  try {
    const response = await fetch("http://localhost:8080/api/session", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(sessionData),
    });

    if (response.ok) {
      const result = await response.json();
      console.log("Session saved:", result);
    } else {
      console.error("Failed to save session:", await response.text());
    }
  } catch (error) {
    console.error("Error submitting data:", error);
  }
}
```

Call this function when the quiz is submitted or at the end of the study flow.
