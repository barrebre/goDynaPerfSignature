package performancesignature

import (
	"fmt"
	"testing"

	"github.com/barrebre/goDynaPerfSignature/datatypes"

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
		testDefs{
			Name: "Pass - valid params",
			Values: datatypes.PerformanceSignature{
				DTEnv:    "",
				DTServer: "asdf1234.live.dynatrace.com",
			},
			ExpectPass: true,
		},
		testDefs{
			Name: "Pass - valid params with env",
			Values: datatypes.PerformanceSignature{
				DTEnv:    "POWEIFJPIOJAPSOIJ",
				DTServer: "asdf1234.live.dynatrace.com",
			},
			ExpectPass: true,
		},
		testDefs{
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
		Name            string
		PerfSignature   datatypes.PerformanceSignature
		MetricsResponse datatypes.ComparisonMetrics
		ExpectPass      bool
		ExpectedError   string
	}

	tests := []testDefs{
		testDefs{
			Name:            "Valid Default Check Failing Data",
			PerfSignature:   datatypes.GetValidDefaultPerformanceSignature(),
			MetricsResponse: datatypes.GetValidFailingComparisonMetrics(),
			ExpectPass:      false,
			ExpectedError:   "",
		},
		testDefs{
			Name:            "Valid Static Check Failing Data",
			PerfSignature:   datatypes.GetValidStaticPerformanceSignature(),
			MetricsResponse: datatypes.GetValidFailingComparisonMetrics(),
			ExpectPass:      false,
			ExpectedError:   "",
		},
		testDefs{
			Name:            "Valid Default Check Passing Data",
			PerfSignature:   datatypes.GetValidDefaultPerformanceSignature(),
			MetricsResponse: datatypes.GetValidPassingComparisonMetrics(),
			ExpectPass:      true,
			ExpectedError:   "",
		},
		testDefs{
			Name:            "Valid Default Check Passing Data",
			PerfSignature:   datatypes.GetValidDefaultPerformanceSignature(),
			MetricsResponse: datatypes.GetValidPassingComparisonMetrics(),
			ExpectPass:      true,
			ExpectedError:   "",
		},
		testDefs{
			Name:            "No Data Returned",
			PerfSignature:   datatypes.GetValidStaticPerformanceSignature(),
			MetricsResponse: datatypes.GetMissingComparisonMetrics(),
			ExpectPass:      false,
			ExpectedError:   "",
		},
		testDefs{
			Name:            "No Previous Deployment Data Returned - Default Check",
			PerfSignature:   datatypes.GetValidDefaultPerformanceSignature(),
			MetricsResponse: datatypes.GetMissingPreviousComparisonMetrics(),
			ExpectPass:      false,
			ExpectedError:   "",
		},
		testDefs{
			Name:            "No Previous Deployment Data Returned - Static Check",
			PerfSignature:   datatypes.GetValidStaticPerformanceSignature(),
			MetricsResponse: datatypes.GetMissingPreviousComparisonMetrics(),
			ExpectPass:      true,
			ExpectedError:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			successText, responseCode, err := checkPerfSignature(test.PerfSignature, test.MetricsResponse)

			fmt.Println(successText, responseCode, err)
			if test.ExpectPass == true {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}

}

func TestPrintDeploymentTimestamps(t *testing.T) {
	type testDefs struct {
		Name   string
		Values []datatypes.Timestamps
	}

	tests := []testDefs{
		testDefs{
			Name: "One previous deployment",
			Values: []datatypes.Timestamps{
				datatypes.Timestamps{
					StartTime: 1574416227,
					EndTime:   1574419827,
				},
			},
		},
		testDefs{
			Name: "Two previous deployments",
			Values: []datatypes.Timestamps{
				datatypes.Timestamps{
					StartTime: 1574416227,
					EndTime:   1574419827,
				},
				datatypes.Timestamps{
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
