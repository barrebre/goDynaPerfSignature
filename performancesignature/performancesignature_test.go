package performancesignature

import (
	"fmt"
	"testing"
	"time"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/barrebre/goDynaPerfSignature/logging"

	"github.com/stretchr/testify/assert"
)

func TestBuildDeploymentRequest(t *testing.T) {
	type testDefs struct {
		Name          string
		Values        datatypes.PerformanceSignature
		ExpectPass    bool
		ExpectedError string
	}

	tests := []testDefs{
		{
			Name:       "Pass - valid params",
			Values:     datatypes.GetValidDefaultPerformanceSignature(),
			ExpectPass: true,
		},
		{
			Name: "Pass - valid params with env",
			Values: datatypes.PerformanceSignature{
				DTEnv:    "POWEIFJPIOJAPSOIJ",
				DTServer: "asdf1234.live.dynatrace.com",
			},
			ExpectPass: true,
		},
		{
			Name: "Fail - couldn't build HTTP request",
			Values: datatypes.PerformanceSignature{
				DTEnv:    "",
				DTServer: "\\",
			},
			ExpectPass:    false,
			ExpectedError: "parse https://\\/api/v1/events?eventType=CUSTOM_DEPLOYMENT&entityId=: invalid character \"\\\\\" in host name",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			_, err := buildDeploymentRequest(test.Values)

			if test.ExpectPass == true {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, err.Error(), test.ExpectedError)
			}
		})
	}
}

func TestCheckPerfSignature(t *testing.T) {
	type testDefs struct {
		Name             string
		PerfSignature    datatypes.PerformanceSignature
		MetricsResponse  datatypes.ComparisonMetrics
		ExpectedResponse []string
	}

	tests := []testDefs{
		{
			Name:             "Valid Default Check Failing Data",
			PerfSignature:    datatypes.GetValidDefaultPerformanceSignature(),
			MetricsResponse:  datatypes.GetValidFailingComparisonMetrics(),
			ExpectedResponse: []string{"Metric degradation found: FAIL - dummy_metric_name:avg had a degradation of 0.88, from 1234.12 to 1235.00"},
		},
		{
			Name:             "Valid Relative Check Failing Data",
			PerfSignature:    datatypes.GetValidSmallRelativePerformanceSignature(),
			MetricsResponse:  datatypes.GetValidFailingComparisonMetrics(),
			ExpectedResponse: []string{"Metric degradation found: FAIL - dummy_metric_name:avg did not meet the relative threshold criteria. the current performance is 1235.00, which is not better than the previous value of 1234.12 plus the relative threshold of 0.00."},
		},
		{
			Name:             "Valid Relative Check Passing Data",
			PerfSignature:    datatypes.GetValidLargeRelativePerformanceSignature(),
			MetricsResponse:  datatypes.GetValidPassingComparisonMetrics(),
			ExpectedResponse: []string{"PASS - dummy_metric_name:avg improvement to 1234.12 from 1235.00. (Difference: -0.88)"},
		},
		{
			Name:             "Valid Static Check Failing Data",
			PerfSignature:    datatypes.GetValidStaticPerformanceSignature(),
			MetricsResponse:  datatypes.GetValidFailingComparisonMetrics(),
			ExpectedResponse: []string{"Metric degradation found: FAIL - dummy_metric_name:percentile(90) was above the static threshold: 1235.00, instead of a desired 1234.12"},
		},
		{
			Name:             "Valid Default Check Passing Data",
			PerfSignature:    datatypes.GetValidDefaultPerformanceSignature(),
			MetricsResponse:  datatypes.GetValidPassingComparisonMetrics(),
			ExpectedResponse: []string{"PASS - Successful deploy! Improvement of 0.88, from 1235.00 to 1234.12"},
		},
		{
			Name:             "No Data Returned",
			PerfSignature:    datatypes.GetValidStaticPerformanceSignature(),
			MetricsResponse:  datatypes.GetMissingComparisonMetrics(),
			ExpectedResponse: []string{"there were no current metrics found for dummy_metric_name:percentile(90)"},
		},
		{
			Name:             "No Previous Deployment Data Returned - Default Check",
			PerfSignature:    datatypes.GetValidDefaultPerformanceSignature(),
			MetricsResponse:  datatypes.GetMissingPreviousComparisonMetrics(),
			ExpectedResponse: []string{"No previous metrics to compare against for metric dummy_metric_name:avg"},
		},
		{
			Name:             "No Previous Deployment Data Returned - Static Check",
			PerfSignature:    datatypes.GetValidStaticPerformanceSignature(),
			MetricsResponse:  datatypes.GetMissingPreviousComparisonMetrics(),
			ExpectedResponse: []string{"PASS - dummy_metric_name:percentile(90) fit static threshold of 1234.12 with value 0.00."},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			response := checkPerfSignature(test.PerfSignature, test.MetricsResponse)

			assert.Equal(t, test.ExpectedResponse, response.Response)
		})
	}

}

func TestPrintDeploymentTimestamps(t *testing.T) {
	type testDefs struct {
		Name   string
		Values []datatypes.Timestamps
	}

	tests := []testDefs{
		{
			Name: "One previous deployment",
			Values: []datatypes.Timestamps{
				{
					StartTime: 1574416227,
					EndTime:   1574419827,
				},
			},
		},
		{
			Name: "Two previous deployments",
			Values: []datatypes.Timestamps{
				{
					StartTime: 1574416227,
					EndTime:   1574419827,
				},
				{
					StartTime: 1574416227,
					EndTime:   1574419827,
				},
			},
		},
	}

	for _, test := range tests {
		printDeploymentTimestamps(test.Values)
	}
}

func printDeploymentTimestamps(timestamps []datatypes.Timestamps) {
	currentStartPretty := time.Unix(timestamps[0].StartTime/1000, 000)
	currentEndPretty := time.Unix(timestamps[0].EndTime/1000, 000)
	if len(timestamps) == 1 {
		deploymentText := fmt.Sprintf("Found current deployment from %v to %v.\n", currentStartPretty, currentEndPretty)
		logging.LogDebug(datatypes.Logging{Message: deploymentText})
	} else if len(timestamps) == 2 {
		previousStartPretty := time.Unix(timestamps[1].StartTime/1000, 000)
		previousEndPretty := time.Unix(timestamps[1].EndTime/1000, 000)
		deploymentText := fmt.Sprintf("Found previous deployment from %v to %v and current deployment from %v to %v.\n", previousStartPretty, previousEndPretty, currentStartPretty, currentEndPretty)
		logging.LogDebug(datatypes.Logging{Message: deploymentText})
	}
}
