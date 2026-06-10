// Package spanstats provides typed access to Cloud Spanner query statistics
// (ResultSetStats.query_stats), which the API exposes as an untyped
// protobuf Struct with undocumented keys.
//
// All values are kept as the wire strings Spanner sends (for example
// "1.23 msecs"); this package does not parse units or numbers. Keys not yet
// modeled are preserved in [QueryStats.Unknown], and known keys whose value
// is not a string are routed there too, so no information is dropped.
package spanstats

import (
	sppb "cloud.google.com/go/spanner/apiv1/spannerpb"
	"google.golang.org/protobuf/types/known/structpb"
)

// QueryStats is the typed view of query_stats. Field coverage follows the
// keys observed from Cloud Spanner; absent keys are zero values.
type QueryStats struct {
	ElapsedTime                string `json:"elapsed_time,omitempty"`
	CPUTime                    string `json:"cpu_time,omitempty"`
	RowsReturned               string `json:"rows_returned,omitempty"`
	RowsScanned                string `json:"rows_scanned,omitempty"`
	DeletedRowsScanned         string `json:"deleted_rows_scanned,omitempty"`
	OptimizerVersion           string `json:"optimizer_version,omitempty"`
	OptimizerStatisticsPackage string `json:"optimizer_statistics_package,omitempty"`
	RemoteServerCalls          string `json:"remote_server_calls,omitempty"`
	MemoryPeakUsageBytes       string `json:"memory_peak_usage_bytes,omitempty"`
	TotalMemoryPeakUsageByte   string `json:"total_memory_peak_usage_byte,omitempty"`
	QueryText                  string `json:"query_text,omitempty"`
	BytesReturned              string `json:"bytes_returned,omitempty"`
	RuntimeCreationTime        string `json:"runtime_creation_time,omitempty"`
	StatisticsLoadTime         string `json:"statistics_load_time,omitempty"`
	MemoryUsagePercentage      string `json:"memory_usage_percentage,omitempty"`
	FilesystemDelaySeconds     string `json:"filesystem_delay_seconds,omitempty"`
	LockingDelay               string `json:"locking_delay,omitempty"`
	QueryPlanCreationTime      string `json:"query_plan_creation_time,omitempty"`
	ServerQueueDelay           string `json:"server_queue_delay,omitempty"`
	DataBytesRead              string `json:"data_bytes_read,omitempty"`
	IsGraphQuery               string `json:"is_graph_query,omitempty"`
	RuntimeCached              string `json:"runtime_cached,omitempty"`
	QueryPlanCached            string `json:"query_plan_cached,omitempty"`

	// Unknown preserves keys this package does not model (and known keys
	// whose value was not a string). A nil map means none were present.
	Unknown map[string]any `json:"-"`
}

// fieldByKey maps a query_stats key to the corresponding string field.
func (s *QueryStats) fieldByKey(key string) *string {
	switch key {
	case "elapsed_time":
		return &s.ElapsedTime
	case "cpu_time":
		return &s.CPUTime
	case "rows_returned":
		return &s.RowsReturned
	case "rows_scanned":
		return &s.RowsScanned
	case "deleted_rows_scanned":
		return &s.DeletedRowsScanned
	case "optimizer_version":
		return &s.OptimizerVersion
	case "optimizer_statistics_package":
		return &s.OptimizerStatisticsPackage
	case "remote_server_calls":
		return &s.RemoteServerCalls
	case "memory_peak_usage_bytes":
		return &s.MemoryPeakUsageBytes
	case "total_memory_peak_usage_byte":
		return &s.TotalMemoryPeakUsageByte
	case "query_text":
		return &s.QueryText
	case "bytes_returned":
		return &s.BytesReturned
	case "runtime_creation_time":
		return &s.RuntimeCreationTime
	case "statistics_load_time":
		return &s.StatisticsLoadTime
	case "memory_usage_percentage":
		return &s.MemoryUsagePercentage
	case "filesystem_delay_seconds":
		return &s.FilesystemDelaySeconds
	case "locking_delay":
		return &s.LockingDelay
	case "query_plan_creation_time":
		return &s.QueryPlanCreationTime
	case "server_queue_delay":
		return &s.ServerQueueDelay
	case "data_bytes_read":
		return &s.DataBytesRead
	case "is_graph_query":
		return &s.IsGraphQuery
	case "runtime_cached":
		return &s.RuntimeCached
	case "query_plan_cached":
		return &s.QueryPlanCached
	}
	return nil
}

// FromMap builds a [QueryStats] from a decoded query_stats map, the shape of
// [cloud.google.com/go/spanner.RowIterator] QueryStats. A nil map returns nil.
func FromMap(m map[string]any) *QueryStats {
	if m == nil {
		return nil
	}
	var s QueryStats
	for key, value := range m {
		if str, ok := value.(string); ok {
			if field := s.fieldByKey(key); field != nil {
				*field = str
				continue
			}
		}
		if s.Unknown == nil {
			s.Unknown = make(map[string]any)
		}
		s.Unknown[key] = value
	}
	return &s
}

// FromStruct builds a [QueryStats] from the raw query_stats protobuf Struct.
// A nil Struct returns nil.
func FromStruct(s *structpb.Struct) *QueryStats {
	if s == nil {
		return nil
	}
	return FromMap(s.AsMap())
}

// FromResultSetStats builds a [QueryStats] from a
// [cloud.google.com/go/spanner/apiv1/spannerpb.ResultSetStats]. A nil stats
// message or absent query_stats returns nil (PLAN-mode and DML responses
// often carry no query_stats).
func FromResultSetStats(stats *sppb.ResultSetStats) *QueryStats {
	return FromStruct(stats.GetQueryStats())
}
