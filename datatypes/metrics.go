package datatypes

//// Definitions

// Metric defines a Dynatrace Service we'd like to investigate and how we'd like to validate it
type Metric struct {
	ID                string
	RelativeThreshold float64
	StaticThreshold   float64
	ValidationMethod  string
}

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

//// Example Values
var (
	validFailingComparisonMetrics = ComparisonMetrics{
		CurrentMetrics: DynatraceMetricsResponse{
			Metrics: map[string]MetricValuesArray{
				"dummy_metric_name:avg": {
					MetricValues: []MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamp: 1234,
							Value:     1235,
						},
					},
				},
				"dummy_metric_name:percentile(90)": {
					MetricValues: []MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamp: 1234,
							Value:     1235,
						},
					},
				},
			},
		},
		PreviousMetrics: DynatraceMetricsResponse{
			Metrics: map[string]MetricValuesArray{
				"dummy_metric_name:avg": {
					MetricValues: []MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamp: 2345,
							Value:     1234.1234,
						},
					},
				},
				"dummy_metric_name:percentile(90)": {
					MetricValues: []MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamp: 2345,
							Value:     1234.1234,
						},
					},
				},
			},
		},
	}

	validPassingComparisonMetrics = ComparisonMetrics{
		CurrentMetrics: DynatraceMetricsResponse{
			Metrics: map[string]MetricValuesArray{
				"dummy_metric_name:avg": {
					MetricValues: []MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamp: 1234,
							Value:     1234.1234,
						},
					},
				},
			},
		},
		PreviousMetrics: DynatraceMetricsResponse{
			Metrics: map[string]MetricValuesArray{
				"dummy_metric_name:avg": {
					MetricValues: []MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamp: 2345,
							Value:     1235,
						},
					},
				},
			},
		},
	}

	missingComparisonMetrics = ComparisonMetrics{
		CurrentMetrics: DynatraceMetricsResponse{
			Metrics: map[string]MetricValuesArray{},
		},
		PreviousMetrics: DynatraceMetricsResponse{
			Metrics: map[string]MetricValuesArray{},
		},
	}

	missingPreviousComparisonMetrics = ComparisonMetrics{
		CurrentMetrics: DynatraceMetricsResponse{
			Metrics: map[string]MetricValuesArray{
				"dummy_metric_name:avg": {
					MetricValues: []MetricValues{
						{},
					},
				},
				"dummy_metric_name:percentile(90)": {
					MetricValues: []MetricValues{
						{},
					},
				},
			},
		},
		PreviousMetrics: DynatraceMetricsResponse{
			Metrics: map[string]MetricValuesArray{
				"dummy_metric_name:avg": {
					MetricValues: nil,
				},
				"dummy_metric_name:percentile(90)": {
					MetricValues: nil,
				},
			},
		},
	}
)

//// Example Accessors

// GetMissingComparisonMetrics returns a ComparisonMetrics missing Metrics
func GetMissingComparisonMetrics() ComparisonMetrics {
	return missingComparisonMetrics
}

// GetMissingPreviousComparisonMetrics returns a ComparisonMetrics missing Metrics
func GetMissingPreviousComparisonMetrics() ComparisonMetrics {
	return missingPreviousComparisonMetrics
}

// GetValidPassingComparisonMetrics returns a valid ComparisonMetrics which passes
func GetValidPassingComparisonMetrics() ComparisonMetrics {
	return validPassingComparisonMetrics
}

// GetValidFailingComparisonMetrics returns a valid ComparisonMetrics which fails
func GetValidFailingComparisonMetrics() ComparisonMetrics {
	return validFailingComparisonMetrics
}
