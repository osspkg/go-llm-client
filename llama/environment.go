/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package llama

import "go.osspkg.com/llm-client/pkg/environment"

// EnvironmentScheme describes the documented llama-server environment
// variables for networking, model loading, memory, generation, routing, and
// endpoint exposure. It is independent of a Client and does not inspect the
// process environment.
func EnvironmentScheme() environment.Scheme {
	return environment.Scheme{Variables: []environment.Variable{
		{
			Name:        "LLAMA_ARG_MODEL",
			Description: "Sets the local GGUF model path loaded by llama-server; leave unset when using a router or a model URL.",
		},
		{
			Name:        "LLAMA_ARG_MODEL_URL",
			Description: "Sets a URL from which llama-server downloads the model when a local model path is not used.",
		},
		{
			Name:        "LLAMA_ARG_HF_REPO",
			Description: "Selects a Hugging Face model repository and optional quantization tag, for example ggml-org/model:Q4_K_M.",
		},
		{
			Name:        "LLAMA_ARG_HF_FILE",
			Description: "Selects the model file inside the Hugging Face repository and overrides the repository's quantization selection.",
		},
		{
			Name:        "HF_TOKEN",
			Description: "Provides the Hugging Face access token used to download private or rate-limited model assets; do not include it in logs or URLs.",
		},
		{
			Name:        "LLAMA_ARG_HOST",
			Default:     "127.0.0.1",
			Description: "Sets the address where llama-server listens; an address ending in .sock selects a Unix socket.",
		},
		{
			Name:        "LLAMA_ARG_PORT",
			Default:     "8080",
			Description: "Sets the TCP port where llama-server listens; use a valid port number from 1 through 65535.",
		},
		{
			Name:        "LLAMA_ARG_API_PREFIX",
			Description: "Adds a URL path prefix to all server routes without a trailing slash; leave unset to serve routes at the root.",
		},
		{
			Name:        "LLAMA_ARG_CORS_ORIGINS",
			Default:     "*",
			Description: "Sets a comma-separated CORS origin allowlist; use localhost instead of * when browser access must stay local.",
		},
		{
			Name:        "LLAMA_ARG_CORS_METHODS",
			Default:     "GET, POST, DELETE, OPTIONS",
			Description: "Sets the comma-separated HTTP methods allowed by the CORS middleware.",
		},
		{
			Name:        "LLAMA_ARG_CORS_HEADERS",
			Default:     "*",
			Description: "Sets the comma-separated request headers allowed by the CORS middleware.",
		},
		{
			Name:          "LLAMA_ARG_CORS_CREDENTIALS",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Controls whether CORS requests may include credentials; avoid combining enabled credentials with a wildcard origin in untrusted deployments.",
		},
		{
			Name:        "LLAMA_ARG_THREADS",
			Default:     "-1",
			Description: "Sets the number of CPU threads used during generation; -1 lets llama.cpp choose automatically.",
		},
		{
			Name:        "LLAMA_ARG_THREADS_BATCH",
			Description: "Sets the number of CPU threads used for prompt and batch processing; unset follows LLAMA_ARG_THREADS.",
		},
		{
			Name:        "LLAMA_ARG_CTX_SIZE",
			Default:     "0",
			Description: "Sets the prompt context size in tokens; zero loads the context size from the model metadata.",
		},
		{
			Name:        "LLAMA_ARG_N_PREDICT",
			Default:     "-1",
			Description: "Sets the maximum number of generated tokens; -1 means no explicit generation limit.",
		},
		{
			Name:        "LLAMA_ARG_BATCH",
			Default:     "2048",
			Description: "Sets the logical maximum batch size used during prompt processing.",
		},
		{
			Name:        "LLAMA_ARG_UBATCH",
			Default:     "512",
			Description: "Sets the physical maximum batch size used by the backend.",
		},
		{
			Name:          "LLAMA_ARG_FLASH_ATTN",
			AllowedValues: []string{"on", "off", "auto"},
			Default:       "auto",
			Description:   "Controls Flash Attention usage during inference.",
		},
		{
			Name:          "LLAMA_ARG_ROPE_SCALING_TYPE",
			AllowedValues: []string{"none", "linear", "yarn"},
			Default:       "model-dependent",
			Description:   "Selects the RoPE frequency scaling algorithm used when extending the model context.",
		},
		{
			Name:        "LLAMA_ARG_ROPE_SCALE",
			Description: "Sets the RoPE context scaling factor; use a positive number to expand the effective context.",
		},
		{
			Name:          "LLAMA_ARG_KV_OFFLOAD",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Controls whether the key/value cache is offloaded to the selected accelerator.",
		},
		{
			Name:          "LLAMA_ARG_CACHE_TYPE_K",
			AllowedValues: []string{"f32", "f16", "bf16", "q8_0", "q4_0", "q4_1", "iq4_nl", "q5_0", "q5_1"},
			Default:       "f16",
			Description:   "Selects the data type used by the key part of the KV cache.",
		},
		{
			Name:          "LLAMA_ARG_CACHE_TYPE_V",
			AllowedValues: []string{"f32", "f16", "bf16", "q8_0", "q4_0", "q4_1", "iq4_nl", "q5_0", "q5_1"},
			Default:       "f16",
			Description:   "Selects the data type used by the value part of the KV cache.",
		},
		{
			Name:          "LLAMA_ARG_LOAD_MODE",
			AllowedValues: []string{"auto", "none", "mmap", "mlock", "mmap+mlock", "dio"},
			Default:       "auto",
			Description:   "Controls how model weights are loaded from storage: memory mapping, memory locking, DirectIO, or automatic selection.",
		},
		{
			Name:          "LLAMA_ARG_LAZY_MODE",
			AllowedValues: []string{"on", "auto", "off"},
			Default:       "auto",
			Description:   "Controls on-demand loading of selected tensors, which can reduce resident memory for large models.",
		},
		{
			Name:          "LLAMA_ARG_NUMA",
			AllowedValues: []string{"distribute", "isolate", "numactl"},
			Description:   "Selects NUMA placement behavior for CPU threads and memory.",
		},
		{
			Name:        "LLAMA_ARG_DEVICE",
			Default:     "automatic",
			Description: "Provides a comma-separated list of accelerator devices for offloading; use none to disable offloading and use device names reported by llama-server.",
		},
		{
			Name:          "LLAMA_ARG_N_GPU_LAYERS",
			AllowedValues: []string{"auto", "all"},
			Default:       "auto",
			Description:   "Sets how many model layers are stored in accelerator memory; it also accepts a non-negative integer.",
		},
		{
			Name:          "LLAMA_ARG_SPLIT_MODE",
			AllowedValues: []string{"none", "layer", "row", "tensor"},
			Default:       "layer",
			Description:   "Controls how model weights and KV cache are split across multiple accelerators.",
		},
		{
			Name:        "LLAMA_ARG_TENSOR_SPLIT",
			Description: "Sets comma-separated per-device offload proportions, for example 3,1 for a 75/25 split.",
		},
		{
			Name:        "LLAMA_ARG_MAIN_GPU",
			Default:     "0",
			Description: "Selects the primary accelerator for single-device or row-split operations; use a non-negative device index.",
		},
		{
			Name:          "LLAMA_ARG_FIT",
			AllowedValues: []string{"on", "off"},
			Default:       "on",
			Description:   "Adjusts unset memory-related options to fit the model in available device memory.",
		},
		{
			Name:        "LLAMA_ARG_FIT_TARGET",
			Default:     "1024",
			Description: "Sets the target free-memory margin per device for automatic fitting, in MiB; a comma-separated list may target individual devices.",
		},
		{
			Name:        "LLAMA_ARG_FIT_CTX",
			Default:     "4096",
			Description: "Sets the minimum context size that automatic fitting may select, in tokens.",
		},
		{
			Name:        "LLAMA_ARG_TOP_K",
			Default:     "40",
			Description: "Sets the number of highest-probability tokens considered by top-k sampling; zero disables top-k filtering.",
		},
		{
			Name:        "LLAMA_ARG_N_PARALLEL",
			Default:     "-1",
			Description: "Sets the number of server slots used for parallel requests; -1 selects the automatic value.",
		},
		{
			Name:          "LLAMA_ARG_CONT_BATCHING",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Enables continuous dynamic batching so multiple requests can share model evaluation work.",
		},
		{
			Name:          "LLAMA_ARG_POOLING",
			AllowedValues: []string{"none", "mean", "cls", "last", "rank"},
			Default:       "model-dependent",
			Description:   "Selects the pooling strategy for embedding models.",
		},
		{
			Name:          "LLAMA_ARG_MMPROJ_AUTO",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Controls whether an available multimodal projector is used automatically, especially with Hugging Face models.",
		},
		{
			Name:          "LLAMA_ARG_MMPROJ_OFFLOAD",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Controls whether the multimodal projector is offloaded to accelerator memory.",
		},
		{
			Name:        "LLAMA_ARG_MMPROJ",
			Description: "Sets the local path to a multimodal projector file used by vision or other multimodal models.",
		},
		{
			Name:        "LLAMA_ARG_MMPROJ_URL",
			Description: "Sets a URL from which llama-server downloads a multimodal projector file.",
		},
		{
			Name:        "LLAMA_ARG_ALIAS",
			Description: "Sets comma-separated model aliases exposed through the API and used for model selection.",
		},
		{
			Name:        "LLAMA_ARG_TAGS",
			Description: "Sets comma-separated informational tags associated with the model; tags do not control routing.",
		},
		{
			Name:          "LLAMA_ARG_UI",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Enables or disables the llama-server Web UI.",
		},
		{
			Name:          "LLAMA_ARG_EMBEDDINGS",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Restricts the server to embedding use cases; enable it only with a dedicated embedding model.",
		},
		{
			Name:          "LLAMA_ARG_RERANKING",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Enables the native reranking endpoint.",
		},
		{
			Name:        "LLAMA_API_KEY",
			Description: "Sets one or more comma-separated API keys accepted by llama-server; keep this secret out of source code and logs.",
		},
		{
			Name:        "LLAMA_ARG_API_KEY_FILE",
			Description: "Sets a file containing accepted API keys, one per line; lines beginning with # are treated as comments.",
		},
		{
			Name:        "LLAMA_ARG_SSL_KEY_FILE",
			Description: "Sets the PEM-encoded private key used for HTTPS listener TLS.",
		},
		{
			Name:        "LLAMA_ARG_SSL_CERT_FILE",
			Description: "Sets the PEM-encoded certificate used for HTTPS listener TLS.",
		},
		{
			Name:        "LLAMA_ARG_CHAT_TEMPLATE_KWARGS",
			Description: "Sets a valid JSON object string with additional arguments passed to the chat-template parser.",
		},
		{
			Name:        "LLAMA_ARG_TIMEOUT",
			Default:     "3600",
			Description: "Sets the server read and write timeout in seconds.",
		},
		{
			Name:        "LLAMA_ARG_SSE_PING_INTERVAL",
			Default:     "30",
			Description: "Sets the interval between SSE keepalive pings in seconds; -1 disables pings.",
		},
		{
			Name:          "LLAMA_ARG_CACHE_PROMPT",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Enables prompt caching so repeated prefixes can reuse KV-cache state.",
		},
		{
			Name:        "LLAMA_ARG_CACHE_REUSE",
			Default:     "0",
			Description: "Sets the minimum chunk size for KV-shift prompt-cache reuse; zero disables this reuse threshold.",
		},
		{
			Name:          "LLAMA_ARG_ENDPOINT_METRICS",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Enables the Prometheus-compatible metrics endpoint.",
		},
		{
			Name:          "LLAMA_ARG_ENDPOINT_PROPS",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Enables changing global server properties through POST /props.",
		},
		{
			Name:          "LLAMA_ARG_ENDPOINT_SLOTS",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Controls whether the slots monitoring endpoint is exposed.",
		},
		{
			Name:        "LLAMA_ARG_MODELS_DIR",
			Description: "Sets the directory containing models available to the router server.",
		},
		{
			Name:        "LLAMA_ARG_MODELS_PRESET",
			Description: "Sets the path to the INI file containing router model presets.",
		},
		{
			Name:        "LLAMA_ARG_MODELS_MAX",
			Default:     "4",
			Description: "Limits the number of router models loaded simultaneously; zero means unlimited.",
		},
		{
			Name:          "LLAMA_ARG_MODELS_AUTOLOAD",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Controls whether the router automatically loads a model when a request selects it.",
		},
		{
			Name:          "LLAMA_ARG_JINJA",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Enables the Jinja template engine for chat message formatting.",
		},
		{
			Name:          "LLAMA_ARG_THINK",
			AllowedValues: []string{"none", "deepseek", "deepseek-legacy", "auto"},
			Default:       "auto",
			Description:   "Controls how reasoning tags are parsed and returned in chat responses.",
		},
		{
			Name:          "LLAMA_ARG_REASONING",
			AllowedValues: []string{"on", "off", "auto"},
			Default:       "auto",
			Description:   "Controls whether the chat template enables reasoning or thinking output.",
		},
		{
			Name:          "LLAMA_ARG_REASONING_EFFORT",
			AllowedValues: []string{"default", "minimal", "low", "medium", "high", "xhigh", "max"},
			Default:       "default",
			Description:   "Selects the reasoning effort level passed to the chat template.",
		},
		{
			Name:        "LLAMA_ARG_THINK_BUDGET",
			Default:     "-1",
			Description: "Sets the reasoning token budget; -1 is unrestricted, 0 ends thinking immediately, and positive values set a limit.",
		},
		{
			Name:        "LLAMA_ARG_THINK_BUDGET_MESSAGE",
			Description: "Sets the message injected before the end-of-thinking tag when the reasoning budget is exhausted.",
		},
		{
			Name:          "LLAMA_ARG_REASONING_PRESERVE",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Controls whether the reasoning trace is preserved across the full conversation history.",
		},
		{
			Name:        "LLAMA_ARG_CHAT_TEMPLATE",
			Description: "Selects a built-in chat template or supplies a supported template name; unset uses the model metadata template.",
		},
		{
			Name:        "LLAMA_ARG_CHAT_TEMPLATE_FILE",
			Description: "Sets a file containing a custom Jinja chat template.",
		},
		{
			Name:          "LLAMA_ARG_SWA_FULL",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Uses a full-size sliding-window attention cache instead of a reduced cache.",
		},
		{
			Name:          "LLAMA_ARG_PERF",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Enables internal llama.cpp performance timing fields in generation results.",
		},
		{
			Name:          "LLAMA_ARG_REPACK",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Enables weight repacking when the backend supports it.",
		},
		{
			Name:          "LLAMA_ARG_NO_HOST",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Bypasses the host buffer so the backend can use additional device buffers.",
		},
		{
			Name:        "LLAMA_ARG_DEFRAG_THOLD",
			Description: "Sets the deprecated KV-cache defragmentation threshold; prefer current cache controls when available.",
		},
		{
			Name:        "LLAMA_ARG_RPC",
			Description: "Sets a comma-separated list of remote RPC servers in host:port form for model computation.",
		},
		{
			Name:        "LLAMA_ARG_ROPE_FREQ_BASE",
			Description: "Overrides the RoPE base frequency used by NTK-aware scaling; unset loads it from model metadata.",
		},
		{
			Name:        "LLAMA_ARG_ROPE_FREQ_SCALE",
			Description: "Sets the RoPE frequency scaling factor; use a positive number when overriding model metadata.",
		},
		{
			Name:        "LLAMA_ARG_YARN_ORIG_CTX",
			Default:     "0",
			Description: "Sets the original model context size used by YaRN; zero uses the training context from model metadata.",
		},
		{
			Name:        "LLAMA_ARG_YARN_EXT_FACTOR",
			Default:     "-1.00",
			Description: "Sets the YaRN extrapolation mix factor; zero requests full interpolation.",
		},
		{
			Name:        "LLAMA_ARG_YARN_ATTN_FACTOR",
			Default:     "-1.00",
			Description: "Sets the YaRN attention scaling factor; negative keeps the backend default.",
		},
		{
			Name:        "LLAMA_ARG_YARN_BETA_SLOW",
			Default:     "-1.00",
			Description: "Sets the YaRN high-correction dimension or alpha parameter.",
		},
		{
			Name:        "LLAMA_ARG_YARN_BETA_FAST",
			Default:     "-1.00",
			Description: "Sets the YaRN low-correction dimension or beta parameter.",
		},
		{
			Name:        "LLAMA_ARG_DOCKER_REPO",
			Description: "Selects a Docker Hub model repository and optional quantization tag for model downloads.",
		},
		{
			Name:        "LLAMA_ARG_LOG_FILE",
			Description: "Writes server logs to the specified file instead of relying only on standard output.",
		},
		{
			Name:          "LLAMA_ARG_LOG_JSONL",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Writes logs as JSON Lines and disables colored log output.",
		},
		{
			Name:          "LLAMA_ARG_LOG_COLORS",
			AllowedValues: []string{"on", "off", "auto"},
			Default:       "auto",
			Description:   "Controls colored log output; auto enables colors when output is attached to a terminal.",
		},
		{
			Name:          "LLAMA_ARG_OFFLINE",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Forces model loading to use local cache and prevents network downloads.",
		},
		{
			Name:          "LLAMA_ARG_LOG_VERBOSITY",
			AllowedValues: []string{"0", "1", "2", "3", "4", "5"},
			Default:       "3",
			Description:   "Sets the log threshold: 0 generic, 1 error, 2 warning, 3 info, 4 trace, or 5 debug.",
		},
		{
			Name:          "LLAMA_ARG_LOG_PREFIX",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Adds a prefix to each log message when enabled.",
		},
		{
			Name:          "LLAMA_ARG_LOG_TIMESTAMPS",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Adds timestamps to each log message when enabled.",
		},
		{
			Name:        "LLAMA_ARG_THREADS_HTTP",
			Default:     "-1",
			Description: "Sets the number of threads used to process HTTP requests; -1 selects the server default.",
		},
		{
			Name:          "LLAMA_ARG_REUSE_PORT",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Allows multiple server processes to bind the same listening port when the operating system supports socket reuse.",
		},
		{
			Name:        "LLAMA_ARG_STATIC_PATH",
			Description: "Sets the directory served as static Web UI content.",
		},
		{
			Name:        "LLAMA_ARG_UI_CONFIG",
			Description: "Provides a JSON object string with default Web UI settings.",
		},
		{
			Name:        "LLAMA_ARG_UI_CONFIG_FILE",
			Description: "Sets a JSON file containing default Web UI settings.",
		},
		{
			Name:          "LLAMA_ARG_UI_MCP_PROXY",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Enables the experimental Web UI MCP CORS proxy; do not enable it in an untrusted environment.",
		},
		{
			Name:        "LLAMA_ARG_TOOLS",
			Default:     "none",
			Description: "Enables selected experimental built-in tools as a comma-separated list or all; do not enable tools in an untrusted environment.",
		},
		{
			Name:        "LLAMA_ARG_TOOLS_RUNTIME",
			Default:     "none",
			Description: "Selects the runtime for experimental built-in tools, such as docker:image, podman:image, an existing container, or ssh:target.",
		},
		{
			Name:        "LLAMA_ARG_MCP_SERVERS_CONFIG",
			Description: "Sets a JSON file containing experimental MCP server definitions; do not enable it in an untrusted environment.",
		},
		{
			Name:        "LLAMA_ARG_MCP_SERVERS_JSON",
			Description: "Provides experimental MCP server definitions as an inline JSON string; do not enable it in an untrusted environment.",
		},
		{
			Name:          "LLAMA_ARG_AGENT",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Enables the experimental CORS proxy and built-in agent tools; do not enable it in an untrusted environment.",
		},
		{
			Name:          "LLAMA_ARG_SKIP_CHAT_PARSING",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Forces a pure content parser so reasoning and tool calls remain in the message content.",
		},
		{
			Name:          "LLAMA_ARG_PREFILL_ASSISTANT",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Controls whether an assistant message at the end of the conversation is used as a generation prefix.",
		},
		{
			Name:        "MTMD_BACKEND_DEVICE",
			Default:     "follows LLAMA_ARG_DEVICE",
			Description: "Selects the device used for a multimodal projector; use none to keep it on the CPU.",
		},
		{
			Name:        "LLAMA_ARG_IMAGE_MIN_TOKENS",
			Description: "Sets the minimum token budget per image for vision models with dynamic resolution.",
		},
		{
			Name:        "LLAMA_ARG_IMAGE_MAX_TOKENS",
			Description: "Sets the maximum token budget per image for vision models with dynamic resolution.",
		},
		{
			Name:        "LLAMA_ARG_MTMD_BATCH_MAX_TOKENS",
			Default:     "1024",
			Description: "Sets the maximum number of image tokens processed in one multimodal batch.",
		},
		{
			Name:        "LLAMA_ARG_VIDEO_FPS",
			Default:     "4.0",
			Description: "Sets the target frame rate used when extracting frames from video input.",
		},
		{
			Name:        "LLAMA_ARG_VIDEO_TIMESTAMP_INTERVAL",
			Default:     "5000",
			Description: "Sets the interval between text timestamps extracted from video, in milliseconds.",
		},
		{
			Name:        "LLAMA_ARG_VIDEO_FFMPEG_DIR",
			Default:     "PATH lookup",
			Description: "Sets the directory containing ffmpeg and ffprobe used for video input; unset searches PATH.",
		},
		{
			Name:          "LLAMA_ARG_BACKEND_SAMPLING",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Enables experimental backend sampling instead of sampling entirely in the server process.",
		},
		{
			Name:        "LLAMA_ARG_KV_UNIFIED_PER_SLOT",
			Description: "Sets the context limit per parallel slot; when context size is unset, the shared KV pool is sized from this value and the slot count.",
		},
		{
			Name:        "LLAMA_ARG_CTX_CHECKPOINTS",
			Default:     "32",
			Description: "Sets the maximum number of context checkpoints created per slot.",
		},
		{
			Name:        "LLAMA_ARG_CHECKPOINT_MIN_SPACING_NT",
			Default:     "8192",
			Description: "Sets the minimum token spacing between context checkpoints; zero removes the minimum.",
		},
		{
			Name:        "LLAMA_ARG_CACHE_RAM",
			Default:     "8192",
			Description: "Sets the maximum prompt-cache size in MiB; zero disables the cache and -1 removes the size limit.",
		},
		{
			Name:          "LLAMA_ARG_KV_UNIFIED",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "automatic",
			Description:   "Controls whether all sequences share one unified KV buffer.",
		},
		{
			Name:          "LLAMA_ARG_CACHE_IDLE_SLOTS",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Saves idle slots to the prompt cache on a new task and clears them when using unified KV.",
		},
		{
			Name:          "LLAMA_ARG_CONTEXT_SHIFT",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "off",
			Description:   "Enables context shifting for effectively unlimited text generation when the context window is full.",
		},
		{
			Name:          "LLAMA_ARG_SPEC_DRAFT_CACHE_TYPE_K",
			AllowedValues: []string{"f32", "f16", "bf16", "q8_0", "q4_0", "q4_1", "iq4_nl", "q5_0", "q5_1"},
			Default:       "f16",
			Description:   "Selects the key-cache data type for a speculative draft model.",
		},
		{
			Name:          "LLAMA_ARG_SPEC_DRAFT_CACHE_TYPE_V",
			AllowedValues: []string{"f32", "f16", "bf16", "q8_0", "q4_0", "q4_1", "iq4_nl", "q5_0", "q5_1"},
			Default:       "f16",
			Description:   "Selects the value-cache data type for a speculative draft model.",
		},
		{
			Name:        "LLAMA_ARG_SPEC_DRAFT_HF_REPO",
			Description: "Selects the Hugging Face repository for a speculative draft model.",
		},
		{
			Name:        "LLAMA_ARG_SPEC_DRAFT_CPU_MOE",
			Description: "Keeps all Mixture-of-Experts weights of the speculative draft model on the CPU.",
		},
		{
			Name:        "LLAMA_ARG_SPEC_DRAFT_N_CPU_MOE",
			Description: "Keeps the first N Mixture-of-Experts layers of the speculative draft model on the CPU.",
		},
		{
			Name:        "LLAMA_ARG_SPEC_DRAFT_N_MAX",
			Default:     "3",
			Description: "Sets the maximum number of tokens proposed by speculative draft decoding.",
		},
		{
			Name:        "LLAMA_ARG_SPEC_DRAFT_N_MIN",
			Default:     "0",
			Description: "Sets the minimum number of draft tokens required before speculative decoding is used.",
		},
		{
			Name:        "LLAMA_ARG_SPEC_SYNTH_LEN",
			Description: "Sets the target mean synthetic acceptance length for speculative-decoding benchmarks.",
		},
		{
			Name:        "LLAMA_ARG_SPEC_SYNTH_RATES",
			Description: "Sets comma-separated per-position synthetic acceptance probabilities for speculative-decoding benchmarks.",
		},
		{
			Name:        "LLAMA_ARG_SPEC_DRAFT_P_SPLIT",
			Default:     "0.10",
			Description: "Sets the speculative-decoding split probability.",
		},
		{
			Name:        "LLAMA_ARG_SPEC_DRAFT_P_MIN",
			Default:     "0.00",
			Description: "Sets the minimum speculative-decoding probability used by greedy sampling.",
		},
		{
			Name:          "LLAMA_ARG_SPEC_DRAFT_BACKEND_SAMPLING",
			AllowedValues: []string{"on", "off", "true", "false", "1", "0"},
			Default:       "on",
			Description:   "Controls whether the speculative draft model performs sampling in the backend.",
		},
		{
			Name:          "LLAMA_ARG_N_GPU_LAYERS_DRAFT",
			AllowedValues: []string{"auto", "all"},
			Default:       "auto",
			Description:   "Sets the number of speculative draft model layers stored in accelerator memory; it also accepts a non-negative integer.",
		},
		{
			Name:        "LLAMA_ARG_SPEC_DRAFT_MODEL",
			Description: "Sets the local path to the speculative draft model.",
		},
		{
			Name:        "LLAMA_ARG_SPEC_TYPE",
			Default:     "none",
			Description: "Sets a comma-separated list of speculative decoding strategies, such as draft-simple, draft-eagle3, draft-mtp, or ngram-simple.",
		},
	}}
}
