package metrics

import (
	"fmt"
	"math"
)

// CompareMetrics compares the metrics from the current and previous timeframe
func CompareMetrics(curr float64, prev float64, metric string) (string, error) {
	delta := curr - prev

	if delta > 0 {
		errorMessage := fmt.Sprintf("FAIL - %v had a degradation of %.2f, from %.2f to %.2f", metric, delta, prev, curr)
		return "", fmt.Errorf(errorMessage)
	}

	successResponse := fmt.Sprintf("PASS - Successful deploy! Improvement of %.2f, from %.2f to %.2f", math.Abs(delta), prev, curr)
	return successResponse, nil
}

// CheckRelativeThreshold compares the metrics from the current and previous timeframe
func CheckRelativeThreshold(curr float64, prev float64, rel float64, metric string) (string, error) {
	delta := curr - prev
	relDiff := delta - rel

	// If the difference including the threshold is still negative, it's a failure
	if relDiff > 0 {
		errorMessage := fmt.Sprintf("FAIL - %v did not meet the relative threshold criteria. the current performance is %.2f, which is not better than the previous value of %.2f plus the relative threshold of %.2f.", metric, curr, prev, rel)
		return "", fmt.Errorf(errorMessage)
	}

	// If the delta is negative, that means there was a performance improvement
	if delta < 0 {
		successResponse := fmt.Sprintf("PASS - %v improvement to %.2f from %.2f. (Difference: %.2f)", metric, curr, prev, delta)
		return successResponse, nil
	}

	// Otherwise, the threshold must've allowed this to pass
	successResponse := fmt.Sprintf("PASS - %v's current value is %.2f, which is passable compared to the previous results (%.2f) plus the tolerance (%.2f).", metric, curr, prev, rel)
	return successResponse, nil
}

// CheckStaticThreshold compares the metrics from the current and previous timeframe
func CheckStaticThreshold(value float64, threshold float64, metric string) (string, error) {
	delta := value - threshold

	if delta > 0 {
		errorMessage := fmt.Sprintf("FAIL - %v was above the static threshold: %.2f, instead of a desired %.2f", metric, value, threshold)
		return "", fmt.Errorf(errorMessage)
	}

	successResponse := fmt.Sprintf("PASS - %v fit static threshold of %.2f with value %.2f.", metric, threshold, value)
	return successResponse, nil
}
