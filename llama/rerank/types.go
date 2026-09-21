// Package rerank provides typed rerank API operations.
package rerank

import "encoding/json"

//go:generate easyjson -all types.go

// Request asks the reranker to score documents against a query.
type Request struct {
	// Query is the natural-language question or search intent.
	Query string `json:"query"`
	// Documents contains candidate texts to rank.
	Documents []string `json:"documents"`
	// TopN limits the number of returned ranked documents when supported.
	TopN int `json:"top_n,omitempty"`
}

// Result is one ranked document.
type Result struct {
	// Index is the zero-based index in the submitted document list.
	Index int `json:"index"`
	// RelevanceScore is the provider's relevance score for the document.
	RelevanceScore float64 `json:"relevance_score"`
	// Document contains the returned document representation when supplied.
	Document json.RawMessage `json:"document,omitempty"`
}

// Response contains ranked documents ordered by relevance.
type Response struct {
	// Results contains the ranked document results.
	Results []Result `json:"results"`
}
