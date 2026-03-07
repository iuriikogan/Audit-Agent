package tools

import (
	"testing"

	"github.com/google/generative-ai-go/genai"
)

func TestToolDefinitions(t *testing.T) {
	toolSets := map[string][]*genai.Tool{
		"ScopeTools":             ScopeTools,
		"VulnTools":              VulnTools,
		"IngestionTools":         IngestionTools,
		"TaggingTools":           TaggingTools,
		"ComplianceTools":        ComplianceTools,
		"RegulatoryCheckerTools": RegulatoryCheckerTools,
	}

	for name, tools := range toolSets {
		t.Run(name, func(t *testing.T) {
			if len(tools) == 0 {
				t.Errorf("%s is empty", name)
			}
			for i, tool := range tools {
				if len(tool.FunctionDeclarations) == 0 {
					t.Errorf("%s[%d] has no function declarations", name, i)
				}
				for j, fn := range tool.FunctionDeclarations {
					if fn.Name == "" {
						t.Errorf("%s[%d].FunctionDeclarations[%d] has no name", name, i, j)
					}
					if fn.Description == "" {
						t.Errorf("%s[%d].FunctionDeclarations[%d] has no description", name, i, j)
					}
					if fn.Parameters == nil {
						t.Errorf("%s[%d].FunctionDeclarations[%d] has no parameters schema", name, i, j)
					}
				}
			}
		})
	}
}
