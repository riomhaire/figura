package file

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/riomhaire/figura/entities"
	"github.com/riomhaire/figura/usecases"
)

const (
	keyExtension = ".key"
)

var validExtensions = []string{"yml", "yaml", "json", "properties"}

type FileBasedConfigurationStorage struct {
	Registry *usecases.Registry
}

func NewFileBasedConfigurationStorage(registry *usecases.Registry) FileBasedConfigurationStorage {
	v := FileBasedConfigurationStorage{registry}
	return v
}

func (c FileBasedConfigurationStorage) Lookup(application string) entities.ApplicationConfiguration {
	foundMatch := false
	matchingFile := ""

	// Step 1 - Find config file
	for _, ext := range validExtensions {
		// Generate prototype filename
		file := fmt.Sprintf("%v/%v.%v", c.Registry.Configuration.ConfigurationLocation, application, ext)
		if _, err := os.Stat(file); err == nil {
			// exists
			matchingFile = file
			foundMatch = true
		}
	}

	// If not found - exit
	if !foundMatch {
		return entities.ApplicationConfiguration{
			ResultType: entities.UnknownApplication,
			Message:    fmt.Sprintf("Application '%v' is not registered within configuration service.", application),
		}

	}

	// OK we have a match - read into data and return
	data, err := ioutil.ReadFile(matchingFile)

	// Did we error ?
	if err != nil {
		return entities.ApplicationConfiguration{
			ResultType: entities.AuthenticationError,
			Message:    err.Error(),
		}

	}
	// No return success
	return entities.ApplicationConfiguration{
		ResultType: entities.NoError,
		Message:    fmt.Sprintf("Read '%v'", matchingFile),
		Data:       data,
	}
}
