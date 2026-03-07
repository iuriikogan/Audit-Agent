package tools

import "github.com/google/generative-ai-go/genai"

// Define tool sets to ensure Least Privilege.
// Only specific agents get access to specific tool lists.

// ScopeTools: Used by classifiers to understand product context.
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

// VulnTools: Used for checking vulnerabilities in software components.
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

// IngestionTools: Used to discover and list cloud assets or local files.
var IngestionTools = []*genai.Tool{
	{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "list_gcp_assets",
				Description: "Lists GCP assets within a given scope (project, folder, or organization).",
				Parameters: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"parent": {
							Type:        genai.TypeString,
							Description: "The parent resource name, e.g., 'projects/my-project', 'folders/123', 'organizations/456'.",
						},
						"asset_types": {
							Type:        genai.TypeArray,
							Items:       &genai.Schema{Type: genai.TypeString},
							Description: "Optional list of asset types to filter by, e.g., ['compute.googleapis.com/Instance', 'storage.googleapis.com/Bucket'].",
						},
					},
					Required: []string{"parent"},
				},
			},
		},
	},
}

// TaggingTools: Used to apply remediation or classification tags to resources.
var TaggingTools = []*genai.Tool{

	{

		FunctionDeclarations: []*genai.FunctionDeclaration{

			{

				Name: "apply_resource_tags",

				Description: "Applies a set of key-value tags to a specified cloud resource.",

				Parameters: &genai.Schema{

					Type: genai.TypeObject,

					Properties: map[string]*genai.Schema{

						"resource_id": {Type: genai.TypeString},

						"tags": {

							Type: genai.TypeObject,

							Description: "Key-value map of tags to apply.",
						},
					},

					Required: []string{"resource_id", "tags"},
				},
			},
		},
	},
}

// ComplianceTools: Used to generate formal compliance documents.
var ComplianceTools = []*genai.Tool{

	{

		FunctionDeclarations: []*genai.FunctionDeclaration{

			{

				Name: "generate_conformity_doc",

				Description: "Generates the official EU Declaration of Conformity PDF.",

				Parameters: &genai.Schema{

					Type: genai.TypeObject,

					Properties: map[string]*genai.Schema{

						"classification": {Type: genai.TypeString},

						"product_name": {Type: genai.TypeString},
					},

					Required: []string{"classification", "product_name"},
				},
			},
		},
	},
}

// RegulatoryCheckerTools: Used to validate compliance against official texts.
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

// VisualTools: Used to generate visual reports and dashboards.
var VisualTools = []*genai.Tool{
	{
		FunctionDeclarations: []*genai.FunctionDeclaration{
			{
				Name:        "generate_visual_dashboard",
				Description: "Generates a visual compliance dashboard image based on a descriptive prompt.",
				Parameters: &genai.Schema{
					Type: genai.TypeObject,
					Properties: map[string]*genai.Schema{
						"prompt": {
							Type:        genai.TypeString,
							Description: "A detailed description of the dashboard to generate.",
						},
						"filename": {
							Type:        genai.TypeString,
							Description: "The output filename (e.g., 'dashboard.png').",
						},
					},
					Required: []string{"prompt", "filename"},
				},
			},
		},
	},
}
