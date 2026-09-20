package uploads

import "io"

//go:generate easyjson -all types.go

// CreateRequest starts a large file upload.
type CreateRequest struct {
	Filename string `json:"filename"`
	Purpose  string `json:"purpose"`
	Bytes    int64  `json:"bytes"`
	MimeType string `json:"mime_type"`
}

// Upload describes an in-progress or completed upload.
type Upload struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int64  `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
	Status    string `json:"status"`
}

// Part describes a completed upload part.
type Part struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt int64  `json:"created_at"`
	UploadID  string `json:"upload_id"`
	Bytes     int64  `json:"bytes"`
}

// CompleteRequest completes an upload with ordered part IDs.
type CompleteRequest struct {
	PartIDs []string `json:"part_ids"`
}

// PartInput contains one bounded upload part reader.
//
//easyjson:skip
type PartInput struct {
	Filename string
	File     io.Reader
}

// DeleteResponse confirms an upload cancellation.
type DeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}
