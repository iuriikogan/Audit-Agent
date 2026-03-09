package store

import (
	"context"
	"time"
)

// ScanResult represents the metadata and findings of a CRA compliance scan.
type ScanResult struct {
	JobID       string     `json:"job_id"`
	Scope       string     `json:"scope"`
	Status      string     `json:"status"`
	Findings    []Finding  `json:"findings"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

// Finding represents a single security finding for a resource.
type Finding struct {
	ResourceName string `json:"resource_name"`
	Status       string `json:"status"`
	Details      string `json:"details"`
}

// Store defines the interface for persisting scan data.
type Store interface {
	CreateScan(ctx context.Context, jobID, scope string) error
	UpdateScanStatus(ctx context.Context, jobID, status string) error
	AddFinding(ctx context.Context, jobID string, f Finding) error
	GetScan(ctx context.Context, jobID string) (*ScanResult, error)
	// GetAllFindings retrieves all historical compliance findings, intended for dashboard display.
	GetAllFindings(ctx context.Context) ([]Finding, error)
	Close() error
}
