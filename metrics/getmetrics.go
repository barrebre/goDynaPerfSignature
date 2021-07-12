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
	metricString := createMetricString(ps.PSMetrics)
	logging.LogDebug(datatypes.Logging{Message: fmt.Sprintf("Escaped safe metric names are: %v", metricString)})

	// Get the metrics from the most recent Deployment Event
	metricResponse, err := queryMetrics(ps.DTServer, ps.DTEnv, metricString, ts[0], ps)
	if err != nil {
		return datatypes.ComparisonMetrics{}, fmt.Errorf("error querying current metrics from Dynatrace: %v", err)
	}

	// If there were two Deployment Events, get the second set of metrics
	if len(ts) == 2 {
		previousMetricResponse, err := queryMetrics(ps.DTServer, ps.DTEnv, metricString, ts[1], ps)
		if err != nil {
			return datatypes.ComparisonMetrics{}, fmt.Errorf("error querying previous metrics from Dynatrace: %v", err)
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
func createMetricString(metricNames map[string]datatypes.PSMetric) string {
	metricString := ""
	for name := range metricNames {
		metricString += name + ","
	}
	logging.LogDebug(datatypes.Logging{Message: fmt.Sprintf("Safe metric names are: %v", metricString)})

	return metricString
}

// queryMetrics actually performs the HTTP request to Dynatrace to get the metrics
func queryMetrics(server string, env string, metricString string, ts datatypes.Timestamps, ps datatypes.PerformanceSignature) (datatypes.DynatraceMetricsResponse, error) {
	url := buildMetricsQueryURL(server, env, metricString, ts, ps)

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
		return datatypes.DynatraceMetricsResponse{}, fmt.Errorf("could not read response body from Dynatrace: %v", err.Error())
	}
	logging.LogDebug(datatypes.Logging{Message: fmt.Sprintf("Full response from Dynatrace is: %v.", string(b))})

	// Check the status code
	if r.StatusCode != 200 {
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Invalid status code from Dynatrace: %v. Message is '%v'", r.StatusCode, string(b))})
		return datatypes.DynatraceMetricsResponse{}, fmt.Errorf("invalid status code from Dynatrace: %v", r.StatusCode)
	}

	// Try to parse the response into MetricsResponses
	var metricsResponse datatypes.DynatraceMetricsResponse
	err = json.Unmarshal(b, &metricsResponse)
	if err != nil {
		return datatypes.DynatraceMetricsResponse{}, err
	}

	return metricsResponse, nil
}

// buildMetricsQueryURL takes all required params and build the URL which will be queried
func buildMetricsQueryURL(server string, env string, metricString string, ts datatypes.Timestamps, ps datatypes.PerformanceSignature) string {
	newURL := url.URL{
		Scheme: "https",
		Host:   server,
	}

	q := newURL.Query()
	q.Set("metricSelector", metricString)
	q.Set("resolution", "Inf")
	q.Set("from", fmt.Sprint(ts.StartTime))
	q.Set("to", fmt.Sprint(ts.EndTime))
	q.Set("entitySelector", fmt.Sprintf("entityId(\"%v\")", ps.ServiceID))
	newURL.RawQuery = q.Encode()

	// Check if there's a Dynatrace environment specified
	if env == "" {
		newURL.Path = "api/v2/metrics/query"
	} else {
		newURL.Path = fmt.Sprintf("/e/%v/api/v2/metrics/query", env)
	}
	logging.LogInfo(datatypes.Logging{Message: fmt.Sprintf("Built URL: %v", newURL.String())})

	return newURL.String()
}
