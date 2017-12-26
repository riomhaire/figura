package web

import (
	"github.com/riomhaire/figura/usecases"
	"github.com/thoas/stats"
	"github.com/urfave/negroni"
)

const HttpSuccess = 200
const HttpUnknown = 404
const HttpUnAuthorized = 403
const HttpApplicationFailure = 500

var bearerPrefix = "bearer "

type RestAPI struct {
	Registry   *usecases.Registry
	Statistics *stats.Stats
	Negroni    *negroni.Negroni
}

func NewRestAPI(registry *usecases.Registry) RestAPI {
	api := RestAPI{}
	api.Registry = registry
	api.Statistics = stats.New()

	return api
}
