# Provenance: CRA Compliance Dashboard

## Specification
User request: Create a CRA compliance dashboard with tabs for live logs and finding summaries. Include filtering by GCP hierarchy (Org, Folder, Project), a shareable link feature, and CSV export.

## Execution Plan
1.  **Backend API**: Added `GetAllFindings()` to the `Store` interface (`pkg/store/store.go`), and implemented it in `CloudSQLStore` and `SQLiteStore`. Created `GET /api/findings` in `cmd/server/main.go`.
2.  **Frontend Layout**: Refactored `web/pages/index.tsx` into an MUI Tabbed interface.
3.  **Frontend Component**: Created `web/components/CRADashboard.tsx`.
    *   Implemented data fetching from `/api/findings`.
    *   Implemented client-side parsing of `resource_name` to extract Organization, Folder, and Project for dropdown filters.
    *   Added Chart.js Doughnut chart for high-level compliance visualization.
    *   Implemented `handleShare()` to copy the current URL with search parameters.
    *   Implemented `handleExportCSV()` to download the table data as a CSV blob.

## Deviations
- Bypassed `npm install recharts` due to a corporate network airlock issue restricting `npm install`. Re-used the existing `chart.js` and `react-chartjs-2` libraries already defined in `package.json` to satisfy the data visualization requirement without requiring new dependencies.

## Outcome
- The backend successfully exposes all historical findings.
- The frontend features two distinct views: an interactive, filterable CRA dashboard with export capabilities, and the existing live-streaming agent logs.

## Learnings
- Utilizing pre-existing UI libraries (MUI and Chart.js) already present in the `node_modules` cache or `package.json` is a robust strategy when network configurations block package managers.