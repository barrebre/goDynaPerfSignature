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
		{
			Name:       "Invalid Endpoint",
			Values:     *passHTTPReq,
			ExpectPass: false,
		},
		{
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
	type TestValues struct {
		DeploymentEvents datatypes.DeploymentEvents
		EvaluationMins   int
	}

	type testDefs struct {
		Name           string
		Values         TestValues
		ExpectPass     bool
		ExpectedError  string
		ExpectedResult []datatypes.Timestamps
	}

	tests := []testDefs{
		{
			Name:           "No Deployment Events",
			ExpectPass:     true,
			ExpectedResult: []datatypes.Timestamps{},
		},
		{
			Name: "One Deployment Event",
			Values: TestValues{
				DeploymentEvents: datatypes.GetSingleEventDeploymentEvent(),
			},
			ExpectPass:     true,
			ExpectedResult: datatypes.GetSingleTimestamps(),
		},
		{
			Name: "One Deployment Event with Eval Time set",
			Values: TestValues{
				DeploymentEvents: datatypes.GetSingleEventDeploymentEvent(),
				EvaluationMins:   5,
			},
			ExpectPass:     false,
			ExpectedResult: datatypes.GetSingleTimestamps(),
		},
		{
			Name: "One Deployment Event with Eval Time set",
			Values: TestValues{
				DeploymentEvents: datatypes.GetSingleEventDeploymentEvent(),
				EvaluationMins:   5,
			},
			ExpectPass: true,
			ExpectedResult: []datatypes.Timestamps{
				{
					StartTime: 1234,
					EndTime:   301234,
				},
			},
		},
		{
			Name: "Two Deployment Events",
			Values: TestValues{
				DeploymentEvents: datatypes.GetMultipleEventDeploymentEvent(),
			},
			ExpectPass:     true,
			ExpectedResult: datatypes.GetMultipleTimestamps(),
		},
		{
			Name: "Two Deployment Events with Eval Time Set",
			Values: TestValues{
				DeploymentEvents: datatypes.GetMultipleEventDeploymentEvent(),
				EvaluationMins:   5,
			},
			ExpectPass: true,
			ExpectedResult: []datatypes.Timestamps{
				{
					StartTime: 1234,
					EndTime:   301234,
				},
				{
					StartTime: 1234,
					EndTime:   301234,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			ts, err := parseDeploymentTimestamps(test.Values.DeploymentEvents, test.Values.EvaluationMins)

			if test.ExpectPass == true {
				if test.ExpectedResult != nil {
					assert.Equal(t, test.ExpectedResult, ts)
				}
			} else {
				if test.ExpectedError != "" {
					assert.Equal(t, err.Error(), test.ExpectedError)
				} else {
					assert.NotEqual(t, test.ExpectedResult, ts)
				}
			}
		})
	}
}
