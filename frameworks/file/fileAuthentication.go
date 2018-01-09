package file

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/riomhaire/figura/usecases"
)

const (
	keyExtension = "key"
)

// FileBasedAuthenticationStorage - contains info for storage based authentication
type FileBasedAuthenticationStorage struct {
	Registry *usecases.Registry
}

func NewFileBasedAuthenticationStorage(registry *usecases.Registry) FileBasedAuthenticationStorage {
	v := FileBasedAuthenticationStorage{registry}
	return v
}

func (a FileBasedAuthenticationStorage) Valid(token []byte, application string) bool {
	foundKeyfile := false
	keyFile := fmt.Sprintf("%v/%v.%v", a.Registry.Configuration.ConfigurationLocation, application, keyExtension)
	if _, err := os.Stat(keyFile); err == nil {
		// exists
		foundKeyfile = true
	}
	a.Registry.Logger.Log("Info", fmt.Sprintf("Keyfile is '%v' and is present ? - %v", keyFile, foundKeyfile))

	// OK if we have no key file the answer is always true - a match
	if foundKeyfile == false {
		return true
	}

	// Otherwise read contents
	data, err := ioutil.ReadFile(keyFile)

	// Did we error ?
	if err != nil {
		return false
	}

	// final step compare
	return bytes.Compare(token, []byte(data)) == 0
}
