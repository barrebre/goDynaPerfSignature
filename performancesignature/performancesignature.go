package performancesignature

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/barrebre/goDynaPerfSignature/metrics"
)

// ProcessRequest handles requests we receive to /performanceSignature
func ProcessRequest(w http.ResponseWriter, r *http.Request, ps datatypes.PerformanceSignature) (string, int, error) {
	// Build the HTTP request object with query for deployments
	req, err := buildDeploymentRequest(ps)
	if err != nil {
		log.Printf("Error building deployment request: %v.\n", err)
		return "", 503, fmt.Errorf("Error building deployment request: %v", err)
	}

	// Query Dt for events on the given service
	deploymentEvents, err := getDeploymentEvents(*req)
	if err != nil {
		log.Printf("Encountered error gathering event timestamps: %v.\n", err)
		return "", 503, fmt.Errorf("Encountered error gathering timestamps: %v", err)
	}

	// Parse those events to determine when the timestamps we should inspect are
	timestamps, err := parseDeploymentTimestamps(deploymentEvents, ps.EvaluationMins)
	if err != nil {
		log.Printf("Error parsing deployment timestamps: %v.\n", err)
		return "", 503, fmt.Errorf("Error parsing deployment timestamps: %v", err)
	}

	if len(timestamps) == 0 {
		return "No deployment events found. Automatic pass\n", 200, nil
	}

	// Get the requested metrics for the discovered timestamp(s)
	metricsResponse, err := metrics.GetMetrics(ps, timestamps)
	if err != nil {
		log.Printf("Encountered error gathering metrics: %v.\n", err)
		return "", 503, fmt.Errorf("Encountered error gathering metrics: %v", err)
	}
	// will be used in future for debug logging
	// log.Printf("Found metrics:\n%v\n", metricsResponse)

	// Ensure the gathered metrics are within the expected perfSignature
	successText, responseCode, err := checkPerfSignature(ps, metricsResponse)
	if err != nil {
		log.Printf("Error occurred when checking performance signature: %v.\n", err)
		return "", responseCode, fmt.Errorf("Error occurred when checking performance signature: %v", err)
	}

	return successText, 0, nil
}

// Builds the request for getting Deployments from Dynatrace
func buildDeploymentRequest(ps datatypes.PerformanceSignature) (*http.Request, error) {
	// Build the URL
	var url string

	if ps.DTEnv == "" {
		url = fmt.Sprintf("https://%v/api/v1/events?eventType=CUSTOM_DEPLOYMENT&entityId=%v", ps.DTServer, ps.ServiceID)
	} else {
		url = fmt.Sprintf("https://%v/e/%v/api/v1/events?eventType=CUSTOM_DEPLOYMENT&entityId=%v", ps.DTServer, ps.DTEnv, ps.ServiceID)
	}

	// Build the request object
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request handler: %v", err)
		return &http.Request{}, err
	}

	// Add the API token
	apiTokenField := fmt.Sprintf("Api-Token %v", ps.APIToken)
	req.Header.Add("Authorization", apiTokenField)

	return req, nil
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
		var previousMetricValues float64
		if len(metricsResponse.PreviousMetrics.Metrics[cleanMetricName].MetricValues) < 1 {
			canCompare = false
		} else {
			previousMetricValues = metricsResponse.PreviousMetrics.Metrics[cleanMetricName].MetricValues[0].Value
		}

		switch checkCounts := metric.ValidationMethod; checkCounts {
		case "static":
			// will be used in future for debug logging
			// fmt.Println("Static check")
			response, err := metrics.CheckStaticThreshold(currentMetricValues, metric.StaticThreshold, cleanMetricName)
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
			response, err := metrics.CompareMetrics(currentMetricValues, previousMetricValues, cleanMetricName)
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
