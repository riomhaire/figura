package web

import (
	"github.com/riomhaire/figura/usecases"
	"github.com/thoas/stats"
	"github.com/urfave/negroni"
)

var bearerPrefix = "bearer "

// RestAPI - struct which contains info
type RestAPI struct {
	Registry   *usecases.Registry
	Statistics *stats.Stats
	Negroni    *negroni.Negroni
}

// NewRestAPI - Create a rest API pseudo object structure and populate it
func NewRestAPI(registry *usecases.Registry) RestAPI {
	api := RestAPI{}
	api.Registry = registry
	api.Statistics = stats.New()

	return api
}
