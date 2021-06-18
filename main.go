package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/barrebre/goDynaPerfSignature/datatypes"
	"github.com/barrebre/goDynaPerfSignature/logging"
	"github.com/barrebre/goDynaPerfSignature/performancesignature"
	"github.com/barrebre/goDynaPerfSignature/utils"

	"github.com/gorilla/mux"
)

const badRequestMessage = "Please ensure your request meets the required specifications found in the readme: https://github.com/barrebre/goDynaPerfSignature#required-parameters"

// Create the paths to access the APIs
func main() {
	logging.LogSystem(datatypes.Logging{Message: fmt.Sprintf("Starting goDynaPerfSig version: %v on port 8080", utils.GetAppVersion())})

	// Get config
	config := utils.GetConfig()

	// Set up server
	var wait time.Duration
	r := mux.NewRouter()

	r.HandleFunc("/performanceSignature", func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			errMessage := fmt.Sprintf("Couldn't parse the body of the request. Error was: %v.", err.Error())
			logging.LogError(datatypes.Logging{Message: errMessage})
			response := datatypes.PerformanceSignatureReturn{
				Error:    true,
				Response: []string{errMessage, badRequestMessage},
			}
			utils.WriteResponse(w, response, datatypes.PerformanceSignature{})
			return
		}

		// Pull out and verify the provided params
		ps, err := performancesignature.ReadAndValidateParams(b, config)
		if err != nil {
			errMessage := fmt.Sprintf("Could not ReadAndValidateParams. Error was: %v.", err.Error())
			logging.LogError(datatypes.Logging{Message: errMessage})
			response := datatypes.PerformanceSignatureReturn{
				Error:    true,
				Response: []string{errMessage, badRequestMessage},
			}
			utils.WriteResponse(w, response, datatypes.PerformanceSignature{})
			return
		}

		// Perform the performance signature
		response := performancesignature.ProcessRequest(w, r, ps)

		utils.WriteResponse(w, response, ps)
	})

	srv := &http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logging.LogSystem(datatypes.Logging{Message: err.Error()})
		}
	}()

	logging.LogSystem(datatypes.Logging{Message: "Finished starting goDynaPerfSig"})

	// Make a channel to wait for an OS shutdown. This helps us keep the app running until ctrl+c
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)
	logging.LogInfo(datatypes.Logging{Message: "Shutting down goDynaPerfSignature."})
	os.Exit(0)
}
