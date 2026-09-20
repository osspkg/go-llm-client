// Package openai provides OpenAI-compatible and OpenAI API clients.
package openai

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go.osspkg.com/llm-client/openai/assistants"
	"go.osspkg.com/llm-client/openai/audio"
	"go.osspkg.com/llm-client/openai/batches"
	"go.osspkg.com/llm-client/openai/capability"
	"go.osspkg.com/llm-client/openai/chat"
	"go.osspkg.com/llm-client/openai/completions"
	"go.osspkg.com/llm-client/openai/containers"
	"go.osspkg.com/llm-client/openai/embeddings"
	"go.osspkg.com/llm-client/openai/evals"
	"go.osspkg.com/llm-client/openai/files"
	finetuning "go.osspkg.com/llm-client/openai/fine_tuning"
	"go.osspkg.com/llm-client/openai/images"
	"go.osspkg.com/llm-client/openai/models"
	"go.osspkg.com/llm-client/openai/moderations"
	"go.osspkg.com/llm-client/openai/organization"
	"go.osspkg.com/llm-client/openai/realtime"
	"go.osspkg.com/llm-client/openai/responses"
	"go.osspkg.com/llm-client/openai/runs"
	"go.osspkg.com/llm-client/openai/stateful"
	"go.osspkg.com/llm-client/openai/threads"
	"go.osspkg.com/llm-client/openai/uploads"
	"go.osspkg.com/llm-client/openai/vector_stores"
	"go.osspkg.com/llm-client/pkg/auth"
	llmerrors "go.osspkg.com/llm-client/pkg/errors"
	"go.osspkg.com/llm-client/pkg/transport"
)

const defaultBaseURL = "https://api.openai.com/v1"

// Option configures an OpenAI client.
type Option func(*config) error

type config struct {
	transportOptions []transport.Option
	capabilities     capability.Matrix
}

// Client is a concurrent-safe facade over OpenAI bounded contexts.
type Client struct {
	Responses    *responses.Client
	Chat         *chat.Client
	Completions  *completions.Client
	Assistants   *assistants.Client
	Threads      *threads.Client
	Runs         *runs.Client
	VectorStores *vector_stores.Client
	Containers   *containers.Client
	Evals        *evals.Client
	Embeddings   *embeddings.Client
	Models       *models.Client
	Files        *files.Client
	Uploads      *uploads.Client
	Batches      *batches.Client
	FineTuning   *finetuning.Client
	Audio        *audio.Client
	Images       *images.Client
	Moderations  *moderations.Client
	Stateful     *stateful.Client
	Organization *organization.Client
	Realtime     *realtime.Client
	capabilities capability.Matrix
}

// New creates an OpenAI-compatible client.
func New(options ...Option) (*Client, error) {
	configuration := config{capabilities: capability.Default()}
	baseURL := defaultBaseURL
	for _, option := range options {
		if option == nil {
			continue
		}
		if err := option(&configuration); err != nil {
			return nil, err
		}
	}
	configuration.transportOptions = append(configuration.transportOptions, transport.WithHeader("User-Agent", "go.osspkg.com/llm-client"))
	configuration.transportOptions = append(configuration.transportOptions, transport.WithOperationGate(func(operation string) error {
		family, _, _ := strings.Cut(operation, ".")
		name := capability.Name(family)
		if family != "" && !configuration.capabilities.Supports(name) {
			return fmt.Errorf("%w: capability %s is disabled", llmerrors.ErrInvalidRequest, name)
		}
		return nil
	}))
	client, err := transport.New(baseURL, configuration.transportOptions...)
	if err != nil {
		return nil, err
	}
	return &Client{
		Responses: responses.New(client), Chat: chat.New(client), Completions: completions.New(client),
		Assistants: assistants.New(client), Threads: threads.New(client), Runs: runs.New(client),
		VectorStores: vector_stores.New(client), Containers: containers.New(client), Evals: evals.New(client),
		Embeddings: embeddings.New(client), Models: models.New(client), Files: files.New(client),
		Uploads: uploads.New(client), Batches: batches.New(client), FineTuning: finetuning.New(client),
		Audio: audio.New(client), Images: images.New(client),
		Moderations: moderations.New(client), Stateful: stateful.New(client),
		Organization: organization.New(client), Realtime: realtime.New(client),
		capabilities: configuration.capabilities.Clone(),
	}, nil
}

// WithBaseURL changes the provider endpoint, including for OpenAI-compatible services.
func WithBaseURL(baseURL string) Option {
	return func(configuration *config) error {
		configuration.transportOptions = append(configuration.transportOptions, transport.WithBaseURL(baseURL))
		return nil
	}
}

// WithAuthProvider configures per-request authentication.
func WithAuthProvider(provider auth.HeaderProvider) Option {
	return func(configuration *config) error {
		configuration.transportOptions = append(configuration.transportOptions, transport.WithAuthProvider(provider))
		return nil
	}
}

// WithHTTPClient supplies a custom HTTP client, normally for tests.
func WithHTTPClient(client *http.Client) Option {
	return func(configuration *config) error {
		configuration.transportOptions = append(configuration.transportOptions, transport.WithHTTPClient(client))
		return nil
	}
}

// WithMaxResponseBody limits buffered provider responses.
func WithMaxResponseBody(limit int64) Option {
	return func(configuration *config) error {
		configuration.transportOptions = append(configuration.transportOptions, transport.WithMaxResponseBody(limit))
		return nil
	}
}

// WithMaxRequestBody limits streamed and buffered request bodies.
func WithMaxRequestBody(limit int64) Option {
	return func(configuration *config) error {
		configuration.transportOptions = append(configuration.transportOptions, transport.WithMaxRequestBody(limit))
		return nil
	}
}

// WithRequestTimeout sets the default non-stream request timeout.
func WithRequestTimeout(timeout time.Duration) Option {
	return func(configuration *config) error {
		configuration.transportOptions = append(configuration.transportOptions, transport.WithRequestTimeout(timeout))
		return nil
	}
}

// WithStreamTimeout sets an optional total stream timeout.
func WithStreamTimeout(timeout time.Duration) Option {
	return func(configuration *config) error {
		configuration.transportOptions = append(configuration.transportOptions, transport.WithStreamTimeout(timeout))
		return nil
	}
}

// WithCapability enables or disables an operation family.
func WithCapability(name capability.Name, enabled bool) Option {
	return func(configuration *config) error {
		configuration.capabilities[name] = enabled
		return nil
	}
}

// Capabilities returns an independent capability snapshot.
func (client *Client) Capabilities() capability.Matrix { return client.capabilities.Clone() }
