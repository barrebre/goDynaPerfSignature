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
		Name          string
		Values        values
		ExpectPass    bool
		ExpectedError string
	}

	validJSON := `{"DTServer":"testserver","DTEnv":"testEnv","APIToken":"S2pMHW_FSlma-PPJIj3l5","Metrics":[{"ID":"builtin:service.response.time:(avg)"},{"ID":"builtin:service.errors.total.rate:(avg)","StaticThreshold":1.0,"ValidationMethod":"static"}],"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`
	invalidJSONNoAPIToken := `{"DTServer":"testserver","DTEnv":"testEnv","Metrics":[{"ID":"builtin:service.response.time:(avg)"},{"ID":"builtin:service.errors.total.rate:(avg)","StaticThreshold":1.0,"ValidationMethod":"static"}],"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`
	invalidJSONNoServer := `{"DTEnv":"testEnv","APIToken":"S2pMHW_FSlma-PPJIj3l5","Metrics":[{"ID":"builtin:service.response.time:(avg)"},{"ID":"builtin:service.errors.total.rate:(avg)","StaticThreshold":1.0,"ValidationMethod":"static"}],"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`
	invalidJSONNoMetrics := `{"DTServer":"testserver","DTEnv":"testEnv","APIToken":"S2pMHW_FSlma-PPJIj3l5","ServiceID":"SERVICE-5D4E743B2BF0CCF5"}`
	invalidJSONNoServices := `{"DTServer":"testserver","DTEnv":"testEnv","APIToken":"S2pMHW_FSlma-PPJIj3l5","Metrics":[{"ID":"builtin:service.response.time:(avg)"},{"ID":"builtin:service.errors.total.rate:(avg)","StaticThreshold":1.0,"ValidationMethod":"static"}]}`

	tests := []testDefs{
		testDefs{
			Name: "Pass - valid JSON",
			Values: values{
				APIString: []byte(validJSON),
				Config:    datatypes.GetConfiguredConfig(),
			},
			ExpectPass: true,
		},
		testDefs{
			Name: "Fail - no DTServer provided",
			Values: values{
				APIString: []byte(invalidJSONNoServer),
				Config:    datatypes.Config{},
			},
			ExpectPass:    false,
			ExpectedError: "There is no DT_SERVER env variable configured and no DTServer was passed with the POST",
		},
		testDefs{
			Name: "Fail - no APIToken provided",
			Values: values{
				APIString: []byte(invalidJSONNoAPIToken),
				Config:    datatypes.Config{},
			},
			ExpectPass:    false,
			ExpectedError: "There is no DT_API_TOKEN env variable configured and no APIToken was passed with the POST",
		},
		testDefs{
			Name: "Fail - no metrics provided",
			Values: values{
				APIString: []byte(invalidJSONNoMetrics),
				Config:    datatypes.Config{},
			},
			ExpectPass:    false,
			ExpectedError: "No Metrics passed with the POST",
		},
		testDefs{
			Name: "Fail - no services provided",
			Values: values{
				APIString: []byte(invalidJSONNoServices),
				Config:    datatypes.Config{},
			},
			ExpectPass:    false,
			ExpectedError: "No ServiceID passed with the POST",
		},
		testDefs{
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
			_, err := ReadAndValidateParams(test.Values.APIString, test.Values.Config)

			if test.ExpectPass == true {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, test.ExpectedError, err.Error())
			}
		})
	}
}
