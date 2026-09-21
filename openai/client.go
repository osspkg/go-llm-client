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
	responses    *responses.Client
	chat         *chat.Client
	completions  *completions.Client
	assistants   *assistants.Client
	threads      *threads.Client
	runs         *runs.Client
	vectorStores *vector_stores.Client
	containers   *containers.Client
	evals        *evals.Client
	embeddings   *embeddings.Client
	models       *models.Client
	files        *files.Client
	uploads      *uploads.Client
	batches      *batches.Client
	fineTuning   *finetuning.Client
	audio        *audio.Client
	images       *images.Client
	moderations  *moderations.Client
	stateful     *stateful.Client
	organization *organization.Client
	realtime     *realtime.Client
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
		responses: responses.New(client), chat: chat.New(client), completions: completions.New(client),
		assistants: assistants.New(client), threads: threads.New(client), runs: runs.New(client),
		vectorStores: vector_stores.New(client), containers: containers.New(client), evals: evals.New(client),
		embeddings: embeddings.New(client), models: models.New(client), files: files.New(client),
		uploads: uploads.New(client), batches: batches.New(client), fineTuning: finetuning.New(client),
		audio: audio.New(client), images: images.New(client),
		moderations: moderations.New(client), stateful: stateful.New(client),
		organization: organization.New(client), realtime: realtime.New(client),
		capabilities: configuration.capabilities.Clone(),
	}, nil
}

// Responses returns the immutable Responses domain client.
func (client *Client) Responses() *responses.Client { return client.responses }

// Chat returns the immutable Chat Completions domain client.
func (client *Client) Chat() *chat.Client { return client.chat }

// Completions returns the immutable legacy Completions domain client.
func (client *Client) Completions() *completions.Client { return client.completions }

// Assistants returns the immutable Assistants domain client.
func (client *Client) Assistants() *assistants.Client { return client.assistants }

// Threads returns the immutable Threads domain client.
func (client *Client) Threads() *threads.Client { return client.threads }

// Runs returns the immutable Runs domain client.
func (client *Client) Runs() *runs.Client { return client.runs }

// VectorStores returns the immutable Vector Stores domain client.
func (client *Client) VectorStores() *vector_stores.Client { return client.vectorStores }

// Containers returns the immutable Containers domain client.
func (client *Client) Containers() *containers.Client { return client.containers }

// Evals returns the immutable Evals domain client.
func (client *Client) Evals() *evals.Client { return client.evals }

// Embeddings returns the immutable Embeddings domain client.
func (client *Client) Embeddings() *embeddings.Client { return client.embeddings }

// Models returns the immutable Models domain client.
func (client *Client) Models() *models.Client { return client.models }

// Files returns the immutable Files domain client.
func (client *Client) Files() *files.Client { return client.files }

// Uploads returns the immutable Uploads domain client.
func (client *Client) Uploads() *uploads.Client { return client.uploads }

// Batches returns the immutable Batches domain client.
func (client *Client) Batches() *batches.Client { return client.batches }

// FineTuning returns the immutable Fine-Tuning domain client.
func (client *Client) FineTuning() *finetuning.Client { return client.fineTuning }

// Audio returns the immutable Audio domain client.
func (client *Client) Audio() *audio.Client { return client.audio }

// Images returns the immutable Images domain client.
func (client *Client) Images() *images.Client { return client.images }

// Moderations returns the immutable Moderations domain client.
func (client *Client) Moderations() *moderations.Client { return client.moderations }

// Stateful returns the immutable stateful-resource domain client.
func (client *Client) Stateful() *stateful.Client { return client.stateful }

// Organization returns the immutable Organization and Admin domain client.
func (client *Client) Organization() *organization.Client { return client.organization }

// Realtime returns the immutable Realtime WebSocket domain client.
func (client *Client) Realtime() *realtime.Client { return client.realtime }

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
