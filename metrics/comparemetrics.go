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

	successResponse := fmt.Sprintf("PASS - %v had an improvement of %.2f, from %.2f to %.2f", metric, math.Abs(delta), prev, curr)
	return successResponse, nil
}

// CheckRelativeThreshold compares the metrics from the current and previous timeframe
func CheckRelativeThreshold(curr float64, prev float64, rel float64, metric string) (string, error) {
	delta := curr - prev
	relDiff := delta - rel

	// If the difference including the threshold is still negative, it's a failure
	if relDiff > 0 {
		errorMessage := fmt.Sprintf("FAIL - %v did not meet the relative threshold criteria. The current performance is %.2f, which is not better than the previous value (%.2f) plus the relative threshold (%.2f).", metric, curr, prev, rel)
		return "", fmt.Errorf(errorMessage)
	}

	// If the delta is negative, that means there was a performance improvement
	if delta < 0 {
		successResponse := fmt.Sprintf("PASS - %v had an improvement of %.2f, from %.2f to %.2f", metric, math.Abs(delta), prev, curr)
		return successResponse, nil
	}

	// Otherwise, the threshold must've allowed this to pass
	successResponse := fmt.Sprintf("PASS - %v's current value is %.2f, which is passable compared to the previous results (%.2f) plus the tolerance (%.2f).", metric, curr, prev, rel)
	return successResponse, nil
}

// CheckStaticThreshold checks the current value against a static threshold
func CheckStaticThreshold(value float64, threshold float64, metric string) (string, error) {
	delta := value - threshold

	if delta > 0 {
		errorMessage := fmt.Sprintf("FAIL - %v is above the static threshold (%.2f) with a value of %.2f", metric, threshold, value)
		return "", fmt.Errorf(errorMessage)
	}

	successResponse := fmt.Sprintf("PASS - %v is below the static threshold (%.2f) with a value of %.2f.", metric, threshold, value)
	return successResponse, nil
}
