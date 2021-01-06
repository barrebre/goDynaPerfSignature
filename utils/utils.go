package utils

import (
	"fmt"
	"net/http"
	"os"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/barrebre/goDynaPerfSignature/logging"
)

const version = "1.4.5"

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
		logging.LogInfo(datatypes.Logging{Message: fmt.Sprintf("A Dynatrace server was not provided. Requests will not work unless a DT_SERVER is given in the POST body.")})
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
func WriteResponse(w http.ResponseWriter, res interface{}, responseText string, err error, errCode int) {
	if err != nil {
		w.WriteHeader(errCode) // Not acceptable - closest applicable
		w.Write([]byte("There was an error: " + err.Error() + "\n"))
	} else {
		w.Write([]byte(responseText)) // Header(200)
	}
}
