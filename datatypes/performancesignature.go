package datatypes

//// Definitions

// PerformanceSignature is a struct defining all of the parameters we need to calculate a performance signature
type PerformanceSignature struct {
	APIToken       string
	DTEnv          string
	DTServer       string
	EvaluationMins int
	EventAge       int
	PSMetrics      map[string]PSMetric
	ServiceID      string
}

// PerformanceSignatureReturn defines the spec for what needs to be returned to the requester
type PerformanceSignatureReturn struct {
	Error    bool
	Pass     bool
	Response []string
}

//// Example Values
var (
	validDefaultPerformanceSignature = PerformanceSignature{
		APIToken:       "asdf1234",
		EvaluationMins: 5,
		EventAge:       1598818148,
		PSMetrics: map[string]PSMetric{
			"dummy_metric_name:avg": {},
		},
		ServiceID: "asdf",
	}

	validLargeRelativePerformanceSignature = PerformanceSignature{
		APIToken:       "asdf1234",
		EvaluationMins: 5,
		EventAge:       10234,
		PSMetrics: map[string]PSMetric{
			"dummy_metric_name:avg": {
				RelativeThreshold: 20,
				ValidationMethod:  "relative",
			},
		},
		ServiceID: "asdf",
	}

	validSmallRelativePerformanceSignature = PerformanceSignature{
		APIToken:       "asdf1234",
		EvaluationMins: 5,
		EventAge:       992348,
		PSMetrics: map[string]PSMetric{
			"dummy_metric_name:avg": {
				RelativeThreshold: 0,
				ValidationMethod:  "relative",
			},
		},
		ServiceID: "asdf",
	}

	validStaticPerformanceSignature = PerformanceSignature{
		APIToken:       "asdf1234",
		EvaluationMins: 5,
		PSMetrics: map[string]PSMetric{
			"dummy_metric_name:percentile(90)": {
				StaticThreshold:  1234.1234,
				ValidationMethod: "static",
			},
		},
		ServiceID: "asdf",
	}

	validPerformanceSignatureReturnSuccess = PerformanceSignatureReturn{
		Error:    false,
		Pass:     true,
		Response: []string{"PASS - builtin:service.response.time:avg improvement to 82122.06 from 150879.00. (Difference: -68756.94)"},
	}

	validPerformanceSignatureReturnFailure = PerformanceSignatureReturn{
		Pass:     false,
		Response: []string{"No previous metrics to compare against for metric dummy_metric_name:avg", "PASS - dummy_metric_name:percentile(90) is below the static threshold (1234.12) with a value of 12.34."},
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

// GetValidPerformanceSignatureReturnSuccess returns a PerformanceSignatureReturn that passed
func GetValidPerformanceSignatureReturnSuccess() PerformanceSignatureReturn {
	return validPerformanceSignatureReturnSuccess
}

// GetValidPerformanceSignatureReturnFailure returns a PerformanceSignatureReturn that failed
func GetValidPerformanceSignatureReturnFailure() PerformanceSignatureReturn {
	return validPerformanceSignatureReturnFailure
}
