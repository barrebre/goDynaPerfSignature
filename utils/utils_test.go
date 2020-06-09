package utils

import (
	"os"
	"testing"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
)

// TODO: Rewrite these tests so they validate the data provided is set in the config. There are no errors to be thrown
func TestGetConfig(t *testing.T) {
	type values struct {
		Config      datatypes.Config
		DtServerSet bool
		DtEnvSet    bool
	}

	type testDefs struct {
		Name          string
		Values        values
		ExpectPass    bool
		ExpectedError string
	}

	tests := []testDefs{
		testDefs{
			Name: "Env variables set",
			Values: values{
				Config: datatypes.Config{
					Env: "abcd1234",
				},
				DtServerSet: true,
				DtEnvSet:    true,
			},
			ExpectPass: true,
		},
		testDefs{
			Name: "Env variables not set",
			Values: values{
				Config: datatypes.Config{
					Env: "abcd1234",
				},
				DtServerSet: false,
				DtEnvSet:    false,
			},
			ExpectPass: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			// Clean up from previous tests
			os.Unsetenv("DT_SERVER")
			os.Unsetenv("DT_ENV")

			// Set envs if necessary
			if test.Values.DtServerSet {
				os.Setenv("DT_SERVER", "test1234.live.dynatrace.com")
			}

			if test.Values.DtEnvSet {
				os.Setenv("DT_ENV", "abc1234")
			}

			GetConfig()
		})
	}
}
