package performancesignature

import (
	"encoding/json"
	"fmt"
	"net/http"

	"barrebre/goDynaPerfSignature/datatypes"
	"barrebre/goDynaPerfSignature/deployments"
	"barrebre/goDynaPerfSignature/metrics"

	"github.com/davecgh/go-spew/spew"
)

// ProcessRequest handles requests we receive to /performanceSignature
func ProcessRequest(w http.ResponseWriter, r *http.Request, config datatypes.Config, b []byte) (string, int, error) {
	// Read POST body params
	var performanceSignature datatypes.PerformanceSignature
	err := json.Unmarshal(b, &performanceSignature)
	if err != nil {
		return "", 400, err
	}

	// Verify all necessary params were sent
	err = checkParams(performanceSignature)
	if err != nil {
		fmt.Println("Encountered error at check params", err)
		return "", 400, err
	}

	// Query Dt for events
	timestamps, err := deployments.GetDeploymentTimestamps(config, performanceSignature.ServiceID, performanceSignature.APIToken)
	if err != nil {
		fmt.Printf("Encountered error gathering event timestamps: %v\n", err)
		return "", 503, fmt.Errorf("Encountered error gathering timestamps: %v", err)
	}
	if len(timestamps) < 2 {
		return fmt.Sprintln("There were not enough deployment events on the service to test."), 200, nil
	}

	metricsResponse, err := metrics.GetMetrics(config, performanceSignature, timestamps)
	if err != nil {
		fmt.Printf("Encountered error gathering metrics: %v\n", err)
		return "", 503, fmt.Errorf("Encountered error gathering metrics: %v", err)
	}

	successText, err := metrics.CompareMetrics(metricsResponse)
	if err != nil {
		responseText := fmt.Sprintf("Metric degradation found: %v\n", err)
		fmt.Printf(responseText)
		return "", 406, fmt.Errorf(responseText)
	}

	return successText, 0, nil
}

// Check the body params sent in with the request to ensure we have all the data we need to query Dt
func checkParams(p datatypes.PerformanceSignature) error {
	if p.APIToken == "" {
		return fmt.Errorf("No API Token found in object: %v", spew.Sdump(p))
	}

	if len(p.MetricIDs) == 0 {
		return fmt.Errorf("No MetricIDs found in object: %v", spew.Sdump(p))
	}

	if p.APIToken == "" {
		return fmt.Errorf("No APIToken found in object: %v", spew.Sdump(p))
	}

	return nil
}
