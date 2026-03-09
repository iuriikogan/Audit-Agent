# Provenance: Database Selection Refactoring

## Specification
User request: Support `CLOUD_SQL` and `SQLITE_MEM` via `DATABASE_TYPE` env var.

## Execution Plan
1.  **Config**: Added `DatabaseType` to `pkg/config`.
2.  **SQLite Implementation**: Created `pkg/store/sqlite.go` using `modernc.org/sqlite` (pure Go).
3.  **Selection Logic**: Updated `cmd/server/main.go` to switch between `CLOUD_SQL` and `SQLITE_MEM`.
4.  **Verification**: Added `pkg/store/sqlite_test.go` and verified with `go test` and `golangci-lint`.

## Deviations
- None.

## Outcome
- Users can now deploy with either Cloud SQL or in-memory SQLite by setting `DATABASE_TYPE`.
- Pure Go SQLite driver ensures cross-platform compatibility without CGO.

## Learnings
- `modernc.org/sqlite` is a convenient drop-in for `sql.DB` without CGO dependencies.
