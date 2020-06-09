package metrics

import (
	"testing"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/stretchr/testify/assert"
)

func TestEscapeMetricNames(t *testing.T) {
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
			Output: "metric1%2Cmetric2%2C",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			safeMetricNames := escapeMetricNames(test.Input)

			assert.Equal(t, safeMetricNames, test.Output)
		})
	}
}
