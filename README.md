# Multi-Agent CRA System

This project implements a Multi-Agent System (MAS) for the EU Cyber Resilience Act (CRA) compliance using Go and Google's Gemini API.

## Architecture

The system uses a **Least Privilege** model where agents are split into "Specialists" (Doers) and "Checkers" (Verifiers).

### Agents

1.  **ScopeClassifier (Specialist)**
    *   **Role:** Determines if a product is Uncritical, Important Class I/II, or Critical.
    *   **Tools:** `get_product_specs`
    *   **Privilege:** Cannot access vulnerability DBs or write official reports.

2.  **RegulatoryAuditor (Checker)**
    *   **Role:** Verifies the `ScopeClassifier`'s output against the actual text of Regulation (EU) 2024/2847.
    *   **Tools:** `read_cra_regulation_text`
    *   **Privilege:** Read-only access to regulation text.

3.  **VulnWatchdog (Specialist)**
    *   **Role:** Checks specific components for vulnerabilities.
    *   **Tools:** `query_cve_database`

## Usage

1.  Set your API Key: `export GEMINI_API_KEY=your_key`
2.  Run the orchestrator: `go run cmd/main.go`
