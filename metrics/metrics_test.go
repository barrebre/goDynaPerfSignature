package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckStaticThreshold(t *testing.T) {
	type values struct {
		Metric    float64
		Threshold float64
	}
	type testDefs struct {
		Name          string
		Values        values
		ExpectPass    bool
		ExpectedError string
	}

	tests := []testDefs{
		testDefs{
			Name: "Metric degradation",
			Values: values{
				Metric:    1.0,
				Threshold: 0.5,
			},
			ExpectPass:    false,
			ExpectedError: "dummy_metric_name:(avg) was above the static threshold: 1, instead of a desired 0.5",
		},
		testDefs{
			Name: "Successful deploy",
			Values: values{
				Metric:    0.0,
				Threshold: 1.0,
			},
			ExpectPass: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			_, err := CheckStaticThreshold(test.Values.Metric, test.Values.Threshold, "dummy_metric_name:(avg)")

			if test.ExpectPass == true {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				if test.ExpectedError != "" {
					assert.EqualError(t, err, test.ExpectedError)
				}
			}
		})
	}
}

func TestCompareMetrics(t *testing.T) {
	type values struct {
		Curr float64
		Prev float64
	}
	type testDefs struct {
		Name          string
		Values        values
		ExpectPass    bool
		ExpectedError string
	}

	tests := []testDefs{
		testDefs{
			Name: "Metric degradation",
			Values: values{
				Curr: float64(4.0 / 3.0),
				Prev: 1,
			},
			ExpectPass:    false,
			ExpectedError: "dummy_metric_name:(avg) degradation of 0.33",
		},
		testDefs{
			Name: "Successful deploy",
			Values: values{
				Curr: 0.0,
				Prev: 1.0,
			},
			ExpectPass: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			_, err := CompareMetrics(test.Values.Curr, test.Values.Prev, "dummy_metric_name:(avg)")

			if test.ExpectPass == true {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				if test.ExpectedError != "" {
					assert.EqualError(t, err, test.ExpectedError)
				}
			}
		})
	}
}
