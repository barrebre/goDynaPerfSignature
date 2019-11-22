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
func GetDeploymentTimestamps(config datatypes.Config, ps datatypes.PerformanceSignature) (timestamps []datatypes.Timestamps, err error) {
	// Build the URL
	var url string

	if config.Env == "" {
		url = fmt.Sprintf("https://%v/api/v1/events?eventType=CUSTOM_DEPLOYMENT&entityId=%v", config.Server, ps.ServiceID)
	} else {
		url = fmt.Sprintf("https://%v/e/%v/api/v1/events?eventType=CUSTOM_DEPLOYMENT&entityId=%v", config.Server, config.Env, ps.ServiceID)
	}
	// fmt.Printf("Made URL: %v\n", url)

	// Build the request object
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request handler: %v", err)
		return make([]datatypes.Timestamps, 0), err
	}

	apiTokenField := fmt.Sprintf("Api-Token %v", ps.APIToken)
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
		fmt.Printf("Invalid status code from Dynatrace: %v.\n", r.StatusCode)
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

	switch deploymentCount := len(deploymentEvents.Events); deploymentCount {
	// If there are no deployment events previously, we can still perform static checks
	case 0:
		fmt.Println("There haven't been enough deployment events. Auto-passing")
		return make([]datatypes.Timestamps, 0), fmt.Errorf("No deployment events found")
	// If there is only one deployment event, we can still perform static checks
	case 1:
		var deploymentTimestamp = []datatypes.Timestamps{
			datatypes.Timestamps{
				StartTime: deploymentEvents.Events[0].StartTime,
				EndTime:   deploymentEvents.Events[0].EndTime,
			},
		}
		return deploymentTimestamp, nil
	// If there are two deployment events, we can perform all types of checks
	case 2:
		var deploymentTimestamps = []datatypes.Timestamps{
			datatypes.Timestamps{
				StartTime: deploymentEvents.Events[0].StartTime,
				EndTime:   deploymentEvents.Events[0].EndTime,
			},
			datatypes.Timestamps{
				StartTime: deploymentEvents.Events[1].StartTime,
				EndTime:   deploymentEvents.Events[1].EndTime,
			},
		}
		return deploymentTimestamps, nil
	}

	return []datatypes.Timestamps{}, fmt.Errorf("Wasn't able to read deployments from Dynatrace")
}
