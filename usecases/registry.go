package usecases

import (
	"fmt"

	"github.com/riomhaire/figura/frameworks/serviceregistry"
)

// Configuration containing data from the environment which is used to define program behaviour
type Configuration struct {
	Application           string
	Version               string
	Profiling             bool
	Host                  string
	Port                  int
	ConfigurationLocation string

	Consul     bool
	ConsulHost string
	ConsulId   string // ID of this client
}

type Registry struct {
	Configuration           Configuration
	Logger                  Logger
	ConfigurationReader     ConfigurationInteractor
	ConfigurationStorage    ConfigurationStorage
	Storage                 Storage
	Authenticator           TokenInteractor
	ExternalServiceRegistry serviceregistry.ServiceRegistry
}

func (c *Configuration) String() string {
	return fmt.Sprintf("\nCONFIGURATION\n\tApplication : '%v'\n\t       Port : '%v'\n\t   Location : '%v'\n",
		c.Application,
		c.Port,
		c.ConfigurationLocation,
	)
}
