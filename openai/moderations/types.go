package moderations

//go:generate easyjson -all types.go

// Request classifies input content.
type Request struct {
	// Model optionally selects a moderation model.
	Model string `json:"model,omitempty"`
	// Input is the content to classify.
	Input string `json:"input"`
}

// Response contains moderation results.
type Response struct {
	// ID identifies the moderation request.
	ID string `json:"id"`
	// Model identifies the model that classified the input.
	Model string `json:"model"`
	// Results contains one classification result per input.
	Results []Result `json:"results"`
}

// Result contains moderation category scores.
type Result struct {
	// Flagged reports whether the provider flagged the input.
	Flagged bool `json:"flagged"`
	// Categories maps category names to their flag state.
	Categories map[string]bool `json:"categories"`
	// Scores maps category names to confidence scores.
	Scores map[string]float64 `json:"category_scores"`
}
