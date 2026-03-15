package store

import (
	"testing"
)

func TestGCS_Paths(t *testing.T) {
	jobID := "test-job"
	resource := "res-1"

	gotMetadata := metadataPath(jobID)
	wantMetadata := "scans/test-job/metadata.json"
	if gotMetadata != wantMetadata {
		t.Errorf("metadataPath() = %q, want %q", gotMetadata, wantMetadata)
	}

	gotFinding := findingPath(jobID, resource)
	wantFinding := "scans/test-job/findings/res-1.json"
	if gotFinding != wantFinding {
		t.Errorf("findingPath() = %q, want %q", gotFinding, wantFinding)
	}
}

func TestGCSStore_Compilation(t *testing.T) {
	// Smoke test for the GCSStore struct and its methods.
	// We don't initialize a real client here.
	var _ = &GCSStore{}
}
