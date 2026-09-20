package files

//go:generate easyjson -all types.go

// File describes an uploaded file.
type File struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int64  `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
	Status    string `json:"status,omitempty"`
}

// ListResponse lists files.
type ListResponse struct {
	Object  string `json:"object"`
	Data    []File `json:"data"`
	HasMore bool   `json:"has_more"`
}

// DeleteResponse confirms deletion.
type DeleteResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}
