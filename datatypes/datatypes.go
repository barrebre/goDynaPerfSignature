package datatypes

// Metric defines a Dynatrace Service we'd like to investigate and how we'd like to validate it
type Metric struct {
	ID               string
	StaticThreshold  float64
	ValidationMethod string
}

// PerformanceSignature is a struct defining all of the parameters we need to calculate a performance signature
type PerformanceSignature struct {
	APIToken  string
	Metrics   []Metric
	ServiceID string
}
