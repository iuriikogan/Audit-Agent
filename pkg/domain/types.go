package domain

// Product represents the item being assessed.
type Product struct {
	ID        string
	Name      string
	Component string // For demo purposes, assumes single main component
	Version   string
}

// AssessmentResult holds the outcome of the multi-agent analysis.
type AssessmentResult struct {
	ProductID      string
	Classification string
	AuditStatus    string // "Verified" or "Rejected"
	VulnReport     string
	Error          error
}
