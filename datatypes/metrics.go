package datatypes

// Both of the following unique definitions could probably be updated to some map syntax

// BuiltinServiceResponseTime has to be defined for some reason
type BuiltinServiceResponseTime struct {
	MetricValues []MetricValues `json:"values"`
}

// BuiltinServiceErrorsTotalRate has to be defined for some reason
type BuiltinServiceErrorsTotalRate struct {
	MetricValues []MetricValues `json:"values"`
}

// ComparisonMetrics has a current and previous set of metrics to compare
type ComparisonMetrics struct {
	CurrentMetrics  Metrics
	PreviousMetrics Metrics
}

// DynatraceMetricsResponse defines what we receive from the Dt Metrics v2 API
type DynatraceMetricsResponse struct {
	TotalCount  int     `json:"totalCount"`
	NextPageKey string  `json:"nextPageKey"`
	Metrics     Metrics `json:"metrics"`
}

// Metric defines a Dynatrace Service we'd like to investigate and how we'd like to validate it
type Metric struct {
	ID              string
	StaticThreshold int64
	Validation      string
}

// Metrics lists the metrics options we can currently use
type Metrics struct {
	BuiltinServiceResponseTime    BuiltinServiceResponseTime    `json:"builtin:service.response.time:avg"`
	BuiltinServiceErrorsTotalRate BuiltinServiceErrorsTotalRate `json:"builtin:service.errors.total.rate:avg"`
}

// MetricValues defines what we receive for each metric
type MetricValues struct {
	Dimensions []string `json:"dimensions"`
	Timestamp  int64    `json:"timestamp"`
	Value      float64  `json:"value"`
}
