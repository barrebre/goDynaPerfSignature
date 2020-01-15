package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"time"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
)

// CheckStaticThreshold compares the metrics from the current and previous timeframe
func CheckStaticThreshold(metric float64, threshold float64) (string, error) {
	delta := metric - threshold

	if delta > 0 {
		errorMessage := fmt.Sprintf("The metric was above the static threshold: %v, instead of a desired %v", metric, threshold)
		return errorMessage, fmt.Errorf("The metric was above the static threshold: %v, instead of a desired %v", metric, threshold)
	}

	successResponse := fmt.Sprintf("Metric fit static threshold: %v instead of %v.\n", metric, threshold)
	return successResponse, nil
}

// CompareMetrics compares the metrics from the current and previous timeframe
func CompareMetrics(curr float64, prev float64) (string, error) {
	delta := curr - prev

	if delta > 0 {
		errorMessage := fmt.Sprintf("The metric had a degradation from %v to %v", prev, curr)
		return errorMessage, fmt.Errorf("Degradation of %v%%", delta)
	}

	successResponse := fmt.Sprintf("Successful deploy! Improvement by %v.\n", math.Abs(delta))
	return successResponse, nil
}

// GetMetrics retrieves the metrics from both Deployment Event times in Dynatrace
func GetMetrics(ps datatypes.PerformanceSignature, ts []datatypes.Timestamps) (datatypes.ComparisonMetrics, error) {
	// Transform the POSTed metrics into escaped strings
	metricNames := ""
	for _, metric := range ps.Metrics {
		metricNames += metric.ID + ","
	}
	safeMetricNames := url.QueryEscape(metricNames)

	// Get the metrics from the most recent Deployment Event
	metricResponse, err := queryMetrics(ps.DTServer, ps.DTEnv, safeMetricNames, ts[0], ps)
	if err != nil {
		return datatypes.ComparisonMetrics{}, fmt.Errorf("Error querying current metrics from Dynatrace: %v", err)
	}

	// If there were two Deployment Events, get the second set of metrics
	if len(ts) == 2 {
		previousMetricResponse, err := queryMetrics(ps.DTServer, ps.DTEnv, safeMetricNames, ts[1], ps)
		if err != nil {
			return datatypes.ComparisonMetrics{}, fmt.Errorf("Error querying previous metrics from Dynatrace: %v", err)
		}
		var bothMetricSets = datatypes.ComparisonMetrics{
			CurrentMetrics:  metricResponse,
			PreviousMetrics: previousMetricResponse,
		}

		return bothMetricSets, nil
	}

	var metrics = datatypes.ComparisonMetrics{
		CurrentMetrics: metricResponse,
	}

	return metrics, nil
}

// queryMetrics actually performs the HTTP request to Dynatrace to get the metrics
func queryMetrics(server string, env string, safeMetricNames string, ts datatypes.Timestamps, ps datatypes.PerformanceSignature) (datatypes.DynatraceMetricsResponse, error) {
	// Build the URL
	var url string

	// Check if there's a Dynatrace environment specified
	if env == "" {
		url = fmt.Sprintf("https://%v/api/v2/metrics/series/%v?resolution=Inf&from=%v&to=%v&scope=entity(%v)", server, safeMetricNames, ts.StartTime, ts.EndTime, ps.ServiceID)
	} else {
		url = fmt.Sprintf("https://%v/e/%v/api/v2/metrics/series/%v?resolution=Inf&from=%v&to=%v&scope=entity(%v)", server, env, safeMetricNames, ts.StartTime, ts.EndTime, ps.ServiceID)
	}

	// Build the request object
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request handler: %v", err)
		return datatypes.DynatraceMetricsResponse{}, err
	}

	// Set the API token
	apiTokenField := fmt.Sprintf("Api-Token %v", ps.APIToken)
	req.Header.Add("Authorization", apiTokenField)
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Perform the request
	r, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error reading metric data from Dynatrace: %v", err)
		return datatypes.DynatraceMetricsResponse{}, err
	}

	// Check the status code
	if r.StatusCode != 200 {
		fmt.Printf("Invalid status code from Dynatrace: %v.\n", r.StatusCode)
		return datatypes.DynatraceMetricsResponse{}, fmt.Errorf("Invalid status code from Dynatrace: %v", r.StatusCode)
	}

	// Read in the body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	// Try to parse the response into MetricsResponses
	var metricsResponse datatypes.DynatraceMetricsResponse
	err = json.Unmarshal(b, &metricsResponse)
	if err != nil {
		return datatypes.DynatraceMetricsResponse{}, err
	}

	return metricsResponse, nil
}
