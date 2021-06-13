package performancesignature

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/barrebre/goDynaPerfSignature/logging"
)

// ReadAndValidateParams validates the body params sent in the request from the user
func ReadAndValidateParams(b []byte, config datatypes.Config) (datatypes.PerformanceSignature, error) {
	// Read POST body params
	var performanceSignature datatypes.PerformanceSignature
	err := json.Unmarshal(b, &performanceSignature)
	if err != nil {
		return datatypes.PerformanceSignature{}, err
	}

	// Verify all necessary params were sent
	updatedPerformanceSignature, err := checkParams(performanceSignature, config)
	if err != nil {
		return datatypes.PerformanceSignature{}, fmt.Errorf(fmt.Sprintf("checkParams - %v", err.Error()))
	}

	return updatedPerformanceSignature, nil
}

// Check the required body params sent in with the request to ensure we have all the data we need to query Dt
func checkParams(params datatypes.PerformanceSignature, config datatypes.Config) (datatypes.PerformanceSignature, error) {
	// Build out the backend variables starting with what is in the goDynaPerfSignature config
	finalQuery := datatypes.PerformanceSignature{
		APIToken:       config.APIToken,
		DTEnv:          config.Env,
		DTServer:       config.Server,
		EvaluationMins: params.EvaluationMins,
		EventAge:       params.EventAge,
		Metrics:        params.Metrics,
		ServiceID:      params.ServiceID,
	}

	// Take the params that were sent in and apply them over the goDynaPerfSignature config
	applyPostParams(params, &finalQuery)

	// Finally, ensure we have all the params we need
	err := validateParams(finalQuery)
	if err != nil {
		return datatypes.PerformanceSignature{}, fmt.Errorf(fmt.Sprintf("Couldn't validate parameters: %v", err.Error()))
	}

	return finalQuery, nil
}

func calculateAgeEpoch(days int) int {
	logging.LogInfo(datatypes.Logging{Message: fmt.Sprintf("Received eventAge for event. Checking %v days back", days)})
	pastTime := time.Now().AddDate(0, 0, -days)
	return int(pastTime.Unix() * 1000)
}

// This function applies all of the POST parameters over the default config of the goDynaPerfSignature
func applyPostParams(params datatypes.PerformanceSignature, finalQuery *datatypes.PerformanceSignature) {
	if params.APIToken != "" {
		finalQuery.APIToken = params.APIToken
	}

	if params.DTServer != "" {
		finalQuery.DTServer = params.DTServer
	}

	if params.DTEnv != "" {
		finalQuery.DTEnv = params.DTEnv
	}

	if params.EventAge != 0 {
		finalQuery.EventAge = calculateAgeEpoch(params.EventAge)
	}
}

// Ensure there are no missing parameters to perform a request to Dynatrace
func validateParams(finalQuery datatypes.PerformanceSignature) error {
	if finalQuery.APIToken == "" {
		return fmt.Errorf("there is no DT_API_TOKEN env variable configured and no APIToken was passed with the POST")
	}

	if finalQuery.DTServer == "" {
		return fmt.Errorf("there is no DT_SERVER env variable configured and no DTServer was passed with the POST")
	}

	if len(finalQuery.Metrics) == 0 {
		return fmt.Errorf("no Metrics passed with the POST")
	}

	if finalQuery.ServiceID == "" {
		return fmt.Errorf("no ServiceID passed with the POST")
	}

	return nil
}
