# Multi-Agent CRA Security Platform

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A scalable, event-driven multi-agent system designed to assess Google Cloud infrastructure against the EU Cyber Resilience Act (CRA). The goal is to provide Security Engineers with a real-time, dashboard-driven tool to monitor, audit, and enforce CRA compliance across their GCP estate.

## 🚀 Key Features

*   **Autonomous Agents:** Specialized AI agents for Discovery (Aggregator), Modeling, Validation, Review, and Tagging.
*   **Real-time Dashboard:** A Next.js frontend embedded in the Go binary featuring live Server-Sent Events (SSE) log streaming and interactive compliance charts.
*   **12-Factor Architecture:** Built for the cloud. Configuration is strictly environment-variable driven. The application scales independently by setting the `ROLE` variable (`server`, `worker`, or `all`).
*   **Event-Driven:** Decoupled architecture using Google Cloud Pub/Sub for resilient, multi-stage agent pipelines.
*   **Flexible Storage:** Choose between robust **Cloud SQL** (PostgreSQL) for production or lightweight in-memory **SQLite** for zero-dependency local development.
*   **AI-Powered:** Leverages Gemini 1.5/3.0 for deep reasoning and compliance mapping via the native Go SDK.

## 🏗️ System Architecture & Data Flow

The system uses a strictly decoupled producer-consumer model:

1.  **Frontend (UI):** Users interact with the embedded React dashboard to initiate scans or view historical CRA findings.
2.  **API Server (`ROLE=server`):** Receives HTTP scan requests, publishes them to Pub/Sub, and serves historical data from the database. It also maintains long-lived SSE connections to broadcast internal monitoring events to the browser.
3.  **Message Broker (Pub/Sub):** Manages discrete topics for every stage of the agent pipeline (`scan-requests` -> `aggregator` -> `modeler` -> `validator` -> `reviewer` -> `tagger`).
4.  **Worker Fleet (`ROLE=worker`):** Stateless background processes that consume Pub/Sub messages, execute Gemini agent logic, interact with GCP APIs (like Cloud Asset Inventory), and write findings to the database.
5.  **State Store:** 
    *   **Cloud SQL (Production):** Persistent storage of scan metadata and compliance findings.
    *   **SQLite (Local):** In-memory ephemeral storage for rapid testing.

### Security Controls
*   **Least Privilege:** Workers operate using dedicated Google Service Accounts with minimal permissions required for Asset Inventory and Pub/Sub.
*   **No Hardcoded Secrets:** API keys and Database URLs are injected securely at runtime via environment variables (Factor III).
*   **Network Isolation:** Cloud SQL instances should be deployed with private IPs. The `cra-worker` does not expose any inbound ports.

## 📂 Project Structure

```
├── cmd/
│   ├── server/      # Unified Entrypoint (API + UI + Config Routing)
│   └── worker/      # Legacy entrypoint (now handled by cmd/server via ROLE)
├── pkg/
│   ├── agent/       # Gemini AI Agent logic
│   ├── config/      # Centralized 12-factor configuration
│   ├── queue/       # Pub/Sub client implementations
│   ├── store/       # Cloud SQL and SQLite implementations
│   ├── tools/       # GCP SDK and LLM tool definitions
│   └── workflow/    # Pub/Sub pipeline orchestrator
├── web/             # Next.js Frontend Dashboard (compiled into Go binary)
└── terraform/       # IaC definitions for GCP deployment
```

## 🛠️ Deployment Instructions

### Local Development (Zero Dependencies)

The easiest way to run the platform locally is using the in-memory SQLite database and running both the server and worker in a single process.

1.  **Set Environment Variables**:
    ```bash
    export GEMINI_API_KEY="your_actual_api_key_here"
    export PROJECT_ID="your-gcp-project-id"
    export ROLE="all" # Runs both API and background workers
    export DATABASE_TYPE="SQLITE_MEM" # Uses in-memory DB
    # Ensure you have valid GCP credentials for Pub/Sub and Asset Inventory:
    # gcloud auth application-default login
    ```

2.  **Run the Application**:
    ```bash
    go run ./cmd/server
    ```
    *   **Dashboard & API:** http://localhost:8080

### Production Deployment (Cloud Run & Cloud SQL)

For production, deploy the `server` and `worker` as separate Cloud Run services to scale them independently.

1.  **Database Setup:** Provision a Cloud SQL (PostgreSQL) instance and obtain the connection string.
2.  **Pub/Sub Setup:** Ensure all topics and subscriptions defined in `pkg/config/config.go` exist in your GCP project.
3.  **Deploy API Server**:
    ```bash
    gcloud run deploy cra-server \
      --source . \
      --set-env-vars="ROLE=server,DATABASE_TYPE=CLOUD_SQL,DATABASE_URL=postgres://user:pass@host/db" \
      --set-secrets="GEMINI_API_KEY=gemini-api-key:latest" \
      --allow-unauthenticated
    ```
4.  **Deploy Worker**:
    ```bash
    gcloud run deploy cra-worker \
      --source . \
      --set-env-vars="ROLE=worker,DATABASE_TYPE=CLOUD_SQL,DATABASE_URL=postgres://user:pass@host/db" \
      --set-secrets="GEMINI_API_KEY=gemini-api-key:latest" \
      --no-allow-unauthenticated
    ```
