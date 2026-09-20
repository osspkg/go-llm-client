package moderations

//go:generate easyjson -all types.go

// Request classifies input content.
type Request struct {
	Model string `json:"model,omitempty"`
	Input string `json:"input"`
}

// Response contains moderation results.
type Response struct {
	ID      string   `json:"id"`
	Model   string   `json:"model"`
	Results []Result `json:"results"`
}

// Result contains moderation category scores.
type Result struct {
	Flagged    bool               `json:"flagged"`
	Categories map[string]bool    `json:"categories"`
	Scores     map[string]float64 `json:"category_scores"`
}
