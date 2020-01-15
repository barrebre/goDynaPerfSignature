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
				DTEnv:    "",
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

func TestReadAndValidateParams(t *testing.T) {
	type values struct {
		APIString []byte
		Config    datatypes.Config
	}

	type testDefs struct {
		Name          string
		Values        values
		ExpectPass    bool
		ExpectedError string
	}

	validJSONstringWithServer := `{"DTServer":"testserver","DTEnv":"testEnv","APIToken":"S2pMHW_FSlma-PPJIj3l5","Metrics":[{"ID":"builtin:service.response.time:(avg)"},{"ID":"builtin:service.errors.total.rate:(avg)","StaticThreshold":1.0,"ValidationMethod":"static"}],"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`
	validJSONstring := `{"APIToken":"S2pMHW_FSlma-PPJIj3l5","Metrics":[{"ID":"builtin:service.response.time:(avg)"},{"ID":"builtin:service.errors.total.rate:(avg)","StaticThreshold":1.0,"ValidationMethod":"static"}],"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`
	invalidJSONNoMetrics := `{"APIToken":"S2pMHW_FSlma-PPJIj3l5","ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`
	invalidJSONNoServices := `{"APIToken":"S2pMHW_FSlma-PPJIj3l5","Metrics":[{"ID":"builtin:service.response.time:(avg)"},{"ID":"builtin:service.errors.total.rate:(avg)","StaticThreshold":1.0,"ValidationMethod":"static"}]}`
	invalidJSONstring := `{"Metrics":[{"ID":"builtin:service.response.time:(avg)"},{"ID":"builtin:service.errors.total.rate:(avg)","StaticThreshold":1.0,"ValidationMethod":"static"}],"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`

	tests := []testDefs{
		testDefs{
			Name: "Pass - valid JSON, no server included",
			Values: values{
				APIString: []byte(validJSONstring),
				Config:    datatypes.GetConfiguredConfig(),
			},
			ExpectPass: true,
		},
		testDefs{
			Name: "Pass - valid JSON with server included",
			Values: values{
				APIString: []byte(validJSONstringWithServer),
				Config:    datatypes.Config{},
			},
			ExpectPass: true,
		},
		testDefs{
			Name: "Fail - no server provided",
			Values: values{
				APIString: []byte(validJSONstring),
				Config:    datatypes.Config{},
			},
			ExpectPass:    false,
			ExpectedError: "There is no default server configured and no server was passed with the POST: (datatypes.PerformanceSignature) {\n APIToken: (string) (len=21) \"S2pMHW_FSlma-PPJIj3l5\",\n DTEnv: (string) \"\",\n DTServer: (string) \"\",\n Metrics: ([]datatypes.Metric) (len=2 cap=4) {\n  (datatypes.Metric) {\n   ID: (string) (len=35) \"builtin:service.response.time:(avg)\",\n   StaticThreshold: (float64) 0,\n   ValidationMethod: (string) \"\"\n  },\n  (datatypes.Metric) {\n   ID: (string) (len=39) \"builtin:service.errors.total.rate:(avg)\",\n   StaticThreshold: (float64) 1,\n   ValidationMethod: (string) (len=6) \"static\"\n  }\n },\n ServiceID: (string) (len=24) \"SERVICE-5D4E743B2BF0CCF5\"\n}\n",
		},
		testDefs{
			Name: "Fail - no metrics provided",
			Values: values{
				APIString: []byte(invalidJSONNoMetrics),
				Config:    datatypes.Config{},
			},
			ExpectPass:    false,
			ExpectedError: "No MetricIDs found in object: (datatypes.PerformanceSignature) {\n APIToken: (string) (len=21) \"S2pMHW_FSlma-PPJIj3l5\",\n DTEnv: (string) \"\",\n DTServer: (string) \"\",\n Metrics: ([]datatypes.Metric) <nil>,\n ServiceID: (string) (len=24) \"SERVICE-5D4E743B2BF0CCF5\"\n}\n",
		},
		testDefs{
			Name: "Fail - no services provided",
			Values: values{
				APIString: []byte(invalidJSONNoServices),
				Config:    datatypes.Config{},
			},
			ExpectPass:    false,
			ExpectedError: "No Services found in object: (datatypes.PerformanceSignature) {\n APIToken: (string) (len=21) \"S2pMHW_FSlma-PPJIj3l5\",\n DTEnv: (string) \"\",\n DTServer: (string) \"\",\n Metrics: ([]datatypes.Metric) (len=2 cap=4) {\n  (datatypes.Metric) {\n   ID: (string) (len=35) \"builtin:service.response.time:(avg)\",\n   StaticThreshold: (float64) 0,\n   ValidationMethod: (string) \"\"\n  },\n  (datatypes.Metric) {\n   ID: (string) (len=39) \"builtin:service.errors.total.rate:(avg)\",\n   StaticThreshold: (float64) 1,\n   ValidationMethod: (string) (len=6) \"static\"\n  }\n },\n ServiceID: (string) \"\"\n}\n",
		},
		testDefs{
			Name: "Fail - invalid JSON",
			Values: values{
				APIString: []byte("Byte stream"),
			},
			ExpectPass:    false,
			ExpectedError: `invalid character 'B' looking for beginning of value`,
		},
		testDefs{
			Name: "Fail - invalid params",
			Values: values{
				APIString: []byte(invalidJSONstring),
			},
			ExpectPass:    false,
			ExpectedError: "No API Token found in object: (datatypes.PerformanceSignature) {\n APIToken: (string) \"\",\n DTEnv: (string) \"\",\n DTServer: (string) \"\",\n Metrics: ([]datatypes.Metric) (len=2 cap=4) {\n  (datatypes.Metric) {\n   ID: (string) (len=35) \"builtin:service.response.time:(avg)\",\n   StaticThreshold: (float64) 0,\n   ValidationMethod: (string) \"\"\n  },\n  (datatypes.Metric) {\n   ID: (string) (len=39) \"builtin:service.errors.total.rate:(avg)\",\n   StaticThreshold: (float64) 1,\n   ValidationMethod: (string) (len=6) \"static\"\n  }\n },\n ServiceID: (string) (len=24) \"SERVICE-5D4E743B2BF0CCF5\"\n}\n",
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			_, err := ReadAndValidateParams(test.Values.APIString, test.Values.Config)

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
