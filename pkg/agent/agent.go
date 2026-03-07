package agent

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
)

// Agent represents the behavioral contract for any AI agent in the system.
type Agent interface {
	Name() string
	Role() string
	Chat(ctx context.Context, input string) (string, error)
}

// GeminiAgent is the concrete implementation using Google's GenAI SDK.
type GeminiAgent struct {
	name              string
	role              string
	model             *genai.GenerativeModel
	systemInstruction string
}

// Option defines a functional option for configuring an agent.
type Option func(*GeminiAgent)

// WithTools adds tools to the agent configuration.
func WithTools(tools ...*genai.Tool) Option {
	return func(a *GeminiAgent) {
		if len(tools) > 0 {
			a.model.Tools = append(a.model.Tools, tools...)
		}
	}
}

// WithSystemInstruction sets the system instruction.
func WithSystemInstruction(instruction string) Option {
	return func(a *GeminiAgent) {
		a.systemInstruction = instruction
		a.model.SystemInstruction = &genai.Content{
			Parts: []genai.Part{genai.Text(instruction)},
		}
	}
}

// New creates a new GeminiAgent with the given options.
func New(client *genai.Client, name, role, modelName string, opts ...Option) *GeminiAgent {
	if modelName == "" {
		modelName = "gemini-3.1-flash-lite"
	}
	model := client.GenerativeModel(modelName)
	model.SetTemperature(0.2) // Low temperature for deterministic/regulatory tasks

	agent := &GeminiAgent{
		name:  name,
		role:  role,
		model: model,
	}

	for _, opt := range opts {
		opt(agent)
	}

	return agent
}

func (a *GeminiAgent) Name() string { return a.name }
func (a *GeminiAgent) Role() string { return a.role }

// Chat executes a single interaction. 
// It handles potential tool calls by the model by executing them and providing results back.
func (a *GeminiAgent) Chat(ctx context.Context, input string) (string, error) {
	session := a.model.StartChat()
	
	msg := genai.Text(input)
	for {
		resp, err := session.SendMessage(ctx, msg)
		if err != nil {
			return "", fmt.Errorf("[%s] messaging failed: %w", a.name, err)
		}

		if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
			return "", fmt.Errorf("[%s] empty response", a.name)
		}

		// Check for tool calls in the response parts
		var functionCalls []genai.FunctionCall
		for _, part := range resp.Candidates[0].Content.Parts {
			if fc, ok := part.(genai.FunctionCall); ok {
				functionCalls = append(functionCalls, fc)
			}
		}

		// If no tool calls, find the first text part and return it
		if len(functionCalls) == 0 {
			for _, part := range resp.Candidates[0].Content.Parts {
				if t, ok := part.(genai.Text); ok {
					return string(t), nil
				}
			}
			return "", fmt.Errorf("[%s] unhandled response format (no text or tools)", a.name)
		}

		// Execute tool calls and prepare response parts
		var responseParts []genai.Part
		for _, fc := range functionCalls {
			// In a production system, you'd route these to actual function implementations.
			// For this demo, we return mock results for the defined tools.
			result := a.executeMockTool(fc.Name, fc.Args)
			responseParts = append(responseParts, genai.FunctionResponse{
				Name:     fc.Name,
				Response: map[string]any{"result": result},
			})
		}
		
		// Send the function responses back to the model
		msg = genai.FunctionResponseList(responseParts)
	}
}

// executeMockTool provides placeholder logic for the demo tools.
func (a *GeminiAgent) executeMockTool(name string, args map[string]any) string {
	switch name {
	case "get_product_specs":
		return fmt.Sprintf("Technical specs for %v: Processor X1, 8GB RAM, Secure Boot enabled.", args["product_id"])
	case "query_cve_database":
		return fmt.Sprintf("No CRITICAL vulnerabilities found for %s %s. 2 LOW found in dependencies.", args["component"], args["version"])
	case "read_cra_regulation_text":
		return "Article X: Products with digital elements shall be designed, developed and produced such that they ensure an appropriate level of cybersecurity."
	default:
		return "Tool executed successfully."
	}
}
