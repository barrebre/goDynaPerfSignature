package exampletestdata

import (
	"github.com/barrebre/goDynaPerfSignature/datatypes"
)

var (
	validFailingComparisonMetrics = datatypes.ComparisonMetrics{
		CurrentMetrics: datatypes.DynatraceMetricsResponse{
			Metrics: map[string]datatypes.MetricValuesArray{
				"TestMetric": datatypes.MetricValuesArray{
					MetricValues: []datatypes.MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamp: 1234,
							Value:     123456789.1234,
						},
					},
				},
			},
		},
		PreviousMetrics: datatypes.DynatraceMetricsResponse{
			Metrics: map[string]datatypes.MetricValuesArray{
				"TestMetric": datatypes.MetricValuesArray{
					MetricValues: []datatypes.MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamp: 0123,
							Value:     0123.0123,
						},
					},
				},
			},
		},
	}

	validPassingComparisonMetrics = datatypes.ComparisonMetrics{
		CurrentMetrics: datatypes.DynatraceMetricsResponse{
			Metrics: map[string]datatypes.MetricValuesArray{
				"TestMetric": datatypes.MetricValuesArray{
					MetricValues: []datatypes.MetricValues{
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
		PreviousMetrics: datatypes.DynatraceMetricsResponse{
			Metrics: map[string]datatypes.MetricValuesArray{
				"TestMetric": datatypes.MetricValuesArray{
					MetricValues: []datatypes.MetricValues{
						{
							Dimensions: []string{
								"dim1",
							},
							Timestamp: 2345,
							Value:     2345.2345,
						},
					},
				},
			},
		},
	}

	missingComparisonMetrics = datatypes.ComparisonMetrics{
		CurrentMetrics: datatypes.DynatraceMetricsResponse{
			Metrics: map[string]datatypes.MetricValuesArray{},
		},
		PreviousMetrics: datatypes.DynatraceMetricsResponse{
			Metrics: map[string]datatypes.MetricValuesArray{},
		},
	}

	missingPreviousComparisonMetrics = datatypes.ComparisonMetrics{
		CurrentMetrics: datatypes.DynatraceMetricsResponse{
			Metrics: map[string]datatypes.MetricValuesArray{
				"TestMetric": datatypes.MetricValuesArray{
					MetricValues: []datatypes.MetricValues{
						{},
					},
				},
			},
		},
		PreviousMetrics: datatypes.DynatraceMetricsResponse{
			Metrics: map[string]datatypes.MetricValuesArray{
				"TestMetric": datatypes.MetricValuesArray{
					MetricValues: nil,
				},
			},
		},
	}

	validStaticPerformanceSignature = datatypes.PerformanceSignature{
		APIToken: "asdf1234",
		Metrics: []datatypes.Metric{
			datatypes.Metric{
				ID:               "TestMetric",
				StaticThreshold:  1234.1234,
				ValidationMethod: "static",
			},
		},
		ServiceID: "asdf",
	}

	validDefaultPerformanceSignature = datatypes.PerformanceSignature{
		APIToken: "asdf1234",
		Metrics: []datatypes.Metric{
			datatypes.Metric{
				ID:              "TestMetric",
				StaticThreshold: 1234.1234,
			},
		},
		ServiceID: "asdf",
	}
)

// GetValidDefaultPerformanceSignature returns a valid PerformanceSignature with default checks
func GetValidDefaultPerformanceSignature() datatypes.PerformanceSignature {
	return validDefaultPerformanceSignature
}

// GetValidStaticPerformanceSignature returns a valid PerformanceSignature with a static check
func GetValidStaticPerformanceSignature() datatypes.PerformanceSignature {
	return validStaticPerformanceSignature
}

// GetMissingComparisonMetrics returns a ComparisonMetrics missing Metrics
func GetMissingComparisonMetrics() datatypes.ComparisonMetrics {
	return missingComparisonMetrics
}

// GetMissingPreviousComparisonMetrics returns a ComparisonMetrics missing Metrics
func GetMissingPreviousComparisonMetrics() datatypes.ComparisonMetrics {
	return missingPreviousComparisonMetrics
}

// GetValidPassingComparisonMetrics returns a valid ComparisonMetrics which passes
func GetValidPassingComparisonMetrics() datatypes.ComparisonMetrics {
	return validPassingComparisonMetrics
}

// GetValidFailingComparisonMetrics returns a valid ComparisonMetrics which fails
func GetValidFailingComparisonMetrics() datatypes.ComparisonMetrics {
	return validFailingComparisonMetrics
}
