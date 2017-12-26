package usecases

import (
	"github.com/riomhaire/figura/entities"
)

// This cointains the business logic for the read configuration use case

type ConfigurationReader struct {
	Registry *Registry
}

func NewConfigurationReader(registry *Registry) ConfigurationInteractor {
	implimentation := ConfigurationReader{registry}
	return implimentation
}

func (c ConfigurationReader) Lookup(authorization, application string) entities.ApplicationConfiguration {
	// Read from backing store config data
	return c.Registry.ConfigurationStorage.Lookup(application)

}
