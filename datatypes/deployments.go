package datatypes

//// Definitions

// DeploymentEvent defines the data needed from a Dt Deployment Event
type DeploymentEvent struct {
	StartTime         int64  `json:"startTime"`
	EndTime           int64  `json:"endTime"`
	DeploymentName    string `json:"deploymentName"`
	DeploymentVersion string `json:"deploymentVersion"`
}

// DeploymentEvents is a collection of Deployment Events
type DeploymentEvents struct {
	Events []DeploymentEvent `json:"events"`
}

// Timestamps represents a start and end time for Deployment events
type Timestamps struct {
	StartTime int64
	EndTime   int64
}

//// Example Variables
var (
	multipleEventDeploymentEvent = DeploymentEvents{
		Events: []DeploymentEvent{
			{
				StartTime: 1234,
				EndTime:   2345,
			},
			{
				StartTime: 1234,
				EndTime:   2345,
			},
		},
	}

	multipleTimestamps = []Timestamps{
		{
			StartTime: 1234,
			EndTime:   2345,
		},
		{
			StartTime: 1234,
			EndTime:   2345,
		},
	}

	singleEventDeploymentEvent = DeploymentEvents{
		Events: []DeploymentEvent{
			{
				StartTime: 1234,
				EndTime:   2345,
			},
		},
	}

	singleTimestamp = []Timestamps{
		{
			StartTime: 1234,
			EndTime:   2345,
		},
	}
)

//// Example Accessors

// GetMultipleEventDeploymentEvent returns a DeploymentEvents with a single event
func GetMultipleEventDeploymentEvent() DeploymentEvents {
	return multipleEventDeploymentEvent
}

// GetMultipleTimestamps returns a []Timestamps with a multiple timestamps
func GetMultipleTimestamps() []Timestamps {
	return multipleTimestamps
}

// GetSingleEventDeploymentEvent returns a DeploymentEvents with a single event
func GetSingleEventDeploymentEvent() DeploymentEvents {
	return singleEventDeploymentEvent
}

// GetSingleTimestamps returns a []Timestamps with a single timestamp
func GetSingleTimestamps() []Timestamps {
	return singleTimestamp
}
