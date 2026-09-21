/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ollama

import "go.osspkg.com/llm-client/pkg/environment"

// EnvironmentScheme describes the Ollama server environment variables
// relevant to installation, networking, model storage, scheduling, and
// runtime behavior. It is independent of a Client and does not inspect the
// process environment.
func EnvironmentScheme() environment.Scheme {
	return environment.Scheme{Variables: []environment.Variable{
		{
			Name:          "OLLAMA_DEBUG",
			AllowedValues: []string{"0", "1", "2", "false", "true"},
			Default:       "0",
			Description:   "Controls server logging: 0 or false for information, 1 or true for debug, and 2 for trace output.",
		},
		{
			Name:          "OLLAMA_DEBUG_LOG_REQUESTS",
			AllowedValues: []string{"false", "true"},
			Default:       "false",
			Description:   "Writes inference request bodies and replay commands to a temporary directory for debugging; avoid enabling it for sensitive prompts.",
		},
		{
			Name:          "OLLAMA_GO_TEMPLATE",
			AllowedValues: []string{"false", "true"},
			Default:       "true",
			Description:   "Enables rendering of a model's Modelfile TEMPLATE when Ollama's Go template renderer is available.",
		},
		{
			Name:          "OLLAMA_FLASH_ATTENTION",
			AllowedValues: []string{"false", "true"},
			Default:       "false",
			Description:   "Enables Flash Attention when the selected backend and model support it.",
		},
		{
			Name:        "OLLAMA_KV_CACHE_TYPE",
			Default:     "f16",
			Description: "Selects the quantization type used by the key/value cache; the accepted types depend on the active llama.cpp backend.",
		},
		{
			Name:        "OLLAMA_GPU_OVERHEAD",
			Default:     "0",
			Description: "Reserves this many bytes of VRAM per GPU so Ollama does not allocate that memory to a model; use a non-negative integer.",
		},
		{
			Name:          "OLLAMA_IGPU_ENABLE",
			AllowedValues: []string{"false", "true"},
			Default:       "false",
			Description:   "Allows Ollama to select integrated GPUs during backend discovery.",
		},
		{
			Name:          "LLAMA_ARG_FIT",
			AllowedValues: []string{"on", "off"},
			Default:       "on",
			Description:   "Enables llama.cpp automatic adjustment of unset memory options so the model fits available device memory.",
		},
		{
			Name:        "LLAMA_ARG_FIT_TARGET",
			Default:     "1024",
			Description: "Sets the target free-memory margin per device for llama.cpp automatic fitting, in MiB; use a non-negative integer.",
		},
		{
			Name:        "OLLAMA_HOST",
			Default:     "127.0.0.1:11434",
			Description: "Sets the scheme, host, port, and optional path where the Ollama server listens, for example http://0.0.0.0:11434.",
		},
		{
			Name:        "OLLAMA_KEEP_ALIVE",
			Default:     "5m",
			Description: "Controls how long an idle model remains loaded; accept a Go duration such as 30s or 5m, 0 for no keep-alive, or a negative value for indefinitely.",
		},
		{
			Name:        "OLLAMA_LLM_LIBRARY",
			Description: "Selects an LLM backend explicitly and bypasses automatic backend detection; leave unset to autodetect.",
		},
		{
			Name:        "OLLAMA_LOAD_TIMEOUT",
			Default:     "5m",
			Description: "Sets the maximum time a model load may stall; accept a Go duration or seconds as an integer, with zero or a negative value meaning unlimited.",
		},
		{
			Name:        "OLLAMA_MAX_LOADED_MODELS",
			Default:     "0",
			Description: "Limits the number of concurrently loaded models per GPU; zero selects Ollama's automatic limit based on available hardware.",
		},
		{
			Name:        "OLLAMA_MAX_TRANSFER_STREAMS",
			Default:     "4",
			Description: "Limits concurrent transfer streams for safetensors pulls and pushes; use a positive integer, and note that it does not affect legacy GGUF transfers.",
		},
		{
			Name:        "OLLAMA_MAX_QUEUE",
			Default:     "512",
			Description: "Limits the number of requests queued while the server is busy; additional requests are rejected when this positive integer is reached.",
		},
		{
			Name:        "OLLAMA_MODELS",
			Default:     "$HOME/.ollama/models",
			Description: "Sets the directory where Ollama stores model blobs and manifests; use an absolute path writable by the server process.",
		},
		{
			Name:          "OLLAMA_NO_CLOUD",
			AllowedValues: []string{"false", "true"},
			Default:       "false",
			Description:   "Disables Ollama cloud features such as remote inference and web search when enabled.",
		},
		{
			Name:          "OLLAMA_NOHISTORY",
			AllowedValues: []string{"false", "true"},
			Default:       "false",
			Description:   "Prevents the interactive CLI from preserving readline history.",
		},
		{
			Name:          "OLLAMA_NOPRUNE",
			AllowedValues: []string{"false", "true"},
			Default:       "false",
			Description:   "Prevents pruning of unused model blobs during server startup.",
		},
		{
			Name:        "OLLAMA_NUM_PARALLEL",
			Default:     "1",
			Description: "Sets the maximum number of requests processed concurrently by each loaded model; use a positive integer because memory use grows with this value and context length.",
		},
		{
			Name:        "OLLAMA_ORIGINS",
			Default:     "localhost,127.0.0.1,0.0.0.0,app://*,file://*,tauri://*,vscode-webview://*,vscode-file://*",
			Description: "Provides a comma-separated allowlist of browser origins that may call the server; keep it narrow when exposing Ollama beyond localhost.",
		},
		{
			Name:          "OLLAMA_SCHED_SPREAD",
			AllowedValues: []string{"false", "true"},
			Default:       "false",
			Description:   "Forces a model to be scheduled across all available GPUs instead of preferring a single GPU when it fits.",
		},
		{
			Name:        "OLLAMA_CONTEXT_LENGTH",
			Default:     "0",
			Description: "Sets the default context length in tokens; zero lets Ollama choose based on available VRAM, otherwise use a positive integer.",
		},
		{
			Name:          "OLLAMA_CREATE_REMOTE",
			AllowedValues: []string{"false", "true"},
			Default:       "false",
			Description:   "Forces model creation through the Ollama server API even when the server is running locally.",
		},
		{
			Name:        "OLLAMA_EDITOR",
			Description: "Sets the executable used by the interactive CLI for prompt editing; leave unset to use the platform default.",
		},
		{
			Name:        "OLLAMA_REMOTES",
			Default:     "ollama.com",
			Description: "Provides a comma-separated allowlist of hosts from which remote models may be used.",
		},
		{
			Name:          "OLLAMA_AUTH",
			AllowedValues: []string{"false", "true"},
			Default:       "false",
			Description:   "Enables authentication between the Ollama client and server when supported by the selected Ollama build.",
		},
		{
			Name:          "OLLAMA_VULKAN",
			AllowedValues: []string{"false", "true"},
			Default:       "true",
			Description:   "Enables Vulkan backend discovery on platforms where the Ollama build supports it.",
		},
		{
			Name:        "CUDA_VISIBLE_DEVICES",
			Description: "Restricts which NVIDIA devices are visible to CUDA; use a comma-separated list of device identifiers.",
		},
		{
			Name:        "HIP_VISIBLE_DEVICES",
			Description: "Restricts which AMD devices are visible to HIP by numeric device identifier.",
		},
		{
			Name:        "ROCR_VISIBLE_DEVICES",
			Description: "Restricts which AMD devices are visible to ROCr by UUID or numeric device identifier.",
		},
		{
			Name:        "GGML_VK_VISIBLE_DEVICES",
			Description: "Restricts which Vulkan devices are visible by numeric device identifier.",
		},
		{
			Name:        "GPU_DEVICE_ORDINAL",
			Description: "Selects visible AMD devices by numeric ordinal for the platform GPU runtime.",
		},
		{
			Name:        "HSA_OVERRIDE_GFX_VERSION",
			Description: "Overrides the GFX architecture reported for detected AMD GPUs; use only when the platform runtime requires it.",
		},
		{
			Name:        "HTTP_PROXY",
			Description: "Sets the HTTP proxy used by Ollama for network operations such as remote model transfers; leave unset to connect directly.",
		},
		{
			Name:        "HTTPS_PROXY",
			Description: "Sets the HTTPS proxy used by Ollama for secure network operations such as remote model transfers; leave unset to connect directly.",
		},
		{
			Name:        "NO_PROXY",
			Description: "Provides a comma-separated list of hosts and address patterns that must bypass HTTP and HTTPS proxies.",
		},
	}}
}
