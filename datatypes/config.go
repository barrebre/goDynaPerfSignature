package datatypes

//// Definitions

// Config contains the config necessary for the app to run
type Config struct {
	APIToken string
	Env      string
	Server   string
}

//// Example Values
var (
	configuredConfig = Config{
		APIToken: "aj0aw9efj0a9wejf09awejf",
		Server:   "1234.live.dynatrace.com",
		Env:      "envSet",
	}
)

//// Example Accessors

// GetConfiguredConfig returns a fully-configured config
func GetConfiguredConfig() Config {
	return configuredConfig
}
