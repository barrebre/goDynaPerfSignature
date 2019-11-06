package utils

import (
	"fmt"
	"net/http"
	"os"

	"barrebre/goDynaPerfSignature/datatypes"
)

// GetConfig retrives the config from the env
func GetConfig() (datatypes.Config, error) {
	server := os.Getenv("DT_SERVER")
	if server == "" {
		return datatypes.Config{}, fmt.Errorf("Error finding the DT_SERVER in the env")
	}
	fmt.Printf("Successfully loaded DT_SERVER: %v.\n", server)

	env := os.Getenv("DT_ENV")
	if env == "" {
		fmt.Printf("A Dynatrace environment was not provided. If your tenant has multiple environments, you will need to include this.\n")
	} else {
		fmt.Printf("Successfully loaded DT_ENV: %v.\n", env)
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
