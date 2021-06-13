package metrics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckRelativeThreshold(t *testing.T) {
	type values struct {
		Curr      float64
		Prev      float64
		Threshold float64
	}
	type testDefs struct {
		Name          string
		Values        values
		ExpectPass    bool
		ExpectedError string
		ExpectedText  string
	}

	tests := []testDefs{
		{
			Name: "FAIL - Metric degradation",
			Values: values{
				Curr:      1.0,
				Prev:      0.0,
				Threshold: 0.5,
			},
			ExpectPass:    false,
			ExpectedError: "FAIL - dummy_metric_name:(avg) did not meet the relative threshold criteria. the current performance is 1.00, which is not better than the previous value of 0.00 plus the relative threshold of 0.50",
		},
		{
			Name: "PASS - Passed because threshold",
			Values: values{
				Curr:      1.0,
				Prev:      0.0,
				Threshold: 2,
			},
			ExpectPass:   true,
			ExpectedText: "PASS - dummy_metric_name:(avg)'s current value is 1.00, which is passable compared to the previous results (0.00) plus the tolerance (2.00).\n",
		},
		{
			Name: "PASS - Passed without threshold",
			Values: values{
				Curr:      1.0,
				Prev:      5.0,
				Threshold: 0.5,
			},
			ExpectPass:   true,
			ExpectedText: "PASS - dummy_metric_name:(avg) improvement to 1.00 from 5.00. (Difference: -4.00)\n",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			successText, err := CheckRelativeThreshold(test.Values.Curr, test.Values.Prev, test.Values.Threshold, "dummy_metric_name:(avg)")

			if test.ExpectPass == true {
				assert.NoError(t, err)
				assert.Equal(t, successText, test.ExpectedText)
			} else {
				assert.Error(t, err)
				if test.ExpectedError != "" {
					assert.EqualError(t, err, test.ExpectedError)
				}
			}
		})
	}
}

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
		{
			Name: "Metric degradation",
			Values: values{
				Metric:    1.0,
				Threshold: 0.5,
			},
			ExpectPass:    false,
			ExpectedError: "dummy_metric_name:(avg) was above the static threshold: 1.00, instead of a desired 0.50",
		},
		{
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
		{
			Name: "Metric degradation",
			Values: values{
				Curr: float64(4.0 / 3.0),
				Prev: 1,
			},
			ExpectPass:    false,
			ExpectedError: "dummy_metric_name:(avg) degradation of 0.33",
		},
		{
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
