package performancesignature

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
)

// Gets the deployment events from Dynatrace
func getDeploymentEvents(config datatypes.Config, ps datatypes.PerformanceSignature, req http.Request) (datatypes.DeploymentEvents, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Perform the request
	r, err := client.Do(&req)
	if err != nil {
		fmt.Printf("Error reading Deployment Event data from Dynatrace: %v", err)
		return datatypes.DeploymentEvents{}, err
	}
	// Check the status code
	if r.StatusCode != 200 {
		fmt.Printf("Invalid status code from Dynatrace: %v.\n", r.StatusCode)
		return datatypes.DeploymentEvents{}, fmt.Errorf("Invalid status code from Dynatrace: %v", r.StatusCode)
	}

	// Read in the body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	// Try to parse the response into DeploymentEvents
	var deploymentEvents datatypes.DeploymentEvents
	err = json.Unmarshal(b, &deploymentEvents)
	if err != nil {
		return datatypes.DeploymentEvents{}, err
	}

	return deploymentEvents, nil
}

// Parses Dynatrace Deployment Events for their timestamps
func parseDeploymentTimestamps(d datatypes.DeploymentEvents) ([]datatypes.Timestamps, error) {
	switch deploymentCount := len(d.Events); deploymentCount {
	// If there are no deployment events previously, we can still perform static checks
	case 0:
		fmt.Println("There haven't been enough deployment events. Auto-passing")
		return []datatypes.Timestamps{}, nil
	// If there is only one deployment event, we can still perform static checks
	case 1:
		var deploymentTimestamp = []datatypes.Timestamps{
			datatypes.Timestamps{
				StartTime: d.Events[0].StartTime,
				EndTime:   d.Events[0].EndTime,
			},
		}
		return deploymentTimestamp, nil
	// If there are two deployment events, we can perform all types of checks
	case 2:
		var deploymentTimestamps = []datatypes.Timestamps{
			datatypes.Timestamps{
				StartTime: d.Events[0].StartTime,
				EndTime:   d.Events[0].EndTime,
			},
			datatypes.Timestamps{
				StartTime: d.Events[1].StartTime,
				EndTime:   d.Events[1].EndTime,
			},
		}
		return deploymentTimestamps, nil
	}

	return []datatypes.Timestamps{}, fmt.Errorf("Wasn't able to read deployments from Dynatrace")
}
