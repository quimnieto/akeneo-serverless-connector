package static

// ConfigurationLoader interface for loading configuration
type ConfigurationLoader interface {
	LoadConfiguration(context string) error
}
