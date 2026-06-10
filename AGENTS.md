# Agent instructions for `spanstats`

Go library (**MIT**): typed access to Cloud Spanner query statistics
(`ResultSetStats.query_stats`, an untyped protobuf Struct with undocumented
keys). Values stay wire strings; unit/number parsing is deliberately out of
scope for now. Unknown keys (and known keys with non-string values) are
preserved in `QueryStats.Unknown` — never drop data.

## Commands

`mise.toml` owns tasks and tool versions; prefer `mise run check`
(fmt-check, vet, build, test, lint), plus `mise run test-race`. Makefile
delegates to mise. CI via jdx/mise-action.

## Conventions

- The key set tracks what Cloud Spanner is observed to send (source of
  truth: real Spanner responses; spanner-mycli's QueryStats struct is the
  historical reference). Adding a field + `fieldByKey` case is the whole
  change; keep field names matching the snake_case keys.
- Entry points stay nil-safe (PLAN/DML responses often lack query_stats).
- Versioning: stay on v0; breaking = minor, otherwise patch. Per-version
  truth: GitHub Releases. Published versions are immutable — never re-tag.
- English only on github.com.
