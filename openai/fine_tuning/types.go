package finetuning

import "encoding/json"

//go:generate easyjson -all types.go

// CreateRequest starts a fine-tuning job.
type CreateRequest struct {
	Model           string            `json:"model"`
	TrainingFile    string            `json:"training_file"`
	ValidationFile  string            `json:"validation_file,omitempty"`
	Hyperparameters json.RawMessage   `json:"hyperparameters,omitempty"`
	Suffix          string            `json:"suffix,omitempty"`
	Integrations    json.RawMessage   `json:"integrations,omitempty"`
	Seed            int               `json:"seed,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// Job describes a fine-tuning job.
type Job struct {
	ID             string `json:"id"`
	Object         string `json:"object"`
	Model          string `json:"model"`
	FineTunedModel string `json:"fine_tuned_model,omitempty"`
	Status         string `json:"status"`
	TrainingFile   string `json:"training_file"`
	ValidationFile string `json:"validation_file,omitempty"`
	CreatedAt      int64  `json:"created_at"`
	FinishedAt     int64  `json:"finished_at,omitempty"`
	Epochs         int    `json:"epochs,omitempty"`
	Error          *Error `json:"error,omitempty"`
}

// Error describes a fine-tuning failure.
type Error struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Param   string `json:"param,omitempty"`
}

// ListResponse lists fine-tuning jobs.
type ListResponse struct {
	Object  string `json:"object"`
	Data    []Job  `json:"data"`
	HasMore bool   `json:"has_more"`
	FirstID string `json:"first_id,omitempty"`
	LastID  string `json:"last_id,omitempty"`
}

// Event is one fine-tuning event.
type Event struct {
	ID        string `json:"id,omitempty"`
	Object    string `json:"object"`
	CreatedAt int64  `json:"created_at"`
	Level     string `json:"level,omitempty"`
	Message   string `json:"message,omitempty"`
	Type      string `json:"type,omitempty"`
}

// ListEventsResponse is a cursor page of fine-tuning job events.
type ListEventsResponse struct {
	Object  string  `json:"object"`
	Data    []Event `json:"data"`
	FirstID string  `json:"first_id,omitempty"`
	LastID  string  `json:"last_id,omitempty"`
	HasMore bool    `json:"has_more"`
}

// CheckpointMetrics contains measurements captured for a fine-tuning checkpoint.
type CheckpointMetrics struct {
	Step                       float64 `json:"step,omitempty"`
	TrainLoss                  float64 `json:"train_loss,omitempty"`
	TrainMeanTokenAccuracy     float64 `json:"train_mean_token_accuracy,omitempty"`
	ValidLoss                  float64 `json:"valid_loss,omitempty"`
	ValidMeanTokenAccuracy     float64 `json:"valid_mean_token_accuracy,omitempty"`
	FullValidLoss              float64 `json:"full_valid_loss,omitempty"`
	FullValidMeanTokenAccuracy float64 `json:"full_valid_mean_token_accuracy,omitempty"`
}

// Checkpoint is a usable fine-tuning model checkpoint.
type Checkpoint struct {
	ID                       string            `json:"id"`
	Object                   string            `json:"object"`
	CreatedAt                int64             `json:"created_at"`
	FineTunedModelCheckpoint string            `json:"fine_tuned_model_checkpoint"`
	FineTuningJobID          string            `json:"fine_tuning_job_id"`
	StepNumber               int               `json:"step_number"`
	Metrics                  CheckpointMetrics `json:"metrics"`
}

// ListCheckpointsResponse is a cursor page of fine-tuning checkpoints.
type ListCheckpointsResponse struct {
	Object  string       `json:"object"`
	Data    []Checkpoint `json:"data"`
	FirstID string       `json:"first_id,omitempty"`
	LastID  string       `json:"last_id,omitempty"`
	HasMore bool         `json:"has_more"`
}

// CheckpointPermission grants a project access to a fine-tuning checkpoint.
type CheckpointPermission struct {
	ID        string `json:"id"`
	Object    string `json:"object"`
	CreatedAt int64  `json:"created_at"`
	ProjectID string `json:"project_id"`
}

// CreateCheckpointPermissionRequest grants one or more projects access to a checkpoint.
type CreateCheckpointPermissionRequest struct {
	ProjectIDs []string `json:"project_ids"`
}

// ListCheckpointPermissionsResponse is a cursor page of checkpoint permissions.
type ListCheckpointPermissionsResponse struct {
	Object  string                 `json:"object"`
	Data    []CheckpointPermission `json:"data"`
	FirstID string                 `json:"first_id,omitempty"`
	LastID  string                 `json:"last_id,omitempty"`
	HasMore bool                   `json:"has_more"`
}

// DeleteCheckpointPermissionResponse confirms deletion of a checkpoint permission.
type DeleteCheckpointPermissionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

// RunGraderRequest executes an upstream-defined alpha grader. Grader and Item
// are arbitrary JSON because the OpenAPI contract deliberately defines them as
// open objects and a discriminator-based union.
type RunGraderRequest struct {
	Grader      json.RawMessage `json:"grader"`
	Item        json.RawMessage `json:"item,omitempty"`
	ModelSample string          `json:"model_sample"`
}

// RunGraderResponse contains a grader reward and its spec-defined metadata.
type RunGraderResponse struct {
	Reward                        float64         `json:"reward"`
	Metadata                      json.RawMessage `json:"metadata,omitempty"`
	SubRewards                    json.RawMessage `json:"sub_rewards,omitempty"`
	ModelGraderTokenUsagePerModel json.RawMessage `json:"model_grader_token_usage_per_model,omitempty"`
}

// ValidateGraderRequest validates an upstream-defined alpha grader union.
type ValidateGraderRequest struct {
	Grader json.RawMessage `json:"grader"`
}

// ValidateGraderResponse returns the normalized grader definition.
type ValidateGraderResponse struct {
	Grader json.RawMessage `json:"grader"`
}
