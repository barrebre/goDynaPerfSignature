package performancesignature

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/barrebre/goDynaPerfSignature/logging"
)

// Gets the deployment events from Dynatrace
func getDeploymentEvents(req http.Request) (datatypes.DeploymentEvents, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Perform the request
	r, err := client.Do(&req)
	if err != nil {
		logging.LogInfo(datatypes.Logging{Message: fmt.Sprintf("Error reading Deployment Event data from Dynatrace: %v", err)})
		return datatypes.DeploymentEvents{}, err
	}

	// Read in the body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Could not read response body from Dynatrace: %v", err.Error())})
		return datatypes.DeploymentEvents{}, fmt.Errorf("could not read response body from Dynatrace: %v", err.Error())
	}

	// Check the status code
	if r.StatusCode != 200 {
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Invalid status code from Dynatrace: %v. Body message is '%v'\n", r.StatusCode, string(b))})
		return datatypes.DeploymentEvents{}, fmt.Errorf("invalid status code from Dynatrace: %v", r.StatusCode)
	}

	// Try to parse the response into DeploymentEvents
	var deploymentEvents datatypes.DeploymentEvents
	err = json.Unmarshal(b, &deploymentEvents)
	if err != nil {
		return datatypes.DeploymentEvents{}, err
	}

	return deploymentEvents, nil
}

// Parses Dynatrace Deployment Events for their timestamps
func parseDeploymentTimestamps(d datatypes.DeploymentEvents, mins int) ([]datatypes.Timestamps, error) {
	eventsFound := len(d.Events)

	// If there are no deployment events previously, we can still perform static checks
	if eventsFound == 0 {
		logging.LogInfo(datatypes.Logging{Message: "There haven't been enough deployment events. Auto-passing"})
		return []datatypes.Timestamps{}, nil
		// If there is only one deployment event, we can still perform static checks
	} else if eventsFound == 1 {
		var deploymentTimestamp []datatypes.Timestamps

		// If there is no evaluation timeframe supplied
		if mins < 1 {
			deploymentTimestamp = []datatypes.Timestamps{
				{
					StartTime: d.Events[0].StartTime,
					EndTime:   d.Events[0].EndTime,
				},
			}
		} else {
			microMins := int64(mins * 60000)
			deploymentTimestamp = []datatypes.Timestamps{
				{
					StartTime: d.Events[0].StartTime,
					EndTime:   d.Events[0].StartTime + microMins,
				},
			}
		}
		return deploymentTimestamp, nil
		// If there are two deployment events, we can perform all types of checks
	} else if eventsFound >= 2 {
		var deploymentTimestamps []datatypes.Timestamps

		// If there is no evaluation timeframe supplied
		if mins < 1 {
			deploymentTimestamps = []datatypes.Timestamps{
				{
					StartTime: d.Events[0].StartTime,
					EndTime:   d.Events[0].EndTime,
				},
				{
					StartTime: d.Events[1].StartTime,
					EndTime:   d.Events[1].EndTime,
				},
			}
		} else {
			microMins := int64(mins * 60000)
			deploymentTimestamps = []datatypes.Timestamps{
				{
					StartTime: d.Events[0].StartTime,
					EndTime:   d.Events[0].StartTime + microMins,
				},
				{
					StartTime: d.Events[1].StartTime,
					EndTime:   d.Events[1].StartTime + microMins,
				},
			}
		}
		return deploymentTimestamps, nil
	}

	return []datatypes.Timestamps{}, fmt.Errorf("wasn't able to read deployments from Dynatrace")
}
