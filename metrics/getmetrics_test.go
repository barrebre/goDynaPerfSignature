package metrics

import (
	"testing"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/stretchr/testify/assert"
)

func TestCreateMetricString(t *testing.T) {
	type testDefs struct {
		Name   string
		Input  []datatypes.Metric
		Output string
	}

	tests := []testDefs{
		testDefs{
			Name: "Metric degradation",
			Input: []datatypes.Metric{
				datatypes.Metric{
					ID: "metric1",
				},
				datatypes.Metric{
					ID: "metric2",
				},
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
		testDefs{
			Name: "Query with ENV",
			Input: inputs{
				Server:       "myserv",
				Env:          "env1234",
				MetricString: "builtin:service.response.time:(avg),builtin:service.errors.total.rate:(avg),",
				TS:           datatypes.GetSingleTimestamp(),
				PS:           datatypes.GetValidStaticPerformanceSignature(),
			},
			Output: "https://myserv/e/env1234/api/v2/metrics/series/builtin:service.response.time:%28avg%29,builtin:service.errors.total.rate:%28avg%29,?from=1234&resolution=Inf&scope=entity%28asdf%29&to=2345",
		},
		testDefs{
			Name: "Query without ENV",
			Input: inputs{
				Server:       "myserv",
				MetricString: "builtin:service.response.time:(avg),builtin:service.errors.total.rate:(avg),",
				TS:           datatypes.GetSingleTimestamp(),
				PS:           datatypes.GetValidStaticPerformanceSignature(),
			},
			Output: "https://myserv/api/v2/metrics/series/builtin:service.response.time:%28avg%29,builtin:service.errors.total.rate:%28avg%29,?from=1234&resolution=Inf&scope=entity%28asdf%29&to=2345",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			metricQueryURL := buildMetricsQueryURL(test.Input.Server, test.Input.Env, test.Input.MetricString, test.Input.TS, test.Input.PS)

			assert.Equal(t, metricQueryURL, test.Output)
		})
	}
}
