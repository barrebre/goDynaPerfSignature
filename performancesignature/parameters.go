package performancesignature

import (
	"encoding/json"
	"fmt"

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
		fmt.Printf("Encountered error at check params - %v.\n", err)
		return datatypes.PerformanceSignature{}, err
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
		Metrics:        params.Metrics,
		ServiceID:      params.ServiceID,
	}

	// Take the params that were sent in and apply them over the goDynaPerfSignature config
	applyPostParams(params, &finalQuery)

	// Finally, ensure we have all the params we need
	err := validateParams(finalQuery)
	if err != nil {
		logging.LogError(datatypes.Logging{Message: fmt.Sprintf("Was not able to validate parameters: %v.\n", err.Error())})
		return datatypes.PerformanceSignature{}, err
	}

	return finalQuery, nil
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
}

// Ensure there are no missing parameters to perform a request to Dynatrace
func validateParams(finalQuery datatypes.PerformanceSignature) error {
	if finalQuery.APIToken == "" {
		return fmt.Errorf("There is no DT_API_TOKEN env variable configured and no APIToken was passed with the POST")
	}

	if finalQuery.DTServer == "" {
		return fmt.Errorf("There is no DT_SERVER env variable configured and no DTServer was passed with the POST")
	}

	if len(finalQuery.Metrics) == 0 {
		return fmt.Errorf("No Metrics passed with the POST")
	}

	if finalQuery.ServiceID == "" {
		return fmt.Errorf("No ServiceID passed with the POST")
	}

	return nil
}
