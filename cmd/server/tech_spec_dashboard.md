# Technical Specification: CRA Compliance Dashboard

## Purpose
This specification details the implementation of a new CRA (Compliance and Risk Assessment) dashboard within the existing multi-agent application. This dashboard will serve as the primary interface for security engineers to monitor, analyze, and share CRA findings.

## Prerequisites
- The existing Go backend with Server-Sent Events (SSE) at `/api/stream` must be running.
- The Next.js frontend framework in the `web/` directory must be used.
- The `spec-driven-development` skill must be loaded to ensure process adherence.

## Context
Security engineers require a centralized and interactive view to understand the organization's security posture against CRA standards. The current application infrastructure provides real-time agent data but lacks a user-facing dashboard to visualize and interact with these findings. This dashboard will bridge that gap.

## Changes

1.  **Backend API Endpoint (`/api/findings`)**:
    *   Create a new HTTP GET endpoint at `/api/findings`.
    *   This endpoint will query the configured data store (`CloudSQL` or `SQLite`) to fetch all historical CRA findings.
    *   It will return a JSON array of findings, structured to be easily consumable by the frontend. The structure should include `job_id`, `resource_name`, `status`, `details`, and `timestamp`.

2.  **Frontend: Main Dashboard View**:
    *   The main page (`web/pages/index.tsx`) will be converted into a tabbed interface.
    *   The first tab, "CRA Dashboard," will be the default view.
    *   This dashboard will fetch data from the `/api/findings` endpoint on initial load.
    *   It will display findings in a filterable and sortable table.
    *   Add controls (e.g., dropdowns) to filter findings by GCP hierarchy: **Organization**, **Folder**, and **Project**. This will require parsing the `resource_name` from the findings.
    *   The dashboard will include data visualization elements, such as a pie chart or bar chart, summarizing the compliance status (e.g., "compliant," "non-compliant," "remediated").

3.  **Frontend: Live Logs Tab**:
    *   A second tab, "Live Agent Logs," will be created.
    *   This tab will contain the existing `Dashboard.tsx` component that connects to the `/api/stream` SSE endpoint to display real-time logs from the agent workflow.

4.  **Frontend: Shareable Links**:
    *   The application state (filters, sorting, etc.) will be encoded into the URL query parameters.
    *   When a user applies a filter, the URL will update automatically (e.g., `/?view=dashboard&filter_project=my-gcp-project`).
    *   A "Share" button will be added to the dashboard, which copies the current URL to the clipboard, allowing users to share their specific view.

5.  **Frontend: CSV Export**:
    *   An "Export to CSV" button will be added to the dashboard.
    *   Clicking this button will trigger a download of the currently filtered and sorted findings as a CSV file. The CSV generation will be handled client-side in the browser.

## Out of Scope
- User authentication and authorization for the dashboard.
- Modifying or remediating findings directly from the dashboard.
- Storing dashboard configurations per user.

## Verification
Every check below is mandatory. Do not skip any.

1.  **API Verification**:
    *   Confirm that a `GET` request to `/api/findings` returns a `200 OK` status and a valid JSON array of findings from the database.
2.  **UI Verification**:
    *   Verify the application loads with two tabs: "CRA Dashboard" and "Live Agent Logs".
    *   Confirm the "CRA Dashboard" tab displays a table of findings fetched from the API.
    *   Test the filtering controls for Organization, Folder, and Project and ensure the table updates correctly.
    *   Verify that clicking the "Share" button copies a URL with the correct query parameters to the clipboard.
    *   Confirm the "Export to CSV" button downloads a CSV file containing the filtered data.
    *   Ensure the "Live Agent Logs" tab correctly streams real-time data from the SSE endpoint.
3.  **Code Quality**:
    *   Ensure all new Go backend code passes `golangci-lint run ./...` with zero issues.
    *   Ensure all new tests pass with `go test ./...`.

## Branch
`feature/cra-dashboard`

## Provenance
`specs/provenance/dashboard-feature.provenance.md` (This file will be overwritten upon execution).
