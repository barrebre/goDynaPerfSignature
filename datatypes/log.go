package datatypes

// Logging defines what can be included in log files
type Logging struct {
	Message string               `json:"message"`
	PerfSig PerformanceSignature `json:"perf_sig"`
}
