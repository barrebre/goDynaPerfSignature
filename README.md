# goDynaPerfSignature

This repo will be a standalone Go application which will allow users to query their Dynatrace environments and compare the performance of current code deployments to previous ones.

**Current State** - This will receive the timestamps of the most recent Dynatrace Deployment Events and compare the Response Time and Error Count of those two deployments to look for regressions.

# Usage
## Requirements
This app requires some ENV vars to be set:
* **DT_SERVER** - The Dynatrace Server to point to
* **DT_ENV** - The Dynatrace environment to point to

## Running Locally via Docker
First, you will need to build the docker image. To do so
1. Run `docker build -t go-dyna-perf-signature .`
1. Then `docker run -expose -p 8080:8080 --name go-dyna-perf-signature --env-file ./docker_env go-dyna-performance-signature`

A clean way to run this in Windows is

```
docker ps -a -q -f name=go-dyna-perf-signature | % { docker stop $_ }; docker ps -a -q -f name=go-dyna-perf-signature | % { docker rm $_ }; docker run -expose -p 8080:8080 --name go-dyna-perf-signature --env-file ./docker_env go-dyna-perf-signature
```

## Running Locally via Go
You can test this locally by simply calling `go run .`

## Interacting with the API
The required parameters on the POST are:
* **APIToken** - Your Dynatrace API token which has the permission `Access problem and event feed, metrics, and topology`
* **MetricIDs** - A comma-delimited array of the metrics you'd like to compare, which can be found from the `Environment API v2` -> `Metrics` -> `GET /metrics/descriptors` API.
  * **Important** - This is not actually implemented at this point. Unfortunately, `["builtin:service.response.time:(avg)","builtin:service.errors.total.rate:(avg)"]` has to be passed through the curl command and it's the only valid operator. This will be fixed soon.
* **ServiceID** - The ID of the Service which you'd like to inspect. This can be found in the UI if you are looking at a Service and pull from its url `id=SERVICE-...`

### Example
From another terminal, you can make requests to the app via a curl like this one:

```
curl -v -XPOST -d '{"APIToken":"S2pMHW_FSlma-PPJIj3l5","MetricIDs":["builtin:service.response.time:(avg)","builtin:service.errors.total.rate:(avg)"],"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}' localhost:8080/performanceSignature
```

## Development Building One-Liner
```
docker build -t barrebre/go-dyna-perf-signature:latest .; docker ps -a -q -f name=go-dyna-perf-signature | % { docker stop $_ }; docker ps -a -q -f name=go-dyna-perf-signature | % { docker rm $_ }; docker run -expose -p 8080:8080 --name go-dyna-perf-signature --env-file ./docker_env go-dyna-perf-signature
```

# Todo
Testing...I know. Any other feedback can be posted in the Github Issues.

