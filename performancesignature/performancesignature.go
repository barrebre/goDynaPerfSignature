package performancesignature

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/barrebre/goDynaPerfSignature/logging"
	"github.com/barrebre/goDynaPerfSignature/metrics"
)

// ProcessRequest handles requests we receive to /performanceSignature
func ProcessRequest(w http.ResponseWriter, r *http.Request, ps datatypes.PerformanceSignature) datatypes.PerformanceSignatureReturn {
	// Build the HTTP request object with query for deployments
	req, err := buildDeploymentRequest(ps)
	if err != nil {
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Error building deployment request: %v.", err)})
		return datatypes.PerformanceSignatureReturn{
			Error:    true,
			Response: []string{"Internal error building deployment request"},
		}
	}

	// Query Dt for events on the given service
	deploymentEvents, err := getDeploymentEvents(*req)
	if err != nil {
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Encountered error gathering event timestamps: %v.", err)})
		return datatypes.PerformanceSignatureReturn{
			Error:    true,
			Response: []string{"Encountered error gathering event timestamps from Dynatrace"},
		}
	}

	// Parse those events to determine when the timestamps we should inspect are
	timestamps, err := parseDeploymentTimestamps(deploymentEvents, ps.EvaluationMins)
	if err != nil {
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Error parsing deployment timestamps: %v.", err)})
		return datatypes.PerformanceSignatureReturn{
			Error:    true,
			Response: []string{"Error parsing deployment timestamps"},
		}
	}

	if len(timestamps) == 0 {
		return datatypes.PerformanceSignatureReturn{
			Error:    false,
			Pass:     true,
			Response: []string{"No deployment events found. Automatic pass"},
		}
	}

	// Get the requested metrics for the discovered timestamp(s)
	metricsResponse, err := metrics.GetMetrics(ps, timestamps)
	if err != nil {
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Encountered error gathering metrics: %v.", err)})
		return datatypes.PerformanceSignatureReturn{
			Error:    true,
			Response: []string{"Encountered error gathering metrics"},
		}
	}
	logging.LogDebug(datatypes.Logging{Message: fmt.Sprintf("Found metrics:\n%v\n", metricsResponse)})

	// Ensure the gathered metrics are within the expected perfSignature
	response := checkPerfSignature(ps, metricsResponse)
	if response.Error {
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Error occurred when checking performance signature: %v.", response.Error)})
		return datatypes.PerformanceSignatureReturn{
			Error:    true,
			Response: []string{fmt.Sprintf("Error occurred when checking performance signature: %v", err)},
		}
	}

	logging.LogInfo(datatypes.Logging{Message: strings.Join(response.Response, "; ")})
	return response
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

	if ps.EventAge != 0 {
		url += fmt.Sprintf("&from=%v", ps.EventAge)
		logging.LogInfo(datatypes.Logging{Message: fmt.Sprintf("Found EventAge with request - new URL is %v", url)})
	}

	// Build the request object
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Error creating request handler: %v", err)})
		return &http.Request{}, err
	}

	// Add the API token
	apiTokenField := fmt.Sprintf("Api-Token %v", ps.APIToken)
	req.Header.Add("Authorization", apiTokenField)

	return req, nil
}

// For each metric, perform its checks
func checkPerfSignature(performanceSignature datatypes.PerformanceSignature, metricsResponse datatypes.ComparisonMetrics) datatypes.PerformanceSignatureReturn {
	// Create the return object, which defaults to a pass
	result := datatypes.PerformanceSignatureReturn{
		Pass: true,
	}

	if len(metricsResponse.CurrentMetrics.Metrics) < 1 {
		result.Response = []string{"PASS - There were no current metrics returned from Dynatrace"}
		return result
	}

	var cleanMetricName string
	for num, metric := range metricsResponse.CurrentMetrics.Metrics {
		if strings.Contains(metric.MetricId, "percentile") {
			cleanMetricName = metric.MetricId
		} else {
			cleanMetricName = strings.ReplaceAll(metric.MetricId, "(", "")
			cleanMetricName = strings.ReplaceAll(cleanMetricName, ")", "")
		}

		localSig := performanceSignature.PSMetrics[cleanMetricName]

		if len(metric.MetricValues) < 1 {
			return datatypes.PerformanceSignatureReturn{
				Pass:     true,
				Response: []string{fmt.Sprintf("PASS - There were no current metrics returned from Dynatrace for %v", cleanMetricName)},
			}
		}

		if len(metric.MetricValues[0].Values) < 1 {
			return datatypes.PerformanceSignatureReturn{
				Pass:     true,
				Response: []string{fmt.Sprintf("PASS - There were no current metric values returned from Dynatrace for %v", cleanMetricName)},
			}
		}
		currentMetricValues := metric.MetricValues[0].Values[0]

		// This is only an issue if trying a comparison
		canCompare := true
		var previousMetricValues float64
		if len(metricsResponse.PreviousMetrics.Metrics[num].MetricValues) < 1 {
			canCompare = false
		} else {
			if len(metricsResponse.PreviousMetrics.Metrics[num].MetricValues[0].Values) < 1 {
				canCompare = false
			} else {
				previousMetricValues = metricsResponse.PreviousMetrics.Metrics[num].MetricValues[0].Values[0]
			}
		}

		switch checkCounts := localSig.ValidationMethod; checkCounts {
		case "relative":
			logging.LogDebug(datatypes.Logging{Message: "Relative Check"})
			response, err := metrics.CheckRelativeThreshold(currentMetricValues, previousMetricValues, localSig.RelativeThreshold, cleanMetricName)
			if err != nil {
				degradationText := fmt.Sprintf("Metric degradation found: %v", err)
				logging.LogInfo(datatypes.Logging{Message: degradationText})
				result.Response = append(result.Response, degradationText)
				result.Pass = false
			} else {
				result.Response = append(result.Response, response)
			}
		case "static":
			logging.LogDebug(datatypes.Logging{Message: "Static Check"})
			response, err := metrics.CheckStaticThreshold(currentMetricValues, localSig.StaticThreshold, cleanMetricName)
			if err != nil {
				degradationText := fmt.Sprintf("Metric degradation found: %v", err)
				logging.LogInfo(datatypes.Logging{Message: degradationText})
				result.Response = append(result.Response, degradationText)
				result.Pass = false
			} else {
				result.Response = append(result.Response, response)
			}
		default:
			logging.LogDebug(datatypes.Logging{Message: "Default Check"})
			if !canCompare {
				degradationText := fmt.Sprintf("No previous metrics to compare against for metric %v", cleanMetricName)
				result.Response = append(result.Response, degradationText)
			} else {
				response, err := metrics.CompareMetrics(currentMetricValues, previousMetricValues, cleanMetricName)
				if err != nil {
					degradationText := fmt.Sprintf("Metric degradation found: %v", err)
					logging.LogInfo(datatypes.Logging{Message: degradationText})
					result.Response = append(result.Response, degradationText)
					result.Pass = false
				} else {
					result.Response = append(result.Response, response)
				}
			}
		}
	}
	return result
}
