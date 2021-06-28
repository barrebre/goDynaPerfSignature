package datatypes

//// Definitions

// Metric defines a Dynatrace Service we'd like to investigate and how we'd like to validate it
type PSMetric struct {
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
	Metrics []MetricValuesArray `json:"result"`
}

// MetricValuesArray - The Dynatrace API always returns an array of timestamps, even though we only need the first value each time
type MetricValuesArray struct {
	MetricId     string         `json:"metricId"`
	MetricValues []MetricValues `json:"data"`
}

// MetricValues defines what we receive for each metric
type MetricValues struct {
	Dimensions []string  `json:"dimensions"`
	Timestamps []int64   `json:"timestamps"`
	Values     []float64 `json:"values"`
}

//// Example Values
var (
	validFailingComparisonMetrics = ComparisonMetrics{
		CurrentMetrics: DynatraceMetricsResponse{
			Metrics: []MetricValuesArray{
				{
					MetricId: "dummy_metric_name:avg",
					MetricValues: []MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamps: []int64{1234},
							Values:     []float64{1235},
						},
					},
				},
				{
					MetricId: "dummy_metric_name:percentile(90)",
					MetricValues: []MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamps: []int64{23456},
							Values:     []float64{23456},
						},
					},
				},
			},
		},
		PreviousMetrics: DynatraceMetricsResponse{
			Metrics: []MetricValuesArray{
				{
					MetricId: "dummy_metric_name:avg",
					MetricValues: []MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamps: []int64{2345},
							Values:     []float64{1234.1234},
						},
					},
				},
				{
					MetricId: "dummy_metric_name:percentile(90)",
					MetricValues: []MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamps: []int64{3456},
							Values:     []float64{2345.1234},
						},
					},
				},
			},
		},
	}

	validPassingComparisonMetrics = ComparisonMetrics{
		CurrentMetrics: DynatraceMetricsResponse{
			Metrics: []MetricValuesArray{
				{
					MetricId: "dummy_metric_name:avg",
					MetricValues: []MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamps: []int64{1234},
							Values:     []float64{1234.1234},
						},
					},
				},
			},
		},
		PreviousMetrics: DynatraceMetricsResponse{
			Metrics: []MetricValuesArray{
				{
					MetricId: "dummy_metric_name:avg",
					MetricValues: []MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamps: []int64{2345},
							Values:     []float64{1235},
						},
					},
				},
			},
		},
	}

	missingComparisonMetrics = ComparisonMetrics{
		CurrentMetrics: DynatraceMetricsResponse{
			Metrics: []MetricValuesArray{},
		},
		PreviousMetrics: DynatraceMetricsResponse{
			Metrics: []MetricValuesArray{},
		},
	}

	missingPreviousComparisonMetrics = ComparisonMetrics{
		CurrentMetrics: DynatraceMetricsResponse{
			Metrics: []MetricValuesArray{
				{
					MetricId: "dummy_metric_name:avg",
					MetricValues: []MetricValues{
						{
							Values: []float64{12.34},
						},
					},
				},
				{
					MetricId: "dummy_metric_name:percentile(90)",
					MetricValues: []MetricValues{
						{
							Values: []float64{12.34},
						},
					},
				},
			},
		},
		PreviousMetrics: DynatraceMetricsResponse{
			Metrics: []MetricValuesArray{
				{
					MetricId:     "dummy_metric_name:avg",
					MetricValues: nil,
				},
				{
					MetricId:     "dummy_metric_name:percentile(90)",
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
