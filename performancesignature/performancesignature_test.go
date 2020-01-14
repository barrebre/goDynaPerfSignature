package performancesignature

import (
	"fmt"
	"testing"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/barrebre/goDynaPerfSignature/exampletestdata"

	"github.com/stretchr/testify/assert"
)

func TestBuildDeploymentRequest(t *testing.T) {
	type values struct {
		Conf      datatypes.Config
		ServiceID string
		APIToken  string
	}
	type testDefs struct {
		Name          string
		Values        values
		ExpectPass    bool
		ExpectedError string
	}

	tests := []testDefs{
		testDefs{
			Name: "Pass - valid params",
			Values: values{
				Conf: datatypes.Config{
					Env:    "",
					Server: "asdf1234.live.dynatrace.com",
				},
			},
			ExpectPass: true,
		},
		testDefs{
			Name: "Pass - valid params with env",
			Values: values{
				Conf: datatypes.Config{
					Env:    "qwer",
					Server: "asdf1234.live.dynatrace.com",
				},
			},
			ExpectPass: true,
		},
		testDefs{
			Name: "Fail - couldn't build HTTP request",
			Values: values{
				Conf: datatypes.Config{
					Env:    "",
					Server: `\`,
				},
			},
			ExpectPass:    false,
			ExpectedError: "parse https://\\/api/v1/events?eventType=CUSTOM_DEPLOYMENT&entityId=: invalid character \"\\\\\" in host name",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			_, err := buildDeploymentRequest(test.Values.Conf, test.Values.ServiceID, test.Values.APIToken)

			if test.ExpectPass == true {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, err.Error(), test.ExpectedError)
			}
		})
	}

}

func TestCheckParams(t *testing.T) {
	type testDefs struct {
		Name          string
		Values        datatypes.PerformanceSignature
		ExpectPass    bool
		ExpectedError string
	}

	tests := []testDefs{
		testDefs{
			Name: "Pass - all params included",
			Values: datatypes.PerformanceSignature{
				APIToken: "S2pMHW_FSlma-PPJIj3l5",
				Metrics: []datatypes.Metric{
					datatypes.Metric{
						ID: "builtin:service.response.time:(avg)",
					},
				},
				ServiceID: "SERVICE-5D4E743B2BF0CCF5",
			},
			ExpectPass: true,
		},
		testDefs{
			Name: "Fail - missing API token",
			Values: datatypes.PerformanceSignature{
				Metrics: []datatypes.Metric{
					datatypes.Metric{
						ID: "builtin:service.response.time:(avg)",
					},
				},
				ServiceID: "SERVICE-5D4E743B2BF0CCF5",
			},
			ExpectPass: false,
		},
		testDefs{
			Name: "Fail - missing metrics",
			Values: datatypes.PerformanceSignature{
				APIToken:  "S2pMHW_FSlma-PPJIj3l5",
				ServiceID: "SERVICE-5D4E743B2BF0CCF5",
			},
			ExpectPass: false,
		},
		testDefs{
			Name: "Pass - missing ServiceID",
			Values: datatypes.PerformanceSignature{
				APIToken: "S2pMHW_FSlma-PPJIj3l5",
				Metrics: []datatypes.Metric{
					datatypes.Metric{
						ID: "builtin:service.response.time:(avg)",
					},
				},
			},
			ExpectPass: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := checkParams(test.Values)

			if test.ExpectPass == true {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
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
			PerfSignature:   exampletestdata.GetValidDefaultPerformanceSignature(),
			MetricsResponse: exampletestdata.GetValidFailingComparisonMetrics(),
			ExpectPass:      false,
			ExpectedError:   "",
		},
		testDefs{
			Name:            "Valid Static Check Failing Data",
			PerfSignature:   exampletestdata.GetValidStaticPerformanceSignature(),
			MetricsResponse: exampletestdata.GetValidFailingComparisonMetrics(),
			ExpectPass:      false,
			ExpectedError:   "",
		},
		testDefs{
			Name:            "Valid Default Check Passing Data",
			PerfSignature:   exampletestdata.GetValidDefaultPerformanceSignature(),
			MetricsResponse: exampletestdata.GetValidPassingComparisonMetrics(),
			ExpectPass:      true,
			ExpectedError:   "",
		},
		testDefs{
			Name:            "Valid Default Check Passing Data",
			PerfSignature:   exampletestdata.GetValidDefaultPerformanceSignature(),
			MetricsResponse: exampletestdata.GetValidPassingComparisonMetrics(),
			ExpectPass:      true,
			ExpectedError:   "",
		},
		testDefs{
			Name:            "No Data Returned",
			PerfSignature:   exampletestdata.GetValidStaticPerformanceSignature(),
			MetricsResponse: exampletestdata.GetMissingComparisonMetrics(),
			ExpectPass:      false,
			ExpectedError:   "",
		},
		testDefs{
			Name:            "No Previous Deployment Data Returned - Default Check",
			PerfSignature:   exampletestdata.GetValidDefaultPerformanceSignature(),
			MetricsResponse: exampletestdata.GetMissingPreviousComparisonMetrics(),
			ExpectPass:      false,
			ExpectedError:   "",
		},
		testDefs{
			Name:            "No Previous Deployment Data Returned - Static Check",
			PerfSignature:   exampletestdata.GetValidStaticPerformanceSignature(),
			MetricsResponse: exampletestdata.GetMissingPreviousComparisonMetrics(),
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

func TestReadAndValidateParams(t *testing.T) {
	type testDefs struct {
		Name          string
		Values        []byte
		ExpectPass    bool
		ExpectedError string
	}

	passJSONstring := `{"APIToken":"S2pMHW_FSlma-PPJIj3l5","Metrics":[{"ID":"builtin:service.response.time:(avg)"},{"ID":"builtin:service.errors.total.rate:(avg)","StaticThreshold":1.0,"ValidationMethod":"static"}],"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`
	failJSONstring := `{"Metrics":[{"ID":"builtin:service.response.time:(avg)"},{"ID":"builtin:service.errors.total.rate:(avg)","StaticThreshold":1.0,"ValidationMethod":"static"}],"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`

	tests := []testDefs{
		testDefs{
			Name:       "Pass - valid JSON",
			Values:     []byte(passJSONstring),
			ExpectPass: true,
		},
		testDefs{
			Name:          "Fail - invalid JSON",
			Values:        []byte("Byte stream"),
			ExpectPass:    false,
			ExpectedError: `invalid character 'B' looking for beginning of value`,
		},
		testDefs{
			Name:          "Fail - invalid params",
			Values:        []byte(failJSONstring),
			ExpectPass:    false,
			ExpectedError: "No API Token found in object: (datatypes.PerformanceSignature) {\n APIToken: (string) \"\",\n Metrics: ([]datatypes.Metric) (len=2 cap=4) {\n  (datatypes.Metric) {\n   ID: (string) (len=35) \"builtin:service.response.time:(avg)\",\n   StaticThreshold: (float64) 0,\n   ValidationMethod: (string) \"\"\n  },\n  (datatypes.Metric) {\n   ID: (string) (len=39) \"builtin:service.errors.total.rate:(avg)\",\n   StaticThreshold: (float64) 1,\n   ValidationMethod: (string) (len=6) \"static\"\n  }\n },\n ServiceID: (string) (len=24) \"SERVICE-5D4E743B2BF0CCF5\"\n}\n",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			_, err := ReadAndValidateParams(test.Values)

			if test.ExpectPass == true {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, err.Error(), test.ExpectedError)
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
