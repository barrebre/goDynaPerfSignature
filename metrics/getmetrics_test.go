package metrics

import (
	"testing"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/stretchr/testify/assert"
)

func TestCreateMetricString(t *testing.T) {
	type testDefs struct {
		Name   string
		Input  map[string]datatypes.PSMetric
		Output string
	}

	tests := []testDefs{
		{
			Name: "Metric degradation",
			Input: map[string]datatypes.PSMetric{
				"metric1": {},
				"metric2": {},
			},
			Output: "metric1,metric2,",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			safeMetricNames := createMetricString(test.Input)

			assert.Equal(t, safeMetricNames, test.Output)
		})
	}
}

func TestBuildMetricsQueryURL(t *testing.T) {
	type inputs struct {
		Server       string
		Env          string
		MetricString string
		TS           datatypes.Timestamps
		PS           datatypes.PerformanceSignature
	}

	type testDefs struct {
		Name   string
		Input  inputs
		Output string
	}

	tests := []testDefs{
		{
			Name: "Query with ENV",
			Input: inputs{
				Server:       "myserv",
				Env:          "env1234",
				MetricString: "builtin:service.response.time:(avg),builtin:service.errors.total.rate:(avg),",
				TS:           datatypes.GetSingleTimestamp(),
				PS:           datatypes.GetValidStaticPerformanceSignature(),
			},
			Output: "https://myserv/e/env1234/api/v2/metrics/query?entitySelector=entityId%28%22asdf%22%29&from=1234&metricSelector=builtin%3Aservice.response.time%3A%28avg%29%2Cbuiltin%3Aservice.errors.total.rate%3A%28avg%29%2C&resolution=Inf&to=2345",
		},
		{
			Name: "Query without ENV",
			Input: inputs{
				Server:       "myserv",
				MetricString: "builtin:service.response.time:(avg),builtin:service.errors.total.rate:(avg),",
				TS:           datatypes.GetSingleTimestamp(),
				PS:           datatypes.GetValidStaticPerformanceSignature(),
			},
			Output: "https://myserv/api/v2/metrics/query?entitySelector=entityId%28%22asdf%22%29&from=1234&metricSelector=builtin%3Aservice.response.time%3A%28avg%29%2Cbuiltin%3Aservice.errors.total.rate%3A%28avg%29%2C&resolution=Inf&to=2345",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			metricQueryURL := buildMetricsQueryURL(test.Input.Server, test.Input.Env, test.Input.MetricString, test.Input.TS, test.Input.PS)

			assert.Equal(t, test.Output, metricQueryURL)
		})
	}
}
