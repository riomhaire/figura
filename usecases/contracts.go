package usecases

import "github.com/riomhaire/figura/entities"

// This file contains the various interface contracts used by the system.

type Logger interface {
	Log(level, message string)
}

type TokenInteractor interface {
	Validate(application string) (string, error)
}

type ConfigurationInteractor interface {
	Lookup(authorization, application string) entities.ApplicationConfiguration
}

type ConfigurationStorage interface {
	Lookup(application string) entities.ApplicationConfiguration
}
