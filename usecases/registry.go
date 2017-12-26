package usecases

import "fmt"

// Configuration containing data from the environment which is used to define program behaviour
type Configuration struct {
	Application           string
	Version               string
	Profiling             bool
	Port                  int
	ConfigurationLocation string
}

type Registry struct {
	Configuration        Configuration
	Logger               Logger
	ConfigurationReader  ConfigurationInteractor
	ConfigurationStorage ConfigurationStorage
}

func (c *Configuration) String() string {
	return fmt.Sprintf("\nCONFIGURATION\n\tApplication : '%v'\n\t       Port : '%v'\n\t   Location : '%v'\n",
		c.Application,
		c.Port,
		c.ConfigurationLocation,
	)
}
