package bootstrap

import (
	"flag"
	"fmt"
	"os"

	"github.com/gorilla/mux"
	"github.com/riomhaire/figura/frameworks"
	"github.com/riomhaire/figura/frameworks/file"
	"github.com/riomhaire/figura/frameworks/serviceregistry/consulagent"
	"github.com/riomhaire/figura/frameworks/serviceregistry/defaultserviceregistry"
	"github.com/riomhaire/figura/frameworks/web"
	"github.com/riomhaire/figura/usecases"
	"github.com/urfave/negroni"
)

const VERSION = "figura Version 1.1.1"

type Application struct {
	registry *usecases.Registry
	restAPI  *web.RestAPI
}

func (a *Application) Initialize() {
	logger := frameworks.ConsoleLogger{}

	logger.Log("INFO", "Initializing")
	// Create Configuration
	configuration := usecases.Configuration{}
	enableProfiling := flag.Bool("profile", false, "Enable profiling endpoint")

	port := flag.Int("port", 3050, "Port to use")
	configurations := flag.String("configs", "configs", "Directory here configurations stored")

	enableConsul := flag.Bool("consul", false, "Enable consul integration")
	consulHost := flag.String("consulHost", "http://localhost:8500", "Consul Host")

	flag.Parse()
	// Set in config
	configuration.Application = "Figura"
	configuration.Version = VERSION
	configuration.Profiling = *enableProfiling
	configuration.Port = *port
	configuration.ConfigurationLocation = *configurations
	configuration.Consul = *enableConsul
	configuration.ConsulHost = *consulHost
	hostname, _ := os.Hostname()
	configuration.Host = hostname

	registry := &usecases.Registry{}
	registry.Configuration = configuration
	registry.Logger = logger
	registry.ConfigurationStorage = file.NewFileBasedConfigurationStorage(registry)
	registry.ConfigurationReader = usecases.NewConfigurationReader(registry)
	registry.Storage = file.NewFileBasedStorage(registry)
	registry.Authenticator = file.NewFileBasedAuthenticationStorage(registry)

	// Do we need external registry
	if configuration.Consul {
		registry.ExternalServiceRegistry = consulagent.NewConsulServiceRegistry(registry, "/api/v1/configuration", "/api/v1/configuration/health")

	} else {
		registry.ExternalServiceRegistry = defaultserviceregistry.NewDefaultServiceRegistry(registry)
	}

	// Create API
	restAPI := web.NewRestAPI(registry)
	a.restAPI = &restAPI
	a.registry = registry

	mux := mux.NewRouter()
	negroni := negroni.Classic()
	restAPI.Negroni = negroni

	// Add handlers
	mux.HandleFunc("/api/v1/configuration/metrics", restAPI.HandleStatistics)
	mux.HandleFunc("/metrics", restAPI.HandleStatistics)
	mux.HandleFunc("/api/v1/configuration/health", restAPI.HandleHealth)
	mux.HandleFunc("/api/v1/configuration/{application}", restAPI.HandleReadConfig)
	mux.HandleFunc("/api/v1/configuration/{application}/{filename}", restAPI.HandleReadFile)

	// Add Middleware
	negroni.Use(restAPI.Statistics)
	negroni.UseFunc(restAPI.AddWorkerHeader)  // Add which instance
	negroni.UseFunc(restAPI.AddWorkerVersion) // Which version
	negroni.UseFunc(restAPI.AddCoorsHeader)   // Add coors
	negroni.UseHandler(mux)

}

func (a *Application) Run() {
	a.registry.Logger.Log("INFO", fmt.Sprintf("Running %s", a.registry.Configuration.Version))
	a.registry.Logger.Log("INFO", a.registry.Configuration.String())
	// Register with external service if required ... default does nothing
	a.registry.ExternalServiceRegistry.Register()
	// Run on port
	a.restAPI.Negroni.Run(fmt.Sprintf(":%d", a.registry.Configuration.Port))
}

func (a *Application) Stop() {
	a.registry.Logger.Log("INFO", "Shutting Down REST API")
	a.registry.ExternalServiceRegistry.Deregister()
}
