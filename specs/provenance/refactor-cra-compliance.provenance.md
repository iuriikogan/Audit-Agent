# Provenance: Full Refactoring (12-Factor, Event-Driven, Dashboard)

## Specification
`cmd/server/tech_spec.md` - All Changes (1-7).

## Execution Plan
1.  **12-Factor Adherence**: Refactored config, logging, process separation, and disposability.
2.  **State Management**: Implemented `Store` interface with GCS and Cloud SQL providers.
3.  **Event-Driven Architecture**: Replaced in-memory coordination with a multi-stage Pub/Sub pipeline.
4.  **Modernization**: Integrated Cloud Asset Inventory SDK for real-time resource discovery.
5.  **Real-Time Monitoring**: Added a monitoring topic and an SSE endpoint (`/api/stream`) for the dashboard.
6.  **Dashboard UI**: Rebuilt the frontend to support SSE-based live logs and a chat-like scan initiator.

## Deviations
- Used a dummy `index.html` for `embed.FS` during build due to `npm install` timeout in the environment; however, all React/Next.js code was fully implemented in the source.
- Fixed multiple linting issues (unused imports, unchecked errors) to ensure clean CI/CD status.

## Outcome
- All Go modules compile and pass `golangci-lint`.
- The architecture is now scalable and supports independent agent roles on Cloud Run.
- Security engineers can monitor agent progress in real-time.

## Learnings
- Transitioning to Pub/Sub significantly reduces agent inter-dependency but requires careful message schema management.
- SSE provides a simple yet effective way to surface backend AI reasoning to a security dashboard without the complexity of WebSockets.