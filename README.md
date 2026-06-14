# spanstats

[![Go Reference](https://pkg.go.dev/badge/github.com/apstndb/spanstats.svg)](https://pkg.go.dev/github.com/apstndb/spanstats)

Typed access to Cloud Spanner query statistics
(`ResultSetStats.query_stats`), which the API exposes as an untyped protobuf
`Struct` with undocumented keys that every tool ends up re-parsing by hand.

```go
stats := spanstats.FromResultSetStats(rss) // or FromStruct / FromMap
if stats != nil {
    fmt.Println(stats.ElapsedTime, stats.RowsScanned)
}
```

- Values stay the wire strings Spanner sends (for example `"1.23 msecs"`);
  no unit or number parsing (deliberately out of MVP scope).
- Keys not yet modeled — and known keys with unexpectedly non-string
  values — are preserved in `Unknown`, so nothing is dropped when Spanner
  adds keys.
- JSON marshaling/unmarshaling preserves those unknown keys too, which keeps
  stored query profile payloads forward-compatible.
- Entry points for all three shapes tools encounter:
  `FromResultSetStats(*sppb.ResultSetStats)`, `FromStruct(*structpb.Struct)`,
  and `FromMap(map[string]any)` (the `spanner.RowIterator.QueryStats` shape).
- Nil-safe: PLAN-mode and DML responses often carry no `query_stats`; all
  entry points return nil for absent input.

Extracted from the duplicated parsing in
[spanner-mycli](https://github.com/apstndb/spanner-mycli) and
[spannersh](https://github.com/apstndb/spannersh).

**Status: experimental.** The key set tracks what Cloud Spanner is observed
to send; additions are non-breaking (new fields, with `Unknown` as the
forward-compatibility net).

## License

MIT
