package usecases

import (
	"fmt"

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

func (c ConfigurationReader) Lookup(authorization []byte, application string) entities.ApplicationConfiguration {
	// Verify key if present
	c.Registry.Logger.Log("Info", fmt.Sprintf("Looking for keyfile for %v using token %v", application, string(authorization)))
	if c.Registry.Authenticator == nil {
		c.Registry.Logger.Log("Warn", "No Authenticator implementation defined so any key matches")
		return c.Registry.ConfigurationStorage.Lookup(application)
	} else if c.Registry.Authenticator != nil && c.Registry.Authenticator.Valid(authorization, application) {
		// Read from backing store config data
		c.Registry.Logger.Log("Info", "Authenticator implementation defined - and key matches")
		return c.Registry.ConfigurationStorage.Lookup(application)
	} else {
		c.Registry.Logger.Log("Error", "Authenticator implementation defined - and no match")
		// return error
		return entities.ApplicationConfiguration{ResultType: entities.AuthenticationError, Message: "Invalid Credentials"}
	}

}
