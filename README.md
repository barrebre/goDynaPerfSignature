# goDynaPerfSignature
This repo is a standalone Go application which will allow users to query their Dynatrace environments and compare the performance of current code deployments to previous ones.

# Running the Application
This application can be run in Docker or in Golang.

## Running the App using Docker
The Dockerhub for the repo can be found [here](https://hub.docker.com/r/barrebre/go-dyna-perf-signature/tags). You can run this by calling
```
docker run -expose -p 8080:8080 --name go-dyna-perf-signature barrebre/go-dyna-perf-signature
```

### Optional Environment Variables
Here are some environment variables you may want to consider setting:
* **DT_ENV** - The Dynatrace environment to point to, if your tenant has multiple environments
* **DT_SERVER** - The Dynatrace Server to point to (FQDN)
*If you would like to set the DT_ENV or DT_SERVER, you can do so by editing the `./docker_env` file and then calling*
```
docker run -expose -p 8080:8080 --name go-dyna-perf-signature --env-file ./docker_env barrebre/go-dyna-perf-signature
```

## Running Locally via Go
You can run this locally by simply calling `go run .`

# Calling the Application
Once you have an instance of the application running, you'll want to make calls to it.

## Parameters
The required parameters on the POST are:
* **APIToken** - Your Dynatrace API token which has the permission `Access problem and event feed, metrics, and topology`
* **Metrics** - A comma-delimited array of the metrics you'd like to inspect. 
  * **ID** - The ID of the metric - The list of metric IDs can be found from the `Environment API v2` -> `Metrics` -> `GET /metrics/descriptors` API.
    * `builtin:service.response.time:(avg)`
  * (Optional) **ValidationMethod** - The type of validation you'd like to perform. This is currently limited to `static`. If no value, the default is the comparison model using the most recent and last deployments.
    * `static`
  * (Optional) **StaticThreshold** - If you chose the ValidationMethod `static`, you will need to provide the threshold value here. If you do not, the value will default to 0.00.
    * `1.25`
* **ServiceID** - The ID of the Service which you'd like to inspect. This can be found in the UI if you are looking at a Service and pull from its url `id=SERVICE-...`
  * `SERVICE-5D4E743B2BF0CCF5`

### Optional Parameters
* **DTEnv** - The Dynatrace environment to point to, if your tenant has multiple environments
* **DTServer** - The Dynatrace Server to point to (FQDN)

## Example
From another terminal, you can make requests to the app via a curl like this one:

```
curl -v -XPOST -d '{"APIToken":"","Metrics":[{"ID":"builtin:service.response.time:(avg)"},{"ID":"builtin:service.errors.total.rate:(avg)","StaticThreshold":1.0,"ValidationMethod":"static"}],"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}' localhost:8080/performanceSignature
```

# Development
Current development notes

## Building the App using Docker
First, you will need to build the docker image. To do so
1. Run `docker build -t go-dyna-perf-signature .`
1. Then `docker run -expose -p 8080:8080 --name go-dyna-perf-signature --env-file ./docker_env go-dyna-performance-signature`

A clean way to run this in Windows is

```
docker ps -a -q -f name=go-dyna-perf-signature | % { docker stop $_; docker rm $_ }; docker run -expose -p 8080:8080 --name go-dyna-perf-signature --env-file ./docker_env go-dyna-perf-signature
```

## Build One-Liner
```
docker ps -a -q -f name=go-dyna-perf-signature | % { docker stop $_ }; docker ps -a -q -f name=go-dyna-perf-signature | % { docker rm $_ }; docker run -expose -p 8080:8080 --name go-dyna-perf-signature --env-file ./docker_env barrebre/go-dyna-perf-signature
```
