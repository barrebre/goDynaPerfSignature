package datatypes

//// Definitions

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
	validDefaultPerformanceSignature = PerformanceSignature{
		APIToken:       "asdf1234",
		EvaluationMins: 5,
		Metrics: []Metric{
			Metric{
				ID: "dummy_metric_name:(avg)",
			},
		},
		ServiceID: "asdf",
	}

	validLargeRelativePerformanceSignature = PerformanceSignature{
		APIToken:       "asdf1234",
		EvaluationMins: 5,
		Metrics: []Metric{
			Metric{
				ID:                "dummy_metric_name:(avg)",
				RelativeThreshold: 20,
				ValidationMethod:  "relative",
			},
		},
		ServiceID: "asdf",
	}

	validSmallRelativePerformanceSignature = PerformanceSignature{
		APIToken:       "asdf1234",
		EvaluationMins: 5,
		Metrics: []Metric{
			Metric{
				ID:                "dummy_metric_name:(avg)",
				RelativeThreshold: 0,
				ValidationMethod:  "relative",
			},
		},
		ServiceID: "asdf",
	}

	validStaticPerformanceSignature = PerformanceSignature{
		APIToken:       "asdf1234",
		EvaluationMins: 5,
		Metrics: []Metric{
			Metric{
				ID:               "dummy_metric_name:percentile(90)",
				StaticThreshold:  1234.1234,
				ValidationMethod: "static",
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

// GetValidLargeRelativePerformanceSignature returns a valid PerformanceSignature with Relative checks and a light sensitivity
func GetValidLargeRelativePerformanceSignature() PerformanceSignature {
	return validLargeRelativePerformanceSignature
}

// GetValidSmallRelativePerformanceSignature returns a valid PerformanceSignature with Relative checks and 0 sensitivity
func GetValidSmallRelativePerformanceSignature() PerformanceSignature {
	return validSmallRelativePerformanceSignature
}

// GetValidStaticPerformanceSignature returns a valid PerformanceSignature with a static check
func GetValidStaticPerformanceSignature() PerformanceSignature {
	return validStaticPerformanceSignature
}
