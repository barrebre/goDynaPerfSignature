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

	// will be used in future for info logging
	// printDeploymentTimestamps(timestamps)

	// Get the requested metrics for the discovered timestamp(s)
	metricsResponse, err := metrics.GetMetrics(config, performanceSignature, timestamps)
	if err != nil {
		fmt.Printf("Encountered error gathering metrics: %v\n", err)
		return "", 503, fmt.Errorf("Encountered error gathering metrics: %v", err)
	}
	// will be used in future for debug logging
	// fmt.Printf("Found metrics:\n%v\n", metricsResponse)

	successText, responseCode, err := checkPerfSignature(performanceSignature, metricsResponse)
	if err != nil {
		fmt.Printf("Error occurred when checking performance signature: %v", err)
		return "", responseCode, fmt.Errorf("Error occurred when checking performance signature: %v", err)
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

func printDeploymentTimestamps(timestamps []datatypes.Timestamps) {
	currentStartPretty := time.Unix(timestamps[0].StartTime/1000, 000)
	currentEndPretty := time.Unix(timestamps[0].EndTime/1000, 000)
	if len(timestamps) == 1 {
		fmt.Printf("Found current deployment from %v to %v.\n", currentStartPretty, currentEndPretty)
	} else if len(timestamps) == 2 {
		previousStartPretty := time.Unix(timestamps[1].StartTime/1000, 000)
		previousEndPretty := time.Unix(timestamps[1].EndTime/1000, 000)
		fmt.Printf("Found previous deployment from %v to %v and current deployment from %v to %v.\n", previousStartPretty, previousEndPretty, currentStartPretty, currentEndPretty)
	}
}

// For each metric, perform its checks
func checkPerfSignature(performanceSignature datatypes.PerformanceSignature, metricsResponse datatypes.ComparisonMetrics) (string, int, error) {
	successText := ""
	for _, metric := range performanceSignature.Metrics {
		// will be used in future for debug logging
		// fmt.Printf("Looking at metric %v\n", metric)
		cleanMetricName := strings.ReplaceAll(metric.ID, "(", "")
		cleanMetricName = strings.ReplaceAll(cleanMetricName, ")", "")
		// will be used in future for debug logging
		// fmt.Printf("Clean name is: %v\n", cleanMetricName)

		if len(metricsResponse.CurrentMetrics.Metrics[cleanMetricName].MetricValues) < 1 {
			return "", 400, fmt.Errorf("There were no current metrics found for %v", cleanMetricName)
		}
		currentMetricValues := metricsResponse.CurrentMetrics.Metrics[cleanMetricName].MetricValues[0].Value

		// This is only an issue if trying a comparison
		canCompare := true
		if len(metricsResponse.PreviousMetrics.Metrics[cleanMetricName].MetricValues) < 1 {
			canCompare = false
		}
		previousMetricValues := metricsResponse.PreviousMetrics.Metrics[cleanMetricName].MetricValues[0].Value

		switch checkCounts := metric.ValidationMethod; checkCounts {
		case "static":
			// will be used in future for debug logging
			// fmt.Println("Static check")
			response, err := metrics.CheckStaticThreshold(currentMetricValues, metric.StaticThreshold)
			if err != nil {
				degradationText := fmt.Sprintf("Metric degradation found: %v\n", err)
				fmt.Printf(degradationText)
				return "", 406, fmt.Errorf(degradationText)
			}
			successText += response
		default:
			// will be used in future for debug logging
			// fmt.Println("Default check")
			if !canCompare {
				return "", 400, fmt.Errorf("No previous metrics to compare against for metric %v", cleanMetricName)
			}
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
