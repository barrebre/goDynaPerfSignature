package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/barrebre/goDynaPerfSignature/logging"
)

const version = "1.5.0"

// GetConfig retrives the config from the env
// TODO: optimize this in the future so it doesn't check the getenv each time
// TODO: change these vars to match the API vars
func GetConfig() datatypes.Config {
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel != "" {
		logging.SetLogLevel(logLevel)
		logging.LogSystem(datatypes.Logging{Message: fmt.Sprintf("LOG_LEVEL set to %v from env var", logLevel)})
	} else {
		logging.LogSystem(datatypes.Logging{Message: "LOG_LEVEL not found in ENV. Defaulting to ERROR"})
	}

	logging.LogInfo(datatypes.Logging{Message: "** Reading in goDynaPerfSignature Config"})

	server := os.Getenv("DT_SERVER")
	if server == "" {
		logging.LogInfo(datatypes.Logging{Message: "A Dynatrace server was not provided. Requests will not work unless a DT_SERVER is given in the POST body."})
	} else {
		logging.LogInfo(datatypes.Logging{Message: fmt.Sprintf("Loaded default DT_SERVER: %v. This can be overridden with any API POST", server)})
	}

	env := os.Getenv("DT_ENV")
	if env == "" {
		logging.LogInfo(datatypes.Logging{Message: "A Dynatrace environment was not provided. If your tenant has multiple environments, you will need to include the DT_ENV in the POST body of requests."})
	} else {
		logging.LogInfo(datatypes.Logging{Message: fmt.Sprintf("Loaded default DT_ENV: %v. This can be overridden with any API POST.", env)})
	}

	apiToken := os.Getenv("DT_API_TOKEN")
	if apiToken == "" {
		logging.LogInfo(datatypes.Logging{Message: "A Dynatrace API token was not provided. DT_API_TOKEN must be given with every API POST."})
	} else {
		logging.LogInfo(datatypes.Logging{Message: fmt.Sprintf("Loaded default DT_API_TOKEN: %v. This can be overridden with any API POST.", apiToken)})
	}

	config := datatypes.Config{
		APIToken: apiToken,
		Env:      env,
		Server:   server,
	}
	return config
}

// GetAppVersion returns the tagged version of goDynaPerfSignature
func GetAppVersion() string {
	return version
}

// WriteResponse helps respond to requests
func WriteResponse(w http.ResponseWriter, response datatypes.PerformanceSignatureReturn, ps datatypes.PerformanceSignature) {
	w.Header().Set("Content-Type", "application/json")
	if response.ErrorCode != 0 {
		w.WriteHeader(response.ErrorCode)
	}

	responseJson, err := json.Marshal(response)
	if err != nil {
		logging.LogError(datatypes.Logging{
			Message: fmt.Sprintf("Couldn't marshal json for response. Error: %v.", err),
			PerfSig: ps,
		})

		w.WriteHeader(513)
		marshalErrorJson := datatypes.PerformanceSignatureReturn{
			ErrorCode: 513,
			Error:     fmt.Sprintf("goDynaPerfSignature internal problem sending the response. The message was supposed to be: %v", response.Response),
			Response:  []string{},
		}

		errorJson, err2 := json.Marshal(marshalErrorJson)
		if err2 != nil {
			logging.LogError(datatypes.Logging{
				Message: fmt.Sprintf("Couldn't marshal failure json for response. Error: %v.", err2),
				PerfSig: ps,
			})
		}
		w.Write(errorJson)
		return
	}

	w.Write(responseJson)
}
