package datatypes

// PerformanceSignature is a struct defining all of the parameters we need to calculate a performance signature
type PerformanceSignature struct {
	APIToken  string
	Metrics   []Metric
	ServiceID string
}
