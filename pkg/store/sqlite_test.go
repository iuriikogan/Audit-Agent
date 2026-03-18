// Package store provides testing for the SQLiteStore implementation.
package store

import (
	"context"
	"testing"
)

// TestSQLiteStore verifies the primary lifecycle of a scan and its associated findings
// within a local SQLite in-memory database.
func TestSQLiteStore(t *testing.T) {
	ctx := context.Background()

	// Create a new in-memory SQLite store for isolated testing.
	s, err := NewSQLite(ctx, ":memory:")
	if err != nil {
		t.Fatalf("failed to create sqlite store: %v", err)
	}
	defer func() { _ = s.Close() }()

	jobID := "test-job-uuid"
	scope := "projects/test-project"
	reg := "DORA"

	// 1. Test CreateScan
	t.Run("CreateScan", func(t *testing.T) {
		if err := s.CreateScan(ctx, jobID, scope, reg); err != nil {
			t.Fatalf("failed to create scan: %v", err)
		}
	})

	// 2. Test AddFinding
	t.Run("AddFinding", func(t *testing.T) {
		finding := Finding{
			ResourceName: "test-compute-instance",
			Status:       "Compliant",
			Details:      map[string]string{"reason": "encryption enabled"},
			Regulation:   reg,
		}
		if err := s.AddFinding(ctx, jobID, finding); err != nil {
			t.Fatalf("failed to add finding: %v", err)
		}
	})

	// 3. Test UpdateScanStatus
	t.Run("UpdateScanStatus", func(t *testing.T) {
		if err := s.UpdateScanStatus(ctx, jobID, "completed"); err != nil {
			t.Fatalf("failed to update status: %v", err)
		}
	})

	// 4. Test GetScan and validation
	t.Run("GetScan", func(t *testing.T) {
		res, err := s.GetScan(ctx, jobID)
		if err != nil {
			t.Fatalf("failed to get scan: %v", err)
		}

		if res.JobID != jobID {
			t.Errorf("expected job ID %s, got %s", jobID, res.JobID)
		}
		if res.Status != "completed" {
			t.Errorf("expected status completed, got %s", res.Status)
		}
		if len(res.Findings) != 1 {
			t.Errorf("expected 1 finding, got %d", len(res.Findings))
		}
		if res.Findings[0].ResourceName != "test-compute-instance" {
			t.Errorf("unexpected finding resource name: %s", res.Findings[0].ResourceName)
		}
	})

	// 5. Test GetAllFindings
	t.Run("GetAllFindings", func(t *testing.T) {
		findings, err := s.GetAllFindings(ctx)
		if err != nil {
			t.Fatalf("failed to get all findings: %v", err)
		}
		if len(findings) == 0 {
			t.Error("expected findings to be returned, got 0")
		}
	})
}
