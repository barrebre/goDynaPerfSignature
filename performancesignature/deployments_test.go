package performancesignature

import (
	"net/http"
	"testing"

	"github.com/barrebre/goDynaPerfSignature/datatypes"

	"github.com/stretchr/testify/assert"
)

func TestGetDeploymentEvents(t *testing.T) {
	type testDefs struct {
		Name           string
		Values         http.Request
		ExpectPass     bool
		ExpectedError  string
		ExpectedResult datatypes.DeploymentEvents
	}

	url := "https://www.google.com"
	passHTTPReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		assert.FailNow(t, "Couldn't query endpoint")
	}

	tests := []testDefs{
		testDefs{
			Name:       "Invalid Endpoint",
			Values:     *passHTTPReq,
			ExpectPass: false,
		},
		testDefs{
			Name:       "Invalid Request",
			Values:     http.Request{},
			ExpectPass: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			_, err := getDeploymentEvents(test.Values)

			if test.ExpectPass == true {
				assert.NoError(t, err)
			} else {
				if test.ExpectedError != "" {
					assert.Equal(t, err.Error(), test.ExpectedError)
				} else {
					assert.Error(t, err)
				}
			}
		})
	}
}

func TestParseDeploymentTimestamps(t *testing.T) {
	type testDefs struct {
		Name           string
		Values         datatypes.DeploymentEvents
		ExpectPass     bool
		ExpectedError  string
		ExpectedResult []datatypes.Timestamps
	}

	tests := []testDefs{
		testDefs{
			Name:           "No Deployment Events",
			ExpectPass:     true,
			ExpectedResult: []datatypes.Timestamps{},
		},
		testDefs{
			Name:           "One Deployment Event",
			Values:         datatypes.GetSingleEventDeploymentEvent(),
			ExpectPass:     true,
			ExpectedResult: datatypes.GetSingleTimestamps(),
		},
		testDefs{
			Name:           "Two Deployment Events",
			Values:         datatypes.GetMultipleEventDeploymentEvent(),
			ExpectPass:     true,
			ExpectedResult: datatypes.GetMultipleTimestamps(),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ts, err := parseDeploymentTimestamps(test.Values)

			if test.ExpectPass == true {
				// assert.NoError(t, err)
				if test.ExpectedResult != nil {
					assert.Equal(t, test.ExpectedResult, ts)
				}
			} else {
				assert.Equal(t, err.Error(), test.ExpectedError)
			}
		})
	}
}
