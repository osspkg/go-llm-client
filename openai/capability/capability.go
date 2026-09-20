// Package capability defines OpenAI-compatible operation capabilities.
package capability

// Name identifies one provider operation family.
type Name string

// Provider capabilities exposed by the OpenAI-compatible client.
const (
	Responses    Name = "responses"
	Chat         Name = "chat"
	Completions  Name = "completions"
	Embeddings   Name = "embeddings"
	Models       Name = "models"
	Files        Name = "files"
	Uploads      Name = "uploads"
	Batches      Name = "batches"
	FineTuning   Name = "fine_tuning"
	Audio        Name = "audio"
	Images       Name = "images"
	Moderations  Name = "moderations"
	Assistants   Name = "assistants"
	Threads      Name = "threads"
	Runs         Name = "runs"
	VectorStores Name = "vector_stores"
	Organization Name = "organization"
	Realtime     Name = "realtime"
	Containers   Name = "containers"
	Evals        Name = "evals"
)

// Matrix is an immutable snapshot of supported operation families.
type Matrix map[Name]bool

// Default returns the baseline OpenAI API capability set.
func Default() Matrix {
	return Matrix{
		Responses: true, Chat: true, Completions: true, Embeddings: true,
		Models: true, Files: true, Uploads: true, Batches: true,
		FineTuning: true, Audio: true, Images: true, Moderations: true,
		Assistants: true, Threads: true, Runs: true, VectorStores: true,
		Organization: true, Realtime: true, Containers: true, Evals: true,
	}
}

// Supports reports whether a capability is enabled.
func (m Matrix) Supports(name Name) bool { return m[name] }

// Clone returns an independent capability snapshot.
func (m Matrix) Clone() Matrix {
	copyMatrix := make(Matrix, len(m))
	for name, enabled := range m {
		copyMatrix[name] = enabled
	}
	return copyMatrix
}
