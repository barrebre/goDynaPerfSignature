package deployments

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"barrebre/goDynaPerfSignature/datatypes"
)

// GetDeploymentTimestamps Gets the timestamps from the two most recent Dynatrace Deployment Events
func GetDeploymentTimestamps(config datatypes.Config, serviceID string, apiToken string) (timestamps []datatypes.Timestamps, err error) {
	// Build the URL
	url := fmt.Sprintf("https://%v/e/%v/api/v1/events?eventType=CUSTOM_DEPLOYMENT&entityId=%v", config.Server, config.Env, serviceID)
	// fmt.Printf("Made URL: %v\n", url)

	// Build the request object
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request handler: %v", err)
		return make([]datatypes.Timestamps, 0), err
	}

	apiTokenField := fmt.Sprintf("Api-Token %v", apiToken)
	req.Header.Add("Authorization", apiTokenField)
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Perform the request
	r, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error reading Deployment Event data from Dynatrace: %v", err)
		return make([]datatypes.Timestamps, 0), err
	}
	// Check the status code
	if r.StatusCode != 200 {
		fmt.Printf("Invalid status code from Dynatrace: %v", r.StatusCode)
		return make([]datatypes.Timestamps, 0), fmt.Errorf("Invalid status code from Dynatrace: %v", r.StatusCode)
	}

	// Read in the body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	// Try to parse the response into DeploymentEvents
	var deploymentEvents datatypes.DeploymentEvents
	err = json.Unmarshal(b, &deploymentEvents)
	if err != nil {
		return nil, err
	}

	// If there is only one deployment event (maybe first time, or only one found), return pass
	if len(deploymentEvents.Events) < 2 {
		fmt.Println("There haven't been enough deployment events. Auto-passing")
		return make([]datatypes.Timestamps, 0), nil
	}

	// Grab the two latest event timestamps
	var timestampsToCompare = []datatypes.Timestamps{
		datatypes.Timestamps{
			StartTime: deploymentEvents.Events[0].StartTime,
			EndTime:   deploymentEvents.Events[0].EndTime,
		},
		datatypes.Timestamps{
			StartTime: deploymentEvents.Events[1].StartTime,
			EndTime:   deploymentEvents.Events[1].EndTime,
		},
	}

	currentStartPretty := time.Unix(timestampsToCompare[0].StartTime/1000, 000)
	currentEndPretty := time.Unix(timestampsToCompare[0].EndTime/1000, 000)
	previousStartPretty := time.Unix(timestampsToCompare[1].StartTime/1000, 000)
	previousEndPretty := time.Unix(timestampsToCompare[1].EndTime/1000, 000)
	fmt.Printf("Found previous deployment from %v to %v and current deployment from %v to %v.\n", previousStartPretty, previousEndPretty, currentStartPretty, currentEndPretty)

	return timestampsToCompare, nil
}
