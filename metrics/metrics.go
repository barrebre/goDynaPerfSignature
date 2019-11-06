package metrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"barrebre/goDynaPerfSignature/datatypes"

	"github.com/davecgh/go-spew/spew"
)

// GetMetrics retrieves the metrics from both Deployment Event times in Dynatrace
// TODO: Make this work with more than the 2 defined metrics
func GetMetrics(config datatypes.Config, ps datatypes.PerformanceSignature, ts []datatypes.Timestamps) (datatypes.ComparisonMetrics, error) {
	// Transform the POSTed metrics into escaped strings
	metricNames := fmt.Sprintf(strings.Join(ps.MetricIDs, ","))
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
		return datatypes.ComparisonMetrics{}, fmt.Errorf("Invalid status code from Dynatrace: %v - ", r.StatusCode)
	}

	// Read in the body
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	// Try to parse the response into MetricsResponses
	var currentMetricsResponse datatypes.DynatraceMetricsResponse
	err = json.Unmarshal(b, &currentMetricsResponse)
	if err != nil {
		return datatypes.ComparisonMetrics{}, err
	}
	numberOfMetrics := len(currentMetricsResponse.Metrics.BuiltinServiceErrorsTotalRate.MetricValues)
	if numberOfMetrics < 1 {
		return datatypes.ComparisonMetrics{}, fmt.Errorf("There were no Error Rate data points found in the current metrics: %v", spew.Sdump(currentMetricsResponse.Metrics))
	}
	numberOfMetrics = len(currentMetricsResponse.Metrics.BuiltinServiceResponseTime.MetricValues)
	if numberOfMetrics < 1 {
		return datatypes.ComparisonMetrics{}, fmt.Errorf("There were no Response Time data points found in the current metrics: %v", spew.Sdump(currentMetricsResponse.Metrics))
	}

	// Get the second set of metrics
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
		return datatypes.ComparisonMetrics{}, fmt.Errorf("Invalid status code from Dynatrace: %v - ", r.StatusCode)
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
	numberOfMetrics = len(previousMetricsResponse.Metrics.BuiltinServiceErrorsTotalRate.MetricValues)
	if numberOfMetrics < 1 {
		return datatypes.ComparisonMetrics{}, fmt.Errorf("There were no Error Rate data points found in the previous event metrics: %v", spew.Sdump(previousMetricsResponse.Metrics))
	}
	numberOfMetrics = len(previousMetricsResponse.Metrics.BuiltinServiceResponseTime.MetricValues)
	if numberOfMetrics < 1 {
		return datatypes.ComparisonMetrics{}, fmt.Errorf("There were no Response Time data points found in the previous event metrics: %v", spew.Sdump(previousMetricsResponse.Metrics))
	}

	var bothMetricSets = datatypes.ComparisonMetrics{
		CurrentMetrics:  currentMetricsResponse.Metrics,
		PreviousMetrics: previousMetricsResponse.Metrics,
	}

	return bothMetricSets, nil
}

// CompareMetrics compares the metrics from the current and previous timeframe
func CompareMetrics(metrics datatypes.ComparisonMetrics) (string, error) {
	currErrRate := metrics.CurrentMetrics.BuiltinServiceErrorsTotalRate.MetricValues[0].Value
	prevErrRate := metrics.PreviousMetrics.BuiltinServiceErrorsTotalRate.MetricValues[0].Value
	errorRateDelta := currErrRate - prevErrRate
	// fmt.Printf("Error rate: %v, from %v to %v.\n", errorRateDelta, prevErrRate, currErrRate)

	if errorRateDelta > 0 {
		errorMessage := fmt.Sprintf("The Error Rate increased since the previous test by %v, from %v to %v", errorRateDelta, prevErrRate, currErrRate)
		return errorMessage, fmt.Errorf("Error rate increase of %v%%", errorRateDelta)
	}

	currResTime := metrics.CurrentMetrics.BuiltinServiceResponseTime.MetricValues[0].Value
	prevResTime := metrics.PreviousMetrics.BuiltinServiceResponseTime.MetricValues[0].Value
	responseTimeDelta := (currResTime - prevResTime) / 1000
	// fmt.Printf("RT: %v, from %v to %v.\n", responseTimeDelta, prevResTime, currResTime)

	if responseTimeDelta > 0 {
		rtMessage := fmt.Sprintf("The Response Time increased since the previous test by %v, from %v to %v", responseTimeDelta, prevResTime, currResTime)
		return rtMessage, fmt.Errorf("Response time increase of %vms", responseTimeDelta)
	}

	successResponse := fmt.Sprintf("Successful deploy! RT decreased by %vms and error rate decreased by %v%%.\n", math.Abs(responseTimeDelta), errorRateDelta)
	return successResponse, nil
}
