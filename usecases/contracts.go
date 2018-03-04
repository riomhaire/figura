package usecases

import (
	"io"

	"github.com/riomhaire/figura/entities"
)

// This file contains the various interface contracts used by the system.

// Logger - The logging system interface
type Logger interface {
	Log(level, message string)
}

// TokenInteractor - Interface to validate token for specific application
type TokenInteractor interface {
	Valid(token []byte, application string) bool
}

// ConfigurationInteractor - lookup configuration for application/token
type ConfigurationInteractor interface {
	Lookup(authorization []byte, application string) entities.ApplicationConfiguration
}

// ConfigurationStorage - lookup configuration
type ConfigurationStorage interface {
	Lookup(application string) entities.ApplicationConfiguration
}

type Storage interface {
	Locate(application, filename string) (io.Reader, error)
}
