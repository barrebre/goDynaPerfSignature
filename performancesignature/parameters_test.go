package performancesignature

import (
	"testing"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/stretchr/testify/assert"
)

func TestReadAndValidateParams(t *testing.T) {
	type values struct {
		APIString []byte
		Config    datatypes.Config
	}

	type testDefs struct {
		Name            string
		Values          values
		ExpectPass      bool
		ExpectedError   string
		ExpectedPerfsig datatypes.PerformanceSignature
	}

	validJSON := `{"DTServer":"testserver","DTEnv":"testEnv","EventAge":180,"APIToken":"S2pMHW_FSlma-PPJIj3l5","PSMetrics":{"builtin:service.response.time:(avg)":{},"builtin:service.errors.total.rate:(avg)":{"StaticThreshold":1.0,"ValidationMethod":"static"}},"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`
	invalidJSONNoAPIToken := `{"DTServer":"testserver","DTEnv":"testEnv","PSMetrics":{"builtin:service.response.time:(avg)":{},"builtin:service.errors.total.rate:(avg)":{"StaticThreshold":1.0,"ValidationMethod":"static"}},"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`
	invalidJSONNoServer := `{"DTEnv":"testEnv","APIToken":"S2pMHW_FSlma-PPJIj3l5","PSMetrics":{"builtin:service.response.time:(avg)":{},"builtin:service.errors.total.rate:(avg)":{"StaticThreshold":1.0,"ValidationMethod":"static"}},"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`
	invalidJSONNoMetrics := `{"DTServer":"testserver","DTEnv":"testEnv","APIToken":"S2pMHW_FSlma-PPJIj3l5","ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`
	invalidJSONNoServices := `{"DTServer":"testserver","DTEnv":"testEnv","APIToken":"S2pMHW_FSlma-PPJIj3l5","PSMetrics":{"builtin:service.response.time:(avg)":{},"builtin:service.errors.total.rate:(avg)":{"StaticThreshold":1.0,"ValidationMethod":"static"}}}`

	tests := []testDefs{
		{
			Name: "Pass - valid JSON",
			Values: values{
				APIString: []byte(validJSON),
				Config:    datatypes.GetConfiguredConfig(),
			},
			ExpectedPerfsig: datatypes.PerformanceSignature{
				APIToken:       "S2pMHW_FSlma-PPJIj3l5",
				DTEnv:          "testEnv",
				DTServer:       "testserver",
				EvaluationMins: 0,
				EventAge:       calculateAgeEpoch(180),
				PSMetrics: map[string]datatypes.PSMetric{
					"builtin:service.response.time:(avg)": {
						RelativeThreshold: 0,
						StaticThreshold:   0,
						ValidationMethod:  "",
					},
					"builtin:service.errors.total.rate:(avg)": {
						RelativeThreshold: 0,
						StaticThreshold:   1,
						ValidationMethod:  "static",
					},
				},
				ServiceID: "SERVICE-5D4E743B2BF0CCF5",
			},
			ExpectPass: true,
		},
		{
			Name: "Fail - no DTServer provided",
			Values: values{
				APIString: []byte(invalidJSONNoServer),
				Config:    datatypes.Config{},
			},
			ExpectPass:    false,
			ExpectedError: "checkParams - Couldn't validate parameters: there is no DT_SERVER env variable configured and no DTServer was passed with the POST",
		},
		{
			Name: "Fail - no APIToken provided",
			Values: values{
				APIString: []byte(invalidJSONNoAPIToken),
				Config:    datatypes.Config{},
			},
			ExpectPass:    false,
			ExpectedError: "checkParams - Couldn't validate parameters: there is no DT_API_TOKEN env variable configured and no APIToken was passed with the POST",
		},
		{
			Name: "PASS - default APIToken configured",
			Values: values{
				APIString: []byte(invalidJSONNoAPIToken),
				Config: datatypes.Config{
					APIToken: "asdf",
				},
			},
			ExpectedPerfsig: datatypes.PerformanceSignature{
				APIToken:       "asdf",
				DTEnv:          "testEnv",
				DTServer:       "testserver",
				EvaluationMins: 0,
				PSMetrics: map[string]datatypes.PSMetric{
					"builtin:service.response.time:(avg)": {
						RelativeThreshold: 0,
						StaticThreshold:   0,
						ValidationMethod:  "",
					},
					"builtin:service.errors.total.rate:(avg)": {
						RelativeThreshold: 0,
						StaticThreshold:   1,
						ValidationMethod:  "static",
					},
				},
				ServiceID: "SERVICE-5D4E743B2BF0CCF5",
			},
			ExpectPass: true,
		},
		{
			Name: "Fail - no metrics provided",
			Values: values{
				APIString: []byte(invalidJSONNoMetrics),
				Config:    datatypes.Config{},
			},
			ExpectPass:    false,
			ExpectedError: "checkParams - Couldn't validate parameters: no Metrics passed with the POST",
		},
		{
			Name: "Fail - no services provided",
			Values: values{
				APIString: []byte(invalidJSONNoServices),
				Config:    datatypes.Config{},
			},
			ExpectPass:    false,
			ExpectedError: "checkParams - Couldn't validate parameters: no ServiceID passed with the POST",
		},
		{
			Name: "Fail - invalid JSON",
			Values: values{
				APIString: []byte("Byte stream"),
			},
			ExpectPass:    false,
			ExpectedError: `invalid character 'B' looking for beginning of value`,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			perfsig, err := ReadAndValidateParams(test.Values.APIString, test.Values.Config)

			if test.ExpectPass == true {
				assert.NoError(t, err)
				assert.Equal(t, test.ExpectedPerfsig, perfsig)
			} else {
				assert.Equal(t, test.ExpectedError, err.Error())
			}
		})
	}
}
