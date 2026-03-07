package tools

import "github.com/google/generative-ai-go/genai"

// Define tool sets to ensure Least Privilege. 
// Only specific agents get access to specific tool lists.

var ScopeTools = []*genai.Tool{
	{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "get_product_specs",
				Description: "Retrieves technical specifications of a product ID.",
				Parameters: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"product_id": {Type: genai.TypeString},
					},
					Required: []string{"product_id"},
				},
			},
		},
	},
}

var VulnTools = []*genai.Tool{
	{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "query_cve_database",
				Description: "Checks public CVE databases for a given component and version.",
				Parameters: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"component": {Type: genai.TypeString},
						"version":   {Type: genai.TypeString},
					},
					Required: []string{"component", "version"},
				},
			},
		},
	},
}

var ComplianceTools = []*genai.Tool{
	{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "generate_conformity_doc",
				Description: "Generates the official EU Declaration of Conformity PDF.",
				Parameters: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"classification": {Type: genai.TypeString},
						"product_name":   {Type: genai.TypeString},
					},
					Required: []string{"classification", "product_name"},
				},
			},
		},
	},
}

// Checker agents often don't need external tools, just logic, 
// OR they need read-only access to a "Source of Truth" (like the CRA PDF text).
var RegulatoryCheckerTools = []*genai.Tool{
	{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "read_cra_regulation_text",
				Description: "Reads specific articles from the EU Cyber Resilience Act.",
				Parameters: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"article_number": {Type: genai.TypeString},
					},
					Required: []string{"article_number"},
				},
			},
		},
	},
}
