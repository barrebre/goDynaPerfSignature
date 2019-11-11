package performancesignature

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

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
	// fmt.Printf("Received params: %v\n", spew.Sdump(performanceSignature))

	// Verify all necessary params were sent
	err = checkParams(performanceSignature)
	if err != nil {
		fmt.Println("Encountered error at check params", err)
		return "", 400, err
	}

	// Query Dt for events
	timestamps, err := deployments.GetDeploymentTimestamps(config, performanceSignature)
	if err != nil {
		fmt.Printf("Encountered error gathering event timestamps: %v\n", err)
		return "", 503, fmt.Errorf("Encountered error gathering timestamps: %v", err)
	}

	// Print out the deployment event(s) we found
	currentStartPretty := time.Unix(timestamps[0].StartTime/1000, 000)
	currentEndPretty := time.Unix(timestamps[0].EndTime/1000, 000)
	if len(timestamps) == 1 {
		fmt.Printf("Found current deployment from %v to %v.\n", currentStartPretty, currentEndPretty)
	} else if len(timestamps) == 2 {
		previousStartPretty := time.Unix(timestamps[1].StartTime/1000, 000)
		previousEndPretty := time.Unix(timestamps[1].EndTime/1000, 000)
		fmt.Printf("Found previous deployment from %v to %v and current deployment from %v to %v.\n", previousStartPretty, previousEndPretty, currentStartPretty, currentEndPretty)
	}

	// Get the requested metrics for the discovered timestamp(s)
	metricsResponse, err := metrics.GetMetrics(config, performanceSignature, timestamps)
	if err != nil {
		fmt.Printf("Encountered error gathering metrics: %v\n", err)
		return "", 503, fmt.Errorf("Encountered error gathering metrics: %v", err)
	}
	// fmt.Printf("Found metrics:\n%v\n", metricsResponse)

	// For each metric, perform its checks
	successText := ""
	for _, metric := range performanceSignature.Metrics {
		fmt.Printf("Looking at metric %v\n", metric)
		cleanMetricName := strings.ReplaceAll(metric.ID, "(", "")
		cleanMetricName = strings.ReplaceAll(cleanMetricName, ")", "")
		// fmt.Printf("Clean name is: %v\n", cleanMetricName)

		currentMetricValues := metricsResponse.CurrentMetrics.Metrics[cleanMetricName].MetricValues[0].Value
		previousMetricValues := metricsResponse.PreviousMetrics.Metrics[cleanMetricName].MetricValues[0].Value

		switch checkCounts := metric.ValidationMethod; checkCounts {
		case "static":
			// fmt.Println("Static check")
			response, err := metrics.CheckStaticThreshold(currentMetricValues, metric.StaticThreshold)
			if err != nil {
				degradationText := fmt.Sprintf("Metric degradation found: %v\n", err)
				fmt.Printf(degradationText)
				return "", 406, fmt.Errorf(degradationText)
			}
			successText += response
		default:
			// fmt.Println("Default check")
			response, err := metrics.CompareMetrics(currentMetricValues, previousMetricValues)
			if err != nil {
				degradationText := fmt.Sprintf("Metric degradation found: %v\n", err)
				fmt.Printf(degradationText)
				return "", 406, fmt.Errorf(degradationText)
			}
			successText += response
		}
	}
	return successText, 0, nil
}

// Check the required body params sent in with the request to ensure we have all the data we need to query Dt
func checkParams(p datatypes.PerformanceSignature) error {
	if p.APIToken == "" {
		return fmt.Errorf("No API Token found in object: %v", spew.Sdump(p))
	}

	if len(p.Metrics) == 0 {
		return fmt.Errorf("No MetricIDs found in object: %v", spew.Sdump(p))
	}

	if p.ServiceID == "" {
		return fmt.Errorf("No Services found in object: %v", spew.Sdump(p))
	}

	return nil
}
