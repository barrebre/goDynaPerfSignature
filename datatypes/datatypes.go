package datatypes

//// Definitions

// Metric defines a Dynatrace Service we'd like to investigate and how we'd like to validate it
type Metric struct {
	ID               string
	StaticThreshold  float64
	ValidationMethod string
}

// PerformanceSignature is a struct defining all of the parameters we need to calculate a performance signature
type PerformanceSignature struct {
	APIToken       string
	DTEnv          string
	DTServer       string
	EvaluationMins int
	Metrics        []Metric
	ServiceID      string
}

//// Example Values
var (
	validStaticPerformanceSignature = PerformanceSignature{
		APIToken:       "asdf1234",
		EvaluationMins: 5,
		Metrics: []Metric{
			Metric{
				ID:               "TestMetric",
				StaticThreshold:  1234.1234,
				ValidationMethod: "static",
			},
		},
		ServiceID: "asdf",
	}

	validDefaultPerformanceSignature = PerformanceSignature{
		APIToken:       "asdf1234",
		EvaluationMins: 5,
		Metrics: []Metric{
			Metric{
				ID:              "TestMetric",
				StaticThreshold: 1234.1234,
			},
		},
		ServiceID: "asdf",
	}
)

//// Example Accessors

// GetValidDefaultPerformanceSignature returns a valid PerformanceSignature with default checks
func GetValidDefaultPerformanceSignature() PerformanceSignature {
	return validDefaultPerformanceSignature
}

// GetValidStaticPerformanceSignature returns a valid PerformanceSignature with a static check
func GetValidStaticPerformanceSignature() PerformanceSignature {
	return validStaticPerformanceSignature
}
