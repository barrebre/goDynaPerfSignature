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

// GetMetrics retrieves the metrics from both Deployment Event times in Dynatrace
// TODO: Make this work with more than the 2 defined metrics
func GetMetrics(config datatypes.Config, ps datatypes.PerformanceSignature, ts []datatypes.Timestamps) (datatypes.ComparisonMetrics, error) {
	// Transform the POSTed metrics into escaped strings
	metricNames := ""
	for _, metric := range ps.Metrics {
		metricNames += metric.ID + ","
	}
	safeMetricNames := url.QueryEscape(metricNames)

	// Build the URL
	var url string

	if config.Env == "" {
		url = fmt.Sprintf("https://%v/api/v2/metrics/series/%v?resolution=Inf&from=%v&to=%v&scope=entity(%v)", config.Server, safeMetricNames, ts[0].StartTime, ts[0].EndTime, ps.ServiceID)
	} else {
		url = fmt.Sprintf("https://%v/e/%v/api/v2/metrics/series/%v?resolution=Inf&from=%v&to=%v&scope=entity(%v)", config.Server, config.Env, safeMetricNames, ts[0].StartTime, ts[0].EndTime, ps.ServiceID)
	}
	// fmt.Printf("Made URL: %v\n", url)

	// Build the request object
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request handler: %v", err)
		return datatypes.ComparisonMetrics{}, err
	}

	apiTokenField := fmt.Sprintf("Api-Token %v", ps.APIToken)
	req.Header.Add("Authorization", apiTokenField)
	client := &http.Client{
		Timeout: time.Second * 10,
	}

	// Perform the request
	r, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error reading metric data from Dynatrace: %v", err)
		return datatypes.ComparisonMetrics{}, err
	}
	// Check the status code
	if r.StatusCode != 200 {
		fmt.Printf("Invalid status code from Dynatrace: %v.\n", r.StatusCode)
		return datatypes.ComparisonMetrics{}, fmt.Errorf("Invalid status code from Dynatrace: %v", r.StatusCode)
	}

	// Read in the body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	// fmt.Println(string(b))

	// Try to parse the response into MetricsResponses
	var currentMetricsResponse datatypes.DynatraceMetricsResponse
	err = json.Unmarshal(b, &currentMetricsResponse)
	if err != nil {
		return datatypes.ComparisonMetrics{}, err
	}
	// fmt.Printf("Found current metrics: %v\n", currentMetricsResponse)

	var metrics = datatypes.ComparisonMetrics{
		CurrentMetrics: currentMetricsResponse,
	}

	// Get the second set of metrics
	if len(ts) == 2 {
		// Build the URL
		if config.Env == "" {
			url = fmt.Sprintf("https://%v/api/v2/metrics/series/%v?resolution=Inf&from=%v&to=%v&scope=entity(%v)", config.Server, safeMetricNames, ts[1].StartTime, ts[1].EndTime, ps.ServiceID)
		} else {
			url = fmt.Sprintf("https://%v/e/%v/api/v2/metrics/series/%v?resolution=Inf&from=%v&to=%v&scope=entity(%v)", config.Server, config.Env, safeMetricNames, ts[1].StartTime, ts[1].EndTime, ps.ServiceID)
		}
		// fmt.Printf("Made URL: %v\n", url)

		// Build the request object
		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Error creating request handler: %v", err)
			return datatypes.ComparisonMetrics{}, err
		}

		apiTokenField = fmt.Sprintf("Api-Token %v", ps.APIToken)
		req.Header.Add("Authorization", apiTokenField)

		// Perform the request
		r, err = client.Do(req)
		if err != nil {
			fmt.Printf("Error reading metric data from Dynatrace: %v", err)
			return datatypes.ComparisonMetrics{}, err
		}
		// Check the status code
		if r.StatusCode != 200 {
			fmt.Printf("Invalid status code from Dynatrace: %v.\n", r.StatusCode)
			return datatypes.ComparisonMetrics{}, fmt.Errorf("Invalid status code from Dynatrace: %v", r.StatusCode)
		}

		// Read in the body
		b, err = ioutil.ReadAll(r.Body)
		defer r.Body.Close()

		// Try to parse the response into MetricsResponses
		var previousMetricsResponse datatypes.DynatraceMetricsResponse
		err = json.Unmarshal(b, &previousMetricsResponse)
		if err != nil {
			return datatypes.ComparisonMetrics{}, err
		}
		// fmt.Printf("Found previous metrics: %v\n", previousMetricsResponse)

		var bothMetricSets = datatypes.ComparisonMetrics{
			CurrentMetrics:  currentMetricsResponse,
			PreviousMetrics: previousMetricsResponse,
		}

		// fmt.Println(spew.Sdump(bothMetricSets))
		return bothMetricSets, nil
	}
	return metrics, nil
}

// CompareMetrics compares the metrics from the current and previous timeframe
func CompareMetrics(curr float64, prev float64) (string, error) {
	delta := curr - prev

	if delta > 0 {
		errorMessage := fmt.Sprintf("The metric had a degradation from %v to %v", prev, curr)
		return errorMessage, fmt.Errorf("Error rate increase of %v%%", delta)
	}

	successResponse := fmt.Sprintf("Successful deploy! Improvement by %v.\n", math.Abs(delta))
	return successResponse, nil
}

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
