package config

import (
	"os"
	"testing"
)

func TestLoad(t *testing.T) {
	// Set environment variables for testing
	_ = os.Setenv("PROJECT_ID", "test-project")
	_ = os.Setenv("REGION", "us-west1")
	_ = os.Setenv("LOG_LEVEL", "DEBUG")
	_ = os.Setenv("GEMINI_API_KEY", "test-key")
	_ = os.Setenv("PUBSUB_TOPIC_SCAN_REQUESTS", "test-topic")
	_ = os.Setenv("PORT", "9090")
	defer func() {
		_ = os.Unsetenv("PROJECT_ID")
		_ = os.Unsetenv("REGION")
		_ = os.Unsetenv("LOG_LEVEL")
		_ = os.Unsetenv("GEMINI_API_KEY")
		_ = os.Unsetenv("PUBSUB_TOPIC_SCAN_REQUESTS")
		_ = os.Unsetenv("PORT")
	}()

	cfg := Load()

	if cfg.ProjectID != "test-project" {
		t.Errorf("expected ProjectID 'test-project', got %s", cfg.ProjectID)
	}
	if cfg.Region != "us-west1" {
		t.Errorf("expected Region 'us-west1', got %s", cfg.Region)
	}
	if cfg.LogLevel != "DEBUG" {
		t.Errorf("expected LogLevel 'DEBUG', got %s", cfg.LogLevel)
	}
	if cfg.PubSub.TopicScanRequests != "test-topic" {
		t.Errorf("expected PubSub topic 'test-topic', got %s", cfg.PubSub.TopicScanRequests)
	}
	if cfg.Server.Port != "9090" {
		t.Errorf("expected Server port '9090', got %s", cfg.Server.Port)
	}
}

func TestLoad_Defaults(t *testing.T) {
	// Ensure environment is clean
	_ = os.Unsetenv("LOG_LEVEL")
	_ = os.Unsetenv("PORT")

	cfg := Load()

	if cfg.LogLevel != "INFO" {
		t.Errorf("expected default LogLevel 'INFO', got %s", cfg.LogLevel)
	}
	if cfg.Server.Port != "8080" {
		t.Errorf("expected default Server port '8080', got %s", cfg.Server.Port)
	}
}
