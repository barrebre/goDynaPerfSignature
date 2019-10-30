# goDynaPerfSignature

This repo will be a standalone Go application which will allow users to query their Dynatrace environments and compare the performance of current code deployments to previous ones.

**Current State** - This will receive the timestamps of the most recent Dynatrace Deployment Events.

# Usage
## Requirements
This app requires some ENV vars to be set:
* **DT_SERVER** - The Dynatrace Server to point to
* **DT_ENV** - The Dynatrace environment to point to

## Running Locally via Docker
First, you will need to build the docker image. To do so
1. Run `docker build -t go-dyna-perf-signature .`
1. Then `docker run --env-file ./docker_env go-dyna-perf-signature` 

## Running Locally via Go
You can test this locally by simply calling `go run .`

### Example CURL command
From another terminal, you can make requests to the app via a curl like this one: 
```curl -v -XPOST -d '{"APIToken":"S2pMHW_FSlma-PPJIj3l5","MetricIDs":["builtin:service.response.time","builtin:service.errors.total.rate"],"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}' localhost:8080/performanceSignature```

# Todo
Testing...I know