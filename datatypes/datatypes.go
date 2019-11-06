package datatypes

// Config is a set of config for the app
type Config struct {
	Env    string
	Server string
}

// DeploymentEvents is a collection of Deployment Events
type DeploymentEvents struct {
	Events []DeploymentEvent `json:"events"`
}

// DeploymentEvent defines the data needed from a Dt Deployment Event
type DeploymentEvent struct {
	StartTime         int64  `json:"startTime"`
	EndTime           int64  `json:"endTime"`
	DeploymentName    string `json:"deploymentName"`
	DeploymentVersion string `json:"deploymentVersion"`
}

// DynatraceMetricsResponse defines what we receive from the Dt Metrics v2 API
type DynatraceMetricsResponse struct {
	TotalCount  int     `json:"totalCount"`
	NextPageKey string  `json:"nextPageKey"`
	Metrics     Metrics `json:"metrics"`
}

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

// Metric defines a Dynatrace Service we'd like to investigate and how we'd like to validate it
type Metric struct {
	ID              string
	StaticThreshold int64
	Validation      string
}

// MetricValues defines what we receive for each metric
type MetricValues struct {
	Dimensions []string `json:"dimensions"`
	Timestamp  int64    `json:"timestamp"`
	Value      float64  `json:"value"`
}

// Metrics lists the metrics options we can currently use
type Metrics struct {
	BuiltinServiceResponseTime    BuiltinServiceResponseTime    `json:"builtin:service.response.time:avg"`
	BuiltinServiceErrorsTotalRate BuiltinServiceErrorsTotalRate `json:"builtin:service.errors.total.rate:avg"`
}

// PerformanceSignature is a struct defining all of the parameters we need to calculate a performance signature
type PerformanceSignature struct {
	APIToken  string
	Metrics   []Metric
	ServiceID string
}

// Timestamps represents a start and end time for Deployment events
type Timestamps struct {
	StartTime int64
	EndTime   int64
}
