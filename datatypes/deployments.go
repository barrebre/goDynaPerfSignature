package datatypes

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
