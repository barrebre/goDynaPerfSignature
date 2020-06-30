package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/barrebre/goDynaPerfSignature/logging"
)

// GetMetrics retrieves the metrics from both Deployment Event times in Dynatrace
func GetMetrics(ps datatypes.PerformanceSignature, ts []datatypes.Timestamps) (datatypes.ComparisonMetrics, error) {
	safeMetricNames := escapeMetricNames(ps.Metrics)
	logging.LogDebug(datatypes.Logging{Message: fmt.Sprintf("Escaped safe metric names are: %v", safeMetricNames)})

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

// Transform the POSTed metrics into escaped strings
func escapeMetricNames(metricNames []datatypes.Metric) string {
	safeMetricNames := ""
	for _, metric := range metricNames {
		safeMetricNames += metric.ID + ","
	}
	logging.LogDebug(datatypes.Logging{Message: fmt.Sprintf("Safe metric names are: %v", safeMetricNames)})

	return url.QueryEscape(safeMetricNames)
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
	logging.LogInfo(datatypes.Logging{Message: fmt.Sprintf("Built URL: %v", url)})

	// Build the request object
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Error creating request handler: %v.", err)})
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
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Error reading metric data from Dynatrace: %v", err)})
		return datatypes.DynatraceMetricsResponse{}, err
	}

	// Read in the body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Could not read response body from Dynatrace: %v", err.Error())})
		return datatypes.DynatraceMetricsResponse{}, fmt.Errorf("Could not read response body from Dynatrace: %v", err.Error())
	}

	// Check the status code
	if r.StatusCode != 200 {
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Invalid status code from Dynatrace: %v. Message is '%v'\n", r.StatusCode, string(b))})
		return datatypes.DynatraceMetricsResponse{}, fmt.Errorf("Invalid status code from Dynatrace: %v", r.StatusCode)
	}

	// Try to parse the response into MetricsResponses
	var metricsResponse datatypes.DynatraceMetricsResponse
	err = json.Unmarshal(b, &metricsResponse)
	if err != nil {
		return datatypes.DynatraceMetricsResponse{}, err
	}

	return metricsResponse, nil
}
