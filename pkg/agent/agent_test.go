package agent

import (
	"context"
	"strings"
	"testing"
)

// Since we cannot easily test the specific GeminiAgent struct without a live client,
// We will create a MockAgent here that can be used by other packages (like workflow)
// if they import it. But usually mocks belong in the test package of the consumer
// or a separate mocks package.
//
// For this file, let's just ensure that GeminiAgent implements the Agent interface.
var _ Agent = (*GeminiAgent)(nil)

func TestGeminiAgent_Name_Role(t *testing.T) {
	// We can construct a struct directly for testing simple getters
	// since we are in the same package.
	a := &GeminiAgent{
		name: "test-agent",
		role: "tester",
	}

	if a.Name() != "test-agent" {
		t.Errorf("expected name 'test-agent', got %s", a.Name())
	}
	if a.Role() != "tester" {
		t.Errorf("expected role 'tester', got %s", a.Role())
	}
}

func TestExecuteMockTool(t *testing.T) {
	// executeMockTool is private, so we can test it here.
	a := &GeminiAgent{}
	ctx := context.Background()

	tests := []struct {
		name     string
		toolName string
		args     map[string]interface{}
		want     string
		contains string // check if output contains string (for variable outputs)
	}{
		{
			name:     "get_product_specs",
			toolName: "get_product_specs",
			args:     map[string]interface{}{"product_id": "123"},
			want:     "Technical specs for 123: Processor X1, 8GB RAM, Secure Boot enabled.",
		},
		{
			name:     "read_cra_regulation_text",
			toolName: "read_cra_regulation_text",
			args:     map[string]interface{}{"article_number": "10"},
			contains: "Article X",
		},
		{
			name:     "ingest_file_system",
			toolName: "ingest_file_system",
			args:     map[string]interface{}{"path": "/tmp/project"},
			want:     "Found: config.yaml, main.go, README.md",
		},
		{
			name:     "ingest_git_repo",
			toolName: "ingest_git_repo",
			args:     map[string]interface{}{"repo_url": "https://github.com/example/repo"},
			contains: "Cloned https://github.com/example/repo",
		},
		{
			name:     "apply_resource_tags",
			toolName: "apply_resource_tags",
			args:     map[string]interface{}{"resource_id": "res-1", "tags": map[string]string{"env": "prod"}},
			contains: "Tags applied successfully",
		},
		{
			name:     "generate_conformity_doc",
			toolName: "generate_conformity_doc",
			args:     map[string]interface{}{"product_name": "Widget", "classification": "Class I"},
			want:     "Generated EU Declaration of Conformity for Widget (Class: Class I)",
		},
		{
			name:     "query_cve_database",
			toolName: "query_cve_database",
			args:     map[string]interface{}{"component": "LibX", "version": "1.0"},
			contains: "No CRITICAL vulnerabilities found",
		},
		{
			name:     "default",
			toolName: "unknown_tool",
			args:     nil,
			want:     "Tool executed successfully.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := a.executeMockTool(ctx, tt.toolName, tt.args)
			if tt.contains != "" {
				if !strings.Contains(got, tt.contains) {
					t.Errorf("executeMockTool() = %q, expected to contain %q", got, tt.contains)
				}
			} else {
				if got != tt.want {
					t.Errorf("executeMockTool() = %q, want %q", got, tt.want)
				}
			}
		})
	}
}
