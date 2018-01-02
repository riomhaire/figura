package file

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/riomhaire/figura/usecases"
)

type FileBasedStorage struct {
	Registry *usecases.Registry
}

func NewFileBasedStorage(registry *usecases.Registry) FileBasedStorage {
	v := FileBasedStorage{registry}
	return v
}

func (c FileBasedStorage) Locate(application, filename string) (io.Reader, error) {
	// Check if application directory present within the base directory
	// Check file is present within that directory
	// return reader to that file if present
	// otherwise return error
	filepath := fmt.Sprintf("%s/%s/%s", c.Registry.Configuration.ConfigurationLocation, application, filename)
	file, err := os.Open(filepath)
	return bufio.NewReader(file), err

}
