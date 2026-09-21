package ollama_test

import (
	"testing"

	"go.osspkg.com/llm-client/ollama"
)

func TestEnvironmentScheme(t *testing.T) {
	scheme := ollama.EnvironmentScheme()
	if len(scheme.Variables) < 20 {
		t.Fatalf("environment variables = %d, want at least 20", len(scheme.Variables))
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

	for _, name := range []string{"OLLAMA_HOST", "OLLAMA_MODELS", "OLLAMA_NUM_PARALLEL"} {
		if _, exists := variables[name]; !exists {
			t.Errorf("environment variable %q is missing", name)
		}
	}

	copyScheme := ollama.EnvironmentScheme()
	copyScheme.Variables[0].Name = "CHANGED"
	if ollama.EnvironmentScheme().Variables[0].Name == "CHANGED" {
		t.Fatal("environment scheme shares mutable state between calls")
	}
}
