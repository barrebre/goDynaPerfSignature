# goDynaPerfSignature
goDynaPerfSignature is an Automated Quality Gate for Dynatrace. It is a standalone Go application which can query Dynatrace environments and compare Service metrics.

[Deployment Events](https://www.dynatrace.com/support/help/shortlink/event-types-info#deployment) must be pushed to Dynatrace for goDynaPerfSignature to know when to evaluated metrics.

This application:
1. Queries Dynatrace for Deployment Events pushed to the provided `ServiceID`
  * If there are no Deployment Events, goDynaPerfSignature auto-passes
2. Queries Dynatrace for the metrics of the provided `ServiceID` when there were Deployment Events
3. Performs the provided `ValidationMethod`
  * If there's only one Deployment Event, goDynaPerfSignature can only use the `StaticThreshold` validation
  * If there is more than one Deployment Event, goDynaPerfSignature will evaluate any of the available `ValidationMethod`s on the last two most recent Deployment Events' timeframes
4. Returns a response code based on the evaluation results

# Running goDynaPerfSignature
The Dockerhub for the repo can be found [here](https://hub.docker.com/r/barrebre/go-dyna-perf-signature/tags). You can run this by calling:
```
docker run -expose -p 8080:8080 --name go-dyna-perf-signature barrebre/go-dyna-perf-signature
```

### Optional Environment Variables
The following parameters can be set at application startup:
* **DT_API_TOKEN** - Your Dynatrace API token which has the permission `Access problem and event feed, metrics, and topology`. By providing the DT_API_TOKEN at startup, requests to goDynaPerfSignature will use the provided value by default. This can be overwritten with any request by providing the `APIToken` in the payload
* **DT_ENV** - The Dynatrace environment to query. Use this only if your tenant has multiple environments. *Ex*:`https://{DT_SERVER}/e/{DT_ENV}/`. This can be overwritten with any request by providing the `APIToken` in the payload
* **DT_SERVER** - The Dynatrace Server to point to (FQDN). *Ex*: `https://{DT_SERVER}.live.dynatrace.com`. This can be overwritten with any request by providing the `APIToken` in the payload
* **LOG_LEVEL** - The logging level (the default is `ERROR`, so only errors will be listed). For greater verbosity, use `INFO` or `DEBUG`

To start with any of these parameters, edit the `docker_env` file and then run:
```
docker run -expose -p 8080:8080 --name go-dyna-perf-signature --env-file ./docker_env barrebre/go-dyna-perf-signature
```

# Calling the Application
Below are the required parameters to query goDynaPerfSignature:

## Required Parameters
* **APIToken** - Your Dynatrace API token which has the permission `Access problem and event feed, metrics, and topology`. This is not actually required if goDynaPerfSignature is started with a `DT_API_TOKEN`
* **DTServer** - The Dynatrace Server to point to (FQDN). *Ex*: `haq1234.live.dynatrace.com`. This is not actually required if goDynaPerfSignature is started with a `DT_SERVER`
* **Metrics** - A comma-delimited array of the metrics you'd like to inspect. 
  * **ID** - The ID of the metric - The list of metric IDs can be found from the `Environment API v2` -> `Metrics` -> `GET /metrics/descriptors` API.
    * `builtin:service.response.time:(avg)`
  * (Optional) **ValidationMethod** - The type of validation you'd like to perform. If no value, the default is the comparison model using the most recent and last deployments. The other options are:
    * `relative` - If you are willing to have some amount of degradation, you can provide a RelativeThreshold for leniancy in the comparison
    * `static` - If you want to use a static hard-corded threshold
  * (Optional) **RelativeThreshold** - If you chose the ValidationMethod `relative`, you will need to provide the threshold value here. If you do not, the value will default to 0.00.
  * (Optional) **StaticThreshold** - If you chose the ValidationMethod `static`, you will need to provide the threshold value here. If you do not, the value will default to 0.00.
    * `1.25`
* **ServiceID** - The ID of the Service which you'd like to inspect. This can be found in the UI if you are looking at a Service and pull from its url `id=SERVICE-...`
  * `SERVICE-5D4E743B2BF0CCF5`

## Optional Parameters
* **DTEnv** - The Dynatrace environment to query. Use this only if your tenant has multiple environments. *Ex*:`https://{DT_SERVER}/e/{DT_ENV}/`
* **EvaluationMins** - If you would rather provide an evaluation timeframe than use the duration of Deployment Events, provide a number of minutes in this field. goDynaPerfSignature will evaluate metrics from the beginning of the discovered Deployment Events for the EvaluationMinutes duration. *Ex*: `5`

## Examples
This example queries two different metrics and also passes along a blank APIToken:
```
curl -v -XPOST -d '{"APIToken":"","Metrics":[{"ID":"builtin:service.response.time:(avg)","RelativeThreshold":1.0,"ValidationMethod":"relative"},{"ID":"builtin:service.errors.total.rate:(avg)","StaticThreshold":1.0,"ValidationMethod":"static"}],"ServiceID":"SERVICE-5D4E743B2BF0CCF5"}' localhost:8080/performanceSignature
```

This example queries for a percentile and does not provide the APIToken. This call will only work if goDynaPerfSignature is started with an DT_API_TOKEN configured:
```
curl -v -XPOST -d '{"EvaluationMins":5,"Metrics":[{"ID":"builtin:service.response.time:(percentile(90))"}],"ServiceID":"SERVICE-FFA6FB5E2FA9FFA8"}' localhost:8080/performanceSignature
```