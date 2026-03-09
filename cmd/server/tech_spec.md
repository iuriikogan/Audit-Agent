# Technical Specification: 12-Factor & CRA Compliance Refactoring

## Purpose
Refactor the `multi-agent-cra` project to strictly adhere to 12-factor app standards, implement an event-driven agent architecture using GCP MCP servers for real-time Cyber Resilience Act (CRA) compliance monitoring, and build a frontend dashboard for real-time security engineer oversight.

## Prerequisites
- Go 1.22+ and `golangci-lint` installed.
- Node.js and `npm` installed (for the `web/` frontend).
- GCP access (Cloud Monitoring, BigQuery, Pub/Sub, Cloud Run, GCP MCP Servers).

## Context
The current implementation combines server and worker roles, relies on local file/GCS state, coordinates agents via in-memory Go channels, and uses `gcloud` CLI calls. To satisfy 12-factor standards and enable security engineers to monitor real-time CRA compliance, the architecture must transition to isolated Cloud Run processes, Pub/Sub event-driven pipelines, GCP MCP servers, and GCS/CloudSQL for state. Additionally, a frontend dashboard is required to surface these real-time streams and provide control over AI agent jobs.

## Changes

1. **Config & Logs (12-Factor Adherence)**: 
   Update `pkg/config/config.go` and `cmd/server/main.go` to strictly use environment variables for configuration. Update `pkg/logger` to ensure `slog` exclusively writes to `Stdout` without local file routing.
2. **Process Separation & Disposability**: 
   Refactor `cmd/server/main.go` into separate execution paths based on a `ROLE` env var (`server` or `worker`). Implement context cancellation on `SIGTERM`/`SIGINT` for graceful shutdown. 
3. **State Management Refactoring**: 
   Replace GCS/local JSON file state management in `pkg/store/gcs.go` with interfaces for cloudsql.
4. **Event-Driven Agent Architecture**: 
   Refactor `pkg/workflow` (Coordinator) to replace in-memory Go channels with Pub/Sub message publishing and subscribing. Each agent will process an event and publish the result back to Pub/Sub.
5. **GCP MCP Servers Integration**: 
   Replace shell-out `gcloud` executions in `pkg/tools/executor.go` with GCP Model Context Protocol (MCP) server integrations to bridge security data and AI agents in real time.
6. **Real-Time CRA Monitoring & API**: 
   Implement real-time metrics emission to Cloud Monitoring. Add a `/api/stream` Server-Sent Events (SSE) endpoint to the server (`cmd/server/main.go`) to stream structured compliance audit trails to the frontend.
7. **Dashboard UI & Frontend**: 
   Update the Next.js application in the `web/` directory to consume the `/api/stream` SSE endpoint. Implement a live logging interface for CRA compliance status and a chat box integration allowing security engineers to trigger new agent jobs.

## Out of Scope
- Provisioning actual GCP resources via Terraform (infrastructure as code is assumed to be handled separately or pre-existing). 

## Verification
Every check below is mandatory. Do not skip any.
- Every Go change MUST be followed by a successful `golangci-lint run ./...`.
- Go unit tests must pass (`go test ./... -count=1`).
- Configurations must strictly use environment variables and not hardcode any credentials.
- The `ROLE` environment variable must determine the backend process mode.
- The Next.js frontend must build successfully (`npm run build` in the `web/` directory).

## Branch
`refactor/12-factor-cra-compliance`

## Provenance
`specs/provenance/refactor-cra-compliance.provenance.md`
