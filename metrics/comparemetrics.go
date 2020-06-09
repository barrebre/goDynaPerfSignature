package metrics

import (
	"fmt"
	"math"
)

// CompareMetrics compares the metrics from the current and previous timeframe
func CompareMetrics(curr float64, prev float64, metric string) (string, error) {
	delta := curr - prev

	if delta > 0 {
		errorMessage := fmt.Sprintf("%v had a degradation of %.2f, from %v to %v", metric, delta, prev, curr)
		return errorMessage, fmt.Errorf("%v degradation of %v", metric, fmt.Sprintf("%.2f", delta))
	}

	successResponse := fmt.Sprintf("Successful deploy! Improvement by %v.\n", math.Abs(delta))
	return successResponse, nil
}

// CheckRelativeThreshold compares the metrics from the current and previous timeframe
func CheckRelativeThreshold(curr float64, prev float64, rel float64, metric string) (string, error) {
	delta := curr - prev
	relDiff := delta - rel

	// If the difference including the threshold is still negative, it's a failure
	if relDiff > 0 {
		errorMessage := fmt.Sprintf("FAIL - %v did not meet the relative threshold criteria. the current performance is %v, which is not better than the previous value of %v plus the relative threshold of %v.", metric, curr, prev, rel)
		return errorMessage, fmt.Errorf("fail - %v degradation of %v, including (%v) relative threshold", metric, fmt.Sprintf("%.2f", delta), rel)
	}

	// If the delta is negative, that means there was a performance improvement
	if delta < 0 {
		successResponse := fmt.Sprintf("PASS - %v improvement to %v from %v. (Difference: %v)\n", metric, curr, prev, delta)
		return successResponse, nil
	}

	// Otherwise, the threshold must've allowed this to pass
	successResponse := fmt.Sprintf("PASS - %v's current value is %v, which is passable compared to the previous results (%v) plus the tolerance (%v).\n", metric, curr, prev, rel)
	return successResponse, nil
}

// CheckStaticThreshold compares the metrics from the current and previous timeframe
func CheckStaticThreshold(value float64, threshold float64, metric string) (string, error) {
	delta := value - threshold

	if delta > 0 {
		errorMessage := fmt.Sprintf("%v was above the static threshold: %v, instead of a desired %v", metric, value, threshold)
		return errorMessage, fmt.Errorf("%v was above the static threshold: %v, instead of a desired %v", metric, value, threshold)
	}

	successResponse := fmt.Sprintf("%v fit static threshold: %v instead of %v.\n", metric, value, threshold)
	return successResponse, nil
}
