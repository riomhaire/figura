package bootstrap

import (
	"flag"
	"fmt"

	"github.com/gorilla/mux"
	"github.com/riomhaire/figura/frameworks"
	"github.com/riomhaire/figura/frameworks/file"
	"github.com/riomhaire/figura/frameworks/web"
	"github.com/riomhaire/figura/usecases"
	"github.com/urfave/negroni"
)

const VERSION = "figura Version 1.0.2"

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

	flag.Parse()
	// Set in config
	configuration.Application = "Figura"
	configuration.Version = VERSION
	configuration.Profiling = *enableProfiling
	configuration.Port = *port
	configuration.ConfigurationLocation = *configurations

	registry := &usecases.Registry{}
	registry.Configuration = configuration
	registry.Logger = logger
	registry.ConfigurationStorage = file.NewFileBasedConfigurationStorage(registry)
	registry.ConfigurationReader = usecases.NewConfigurationReader(registry)
	registry.Storage = file.NewFileBasedStorage(registry)

	// Create API
	restAPI := web.NewRestAPI(registry)
	a.restAPI = &restAPI
	a.registry = registry

	mux := mux.NewRouter()
	negroni := negroni.Classic()
	restAPI.Negroni = negroni

	// Add handlers
	mux.HandleFunc("/api/v1/configuration/statistics", restAPI.HandleStatistics)
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
	a.restAPI.Negroni.Run(fmt.Sprintf(":%d", a.registry.Configuration.Port))
}
