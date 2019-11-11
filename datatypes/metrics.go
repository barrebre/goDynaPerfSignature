package datatypes

// ComparisonMetrics has a current and previous set of metrics to compare
type ComparisonMetrics struct {
	CurrentMetrics  DynatraceMetricsResponse
	PreviousMetrics DynatraceMetricsResponse
}

// DynatraceMetricsResponse defines what we receive from the Dt Metrics v2 API
type DynatraceMetricsResponse struct {
	Metrics map[string]MetricValuesArray `json:"metrics"`
}

// MetricValuesArray - The Dynatrace API always returns an array of timestamps, even though we only need the first value each time
type MetricValuesArray struct {
	MetricValues []MetricValues `json:"values"`
}

// MetricValues defines what we receive for each metric
type MetricValues struct {
	Dimensions []string `json:"dimensions"`
	Timestamp  int64    `json:"timestamp"`
	Value      float64  `json:"value"`
}
