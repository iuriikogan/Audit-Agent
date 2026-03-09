# Technical Specification: 12-Factor & CRA Compliance Refactoring

## Overview
This document outlines the architectural changes required to align the Multi-Agent CRA application with 12-Factor App standards and to implement real-time security monitoring capabilities.

## 12-Factor App Adherence

### 1. Codebase
- **Current**: Single codebase.
- **Change**: Maintain single repository, but separate build targets or execution roles for the API Server and AI Worker.

### 2. Dependencies
- **Current**: Managed via `go.mod`.
- **Change**: Ensure explicit declaration across all modules "github.com/iuriikogan/multi-agent-cra" 

### 3. Config
- **Current**: Mixed approach (some via `config` package, some directly in `main.go`).
- **Change**: All configuration strictly pulled from environment variables using the os.env package

### 4. Backing Services
- **Current**: Pub/Sub (Messaging), GCS (Storage).
- **Change**: Treat all backing services as attached resources. Use stdout for cloud logging. Use GCS or CloudSQL, whichever is easier and cheaper

### 5. Build, Release, Run
- **Change**: Strictly separate stages using CI/CD.

### 6. Processes
- **Current**: `main.go` runs both the API server and the Pub/Sub subscriber (worker) in the same process.
- **Change**: Separate processes. Introduce a `ROLE` environment variable (`server` or `worker`) so the same image can be deployed differently, allowing independent scaling on Cloud Run.

### 7. Port Binding
- **Current**: Implemented.
- **Change**: Ensure the web server strictly binds to `$PORT`.

### 8. Concurrency
- **Change**: Scale out the worker processes via Cloud Run based on Pub/Sub queue depth.

### 9. Disposability
- **Change**: Implement graceful shutdown for both the server and worker processes using context cancellation on `SIGTERM`.

### 10. Dev/Prod Parity
- **Change**: Keep development, staging, and production as similar as possible using containerization.

### 11. Logs
- **Current**: `slog` is used.
- **Change**: Ensure all logs are written to `stdout` without local file routing. Forward to Cloud Logging.

### 12. Admin Processes
- **Change**: Run admin/management tasks as one-off processes.

## Security Engineer Monitoring (CRA Compliance)

### 1. Real-time Telemetry (Pub/Sub)
- **Change**: Emit granular events (`StepStarted`, `StepCompleted`, `AgentInsight`, `AuditTrail`) to a dedicated Pub/Sub topic (`compliance-telemetry`).

### 2. Live Dashboard Interface (Server)
- **Change**: The server will provide a `/api/stream` (SSE or WebSocket) endpoint to push real-time telemetry to the security dashboard.

### 3. Structured Audit Trail
- **Change**: AI agents must output structured JSON including reasoning and evidence.

## Proposed Changes to `cmd/server/main.go`
1. **Process Separation**: Read the `ROLE` env var. If `ROLE==worker`, start only the Pub/Sub subscriber. If `ROLE==server`, start only the HTTP server.
2. **Graceful Shutdown**: Trap `SIGINT` and `SIGTERM` and pass a cancellable context to the worker/server.
3. **Configuration**: Rely explicitly on `os.Getenv()` for top-level configs in `main.go` to adhere to factor III.
4. **Telemetry Setup**: Add a mock or skeleton for emitting compliance telemetry to BigQuery/Cloud Monitoring.

## Action Plan
1. **Step 1**: Refactor `cmd/server/main.go` to implement Process Separation (ROLE), Config via Env, and Graceful Shutdown.
2. **Step 2**: Add Realtime Monitoring logic (SSE endpoint for dashboard).
