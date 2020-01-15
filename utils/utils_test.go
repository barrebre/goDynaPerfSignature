package utils

import (
	"os"
	"testing"

	"github.com/barrebre/goDynaPerfSignature/datatypes"

	"github.com/stretchr/testify/assert"
)

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
			ExpectPass:    false,
			ExpectedError: "Error finding the DT_SERVER in the env",
		},
		testDefs{
			Name: "Env variables not set",
			Values: values{
				Config: datatypes.Config{
					Env: "abcd1234",
				},
				DtServerSet: true,
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

			// Get Config
			_, err := GetConfig()

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
