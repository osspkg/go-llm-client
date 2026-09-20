package finetuning

import "encoding/json"

//go:generate easyjson -all types.go

// CreateRequest starts a fine-tuning job.
type CreateRequest struct {
	// Model selects the model used for the operation.
	Model string `json:"model"`
	// TrainingFile identifies the fine-tuning training file.
	TrainingFile string `json:"training_file"`
	// ValidationFile identifies the optional fine-tuning validation file.
	ValidationFile string `json:"validation_file,omitempty"`
	// Hyperparameters contains the provider value used for the `Hyperparameters` field.
	Hyperparameters json.RawMessage `json:"hyperparameters,omitempty"`
	// Suffix contains the provider value used for the `Suffix` field.
	Suffix string `json:"suffix,omitempty"`
	// Integrations contains the provider value used for the `Integrations` field.
	Integrations json.RawMessage `json:"integrations,omitempty"`
	// Seed contains the provider value used for the `Seed` field.
	Seed int `json:"seed,omitempty"`
	// Metadata contains caller-defined key-value labels attached to the resource.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// Job describes a fine-tuning job.
type Job struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Model selects the model used for the operation.
	Model string `json:"model"`
	// FineTunedModel identifies the resulting fine-tuned model.
	FineTunedModel string `json:"fine_tuned_model,omitempty"`
	// Status reports the current provider processing or lifecycle state.
	Status string `json:"status"`
	// TrainingFile identifies the fine-tuning training file.
	TrainingFile string `json:"training_file"`
	// ValidationFile identifies the optional fine-tuning validation file.
	ValidationFile string `json:"validation_file,omitempty"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at"`
	// FinishedAt contains the Unix timestamp when processing finished.
	FinishedAt int64 `json:"finished_at,omitempty"`
	// Epochs contains the number of training epochs.
	Epochs int `json:"epochs,omitempty"`
	// Error contains provider error details.
	Error *Error `json:"error,omitempty"`
}

// Error describes a fine-tuning failure.
type Error struct {
	// Code contains the provider error code.
	Code string `json:"code,omitempty"`
	// Message contains the provider error message.
	Message string `json:"message,omitempty"`
	// Param identifies the request parameter associated with the error.
	Param string `json:"param,omitempty"`
}

// ListResponse lists fine-tuning jobs.
type ListResponse struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Data contains the returned records.
	Data []Job `json:"data"`
	// HasMore reports whether another page is available.
	HasMore bool `json:"has_more"`
	// FirstID identifies the first resource in a cursor page.
	FirstID string `json:"first_id,omitempty"`
	// LastID identifies the last resource in a cursor page.
	LastID string `json:"last_id,omitempty"`
}

// Event is one fine-tuning event.
type Event struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id,omitempty"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at"`
	// Level identifies the event severity.
	Level string `json:"level,omitempty"`
	// Message contains the provider error message.
	Message string `json:"message,omitempty"`
	// Type identifies the object or event variant.
	Type string `json:"type,omitempty"`
}

// ListEventsResponse is a cursor page of fine-tuning job events.
type ListEventsResponse struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Data contains the returned records.
	Data []Event `json:"data"`
	// FirstID identifies the first resource in a cursor page.
	FirstID string `json:"first_id,omitempty"`
	// LastID identifies the last resource in a cursor page.
	LastID string `json:"last_id,omitempty"`
	// HasMore reports whether another page is available.
	HasMore bool `json:"has_more"`
}

// CheckpointMetrics contains measurements captured for a fine-tuning checkpoint.
type CheckpointMetrics struct {
	// Step contains the training step number.
	Step float64 `json:"step,omitempty"`
	// TrainLoss contains training loss at the checkpoint.
	TrainLoss float64 `json:"train_loss,omitempty"`
	// TrainMeanTokenAccuracy contains training token accuracy.
	TrainMeanTokenAccuracy float64 `json:"train_mean_token_accuracy,omitempty"`
	// ValidLoss contains validation loss at the checkpoint.
	ValidLoss float64 `json:"valid_loss,omitempty"`
	// ValidMeanTokenAccuracy contains validation token accuracy.
	ValidMeanTokenAccuracy float64 `json:"valid_mean_token_accuracy,omitempty"`
	// FullValidLoss contains full-validation loss at the checkpoint.
	FullValidLoss float64 `json:"full_valid_loss,omitempty"`
	// FullValidMeanTokenAccuracy contains full-validation token accuracy.
	FullValidMeanTokenAccuracy float64 `json:"full_valid_mean_token_accuracy,omitempty"`
}

// Checkpoint is a usable fine-tuning model checkpoint.
type Checkpoint struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at"`
	// FineTunedModelCheckpoint identifies the checkpoint model.
	FineTunedModelCheckpoint string `json:"fine_tuned_model_checkpoint"`
	// FineTuningJobID identifies the fine-tuning job.
	FineTuningJobID string `json:"fine_tuning_job_id"`
	// StepNumber contains the training step represented by the checkpoint.
	StepNumber int `json:"step_number"`
	// Metrics contains measurements captured at the checkpoint.
	Metrics CheckpointMetrics `json:"metrics"`
}

// ListCheckpointsResponse is a cursor page of fine-tuning checkpoints.
type ListCheckpointsResponse struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Data contains the returned records.
	Data []Checkpoint `json:"data"`
	// FirstID identifies the first resource in a cursor page.
	FirstID string `json:"first_id,omitempty"`
	// LastID identifies the last resource in a cursor page.
	LastID string `json:"last_id,omitempty"`
	// HasMore reports whether another page is available.
	HasMore bool `json:"has_more"`
}

// CheckpointPermission grants a project access to a fine-tuning checkpoint.
type CheckpointPermission struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// CreatedAt contains the Unix creation timestamp.
	CreatedAt int64 `json:"created_at"`
	// ProjectID identifies the project receiving access.
	ProjectID string `json:"project_id"`
}

// CreateCheckpointPermissionRequest grants one or more projects access to a checkpoint.
type CreateCheckpointPermissionRequest struct {
	// ProjectIDs lists projects receiving access.
	ProjectIDs []string `json:"project_ids"`
}

// ListCheckpointPermissionsResponse is a cursor page of checkpoint permissions.
type ListCheckpointPermissionsResponse struct {
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Data contains the returned records.
	Data []CheckpointPermission `json:"data"`
	// FirstID identifies the first resource in a cursor page.
	FirstID string `json:"first_id,omitempty"`
	// LastID identifies the last resource in a cursor page.
	LastID string `json:"last_id,omitempty"`
	// HasMore reports whether another page is available.
	HasMore bool `json:"has_more"`
}

// DeleteCheckpointPermissionResponse confirms deletion of a checkpoint permission.
type DeleteCheckpointPermissionResponse struct {
	// ID identifies the resource, event, or request.
	ID string `json:"id"`
	// Object identifies the provider resource or collection type.
	Object string `json:"object"`
	// Deleted confirms that the provider deleted the resource.
	Deleted bool `json:"deleted"`
}

// RunGraderRequest executes an upstream-defined alpha grader. Grader and Item
// are arbitrary JSON because the OpenAPI contract deliberately defines them as
// open objects and a discriminator-based union.
type RunGraderRequest struct {
	// Grader contains the provider-defined grader configuration.
	Grader json.RawMessage `json:"grader"`
	// Item contains the item submitted to a grader.
	Item json.RawMessage `json:"item,omitempty"`
	// ModelSample contains the model output evaluated by the grader.
	ModelSample string `json:"model_sample"`
}

// RunGraderResponse contains a grader reward and its spec-defined metadata.
type RunGraderResponse struct {
	// Reward contains the grader reward.
	Reward float64 `json:"reward"`
	// Metadata contains caller-defined key-value labels attached to the resource.
	Metadata json.RawMessage `json:"metadata,omitempty"`
	// SubRewards contains per-component grader rewards.
	SubRewards json.RawMessage `json:"sub_rewards,omitempty"`
	// ModelGraderTokenUsagePerModel contains the provider value used for the `ModelGraderTokenUsagePerModel` field.
	ModelGraderTokenUsagePerModel json.RawMessage `json:"model_grader_token_usage_per_model,omitempty"`
}

// ValidateGraderRequest validates an upstream-defined alpha grader union.
type ValidateGraderRequest struct {
	// Grader contains the provider-defined grader configuration.
	Grader json.RawMessage `json:"grader"`
}

// ValidateGraderResponse returns the normalized grader definition.
type ValidateGraderResponse struct {
	// Grader contains the provider-defined grader configuration.
	Grader json.RawMessage `json:"grader"`
}
