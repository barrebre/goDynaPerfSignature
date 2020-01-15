package utils

import (
	"fmt"
	"net/http"
	"os"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
)

// GetConfig retrives the config from the env
func GetConfig() (datatypes.Config, error) {
	server := os.Getenv("DT_SERVER")
	if server == "" {
		fmt.Printf("A Dynatrace server was not provided. Requests will not work unless a DT_SERVER is given in the POST body.\n")
	} else {
		fmt.Printf("Loaded default Dynatrace Server: %v. This can be overridden with any API POST.\n", server)
	}

	env := os.Getenv("DT_ENV")
	if env == "" {
		fmt.Printf("A Dynatrace environment was not provided. If your tenant has multiple environments, you will need to include this.\n")
	} else {
		fmt.Printf("Loaded default Dynatrace Env: %v. This can be overridden with any API POST.\n", env)
	}

	config := datatypes.Config{
		Env:    env,
		Server: server,
	}
	return config, nil
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
