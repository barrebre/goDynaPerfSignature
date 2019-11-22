package performancesignature

import (
	"testing"

	"barrebre/goDynaPerfSignature/datatypes"

	"github.com/stretchr/testify/assert"
)

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
