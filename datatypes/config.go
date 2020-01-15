package datatypes

//// Definitions

// Config contains the config necessary for the app to run
type Config struct {
	Env    string
	Server string
}

//// Example Values
var (
	configuredConfig = Config{
		Server: "1234.live.dynatrace.com",
		Env:    "envSet",
	}
)

//// Example Accessors

// GetConfiguredConfig returns a fully-configured config
func GetConfiguredConfig() Config {
	return configuredConfig
}
