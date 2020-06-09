package utils

import (
	"log"
	"net/http"
	"os"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
)

// GetConfig retrives the config from the env
// TODO: optimize this in the future so it doesn't check the getenv each time
// TODO: change these vars to match the API vars
func GetConfig() datatypes.Config {
	log.Printf("** Reading in goDynaPerfSignature Config")
	server := os.Getenv("DT_SERVER")
	if server == "" {
		log.Printf("* A Dynatrace server was not provided. Requests will not work unless a DT_SERVER is given in the POST body.\n")
	} else {
		log.Printf("* Loaded default DT_SERVER: %v. This can be overridden with any API POST.\n", server)
	}

	env := os.Getenv("DT_ENV")
	if env == "" {
		log.Printf("* A Dynatrace environment was not provided. If your tenant has multiple environments, you will need to include the DT_ENV in the POST body of requests.\n")
	} else {
		log.Printf("* Loaded default DT_ENV: %v. This can be overridden with any API POST.\n", env)
	}

	apiToken := os.Getenv("DT_API_TOKEN")
	if apiToken == "" {
		log.Printf("* A Dynatrace API token was not provided. DT_API_TOKEN must be given with every API POST.\n")
	} else {
		log.Printf("* Loaded default DT_API_TOKEN: %v. This can be overridden with any API POST.\n", apiToken)
	}

	config := datatypes.Config{
		APIToken: apiToken,
		Env:      env,
		Server:   server,
	}
	return config
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
