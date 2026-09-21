package llama_test

import (
	"testing"

	"go.osspkg.com/llm-client/llama"
)

func TestEnvironmentScheme(t *testing.T) {
	scheme := llama.EnvironmentScheme()
	if len(scheme.Variables) < 40 {
		t.Fatalf("environment variables = %d, want at least 40", len(scheme.Variables))
	}

	variables := make(map[string]struct{}, len(scheme.Variables))
	for _, variable := range scheme.Variables {
		if variable.Name == "" {
			t.Error("environment variable has an empty name")
		}
		if variable.Description == "" {
			t.Errorf("%s has an empty description", variable.Name)
		}
		if _, exists := variables[variable.Name]; exists {
			t.Errorf("environment variable %q is listed more than once", variable.Name)
		}
		variables[variable.Name] = struct{}{}
	}

	for _, name := range []string{"LLAMA_ARG_HOST", "LLAMA_ARG_MODEL", "LLAMA_ARG_CTX_SIZE", "LLAMA_API_KEY"} {
		if _, exists := variables[name]; !exists {
			t.Errorf("environment variable %q is missing", name)
		}
	}

	copyScheme := llama.EnvironmentScheme()
	copyScheme.Variables[0].Name = "CHANGED"
	if llama.EnvironmentScheme().Variables[0].Name == "CHANGED" {
		t.Fatal("environment scheme shares mutable state between calls")
	}
}
