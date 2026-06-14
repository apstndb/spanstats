package spanstats_test

import (
	"encoding/json"
	"testing"

	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/types/known/structpb"

	"github.com/apstndb/spanstats"
)

func TestFromResultSetStats(t *testing.T) {
	t.Parallel()

	qs, err := structpb.NewStruct(map[string]any{
		"elapsed_time":      "1.23 msecs",
		"cpu_time":          "0.5 msecs",
		"rows_returned":     "42",
		"rows_scanned":      "100",
		"optimizer_version": "7",
		"query_text":        "SELECT 1",
		"future_key":        "future value",
		"odd_typed":         true,
	})
	if err != nil {
		t.Fatal(err)
	}

	got := spanstats.FromResultSetStats(&sppb.ResultSetStats{QueryStats: qs})
	want := &spanstats.QueryStats{
		ElapsedTime:      "1.23 msecs",
		CPUTime:          "0.5 msecs",
		RowsReturned:     "42",
		RowsScanned:      "100",
		OptimizerVersion: "7",
		QueryText:        "SELECT 1",
		Unknown: map[string]any{
			"future_key": "future value",
			"odd_typed":  true,
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("FromResultSetStats mismatch (-want +got):\n%s", diff)
	}
}

func TestNilSafety(t *testing.T) {
	t.Parallel()

	if got := spanstats.FromResultSetStats(nil); got != nil {
		t.Errorf("FromResultSetStats(nil) = %v, want nil", got)
	}
	if got := spanstats.FromResultSetStats(&sppb.ResultSetStats{}); got != nil {
		t.Errorf("FromResultSetStats(no query_stats) = %v, want nil", got)
	}
	if got := spanstats.FromStruct(nil); got != nil {
		t.Errorf("FromStruct(nil) = %v, want nil", got)
	}
	if got := spanstats.FromMap(nil); got != nil {
		t.Errorf("FromMap(nil) = %v, want nil", got)
	}
}

func TestFromMapEmpty(t *testing.T) {
	t.Parallel()

	got := spanstats.FromMap(map[string]any{})
	if got == nil {
		t.Fatal("FromMap(empty) = nil, want non-nil zero value")
	}
	if got.Unknown != nil {
		t.Errorf("Unknown = %v, want nil", got.Unknown)
	}
}

func TestQueryStatsUnmarshalJSONPreservesUnknown(t *testing.T) {
	t.Parallel()

	var got spanstats.QueryStats
	err := json.Unmarshal([]byte(`{
		"elapsed_time": "1.23 msecs",
		"rows_returned": 42,
		"future_key": "future value",
		"future_object": {"nested": true}
	}`), &got)
	if err != nil {
		t.Fatal(err)
	}

	want := spanstats.QueryStats{
		ElapsedTime: "1.23 msecs",
		Unknown: map[string]any{
			"rows_returned": float64(42),
			"future_key":    "future value",
			"future_object": map[string]any{"nested": true},
		},
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("QueryStats JSON unmarshal mismatch (-want +got):\n%s", diff)
	}
}

func TestQueryStatsMarshalJSONPreservesUnknown(t *testing.T) {
	t.Parallel()

	input := spanstats.QueryStats{
		ElapsedTime:  "1.23 msecs",
		RowsReturned: "42",
		Unknown: map[string]any{
			"future_key":        "future value",
			"query_plan_cached": true,
			"optimizer_version": 7,
		},
	}

	b, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := json.Unmarshal(b, &got); err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"elapsed_time":      "1.23 msecs",
		"rows_returned":     "42",
		"future_key":        "future value",
		"query_plan_cached": true,
		"optimizer_version": float64(7),
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("QueryStats JSON marshal mismatch (-want +got):\n%s", diff)
	}
}

func TestQueryStatsUnmarshalJSONNull(t *testing.T) {
	t.Parallel()

	got := spanstats.QueryStats{ElapsedTime: "previous"}
	if err := json.Unmarshal([]byte(" \n null \t "), &got); err != nil {
		t.Fatal(err)
	}
	want := spanstats.QueryStats{ElapsedTime: "previous"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("QueryStats JSON null mismatch (-want +got):\n%s", diff)
	}
}
