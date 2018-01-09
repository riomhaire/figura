package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/riomhaire/figura/entities"
)

// HandleReadConfig - REST API endpoint for reading configuration file
func (rest *RestAPI) HandleReadConfig(w http.ResponseWriter, r *http.Request) {
	application := mux.Vars(r)["application"]

	authorization := r.Header.Get("Authorization")
	response := rest.Registry.ConfigurationReader.Lookup([]byte(authorization), application)
	code := http.StatusNotImplemented
	data := []byte("Internal Error")

	switch response.ResultType {
	case entities.NoError:
		code = http.StatusOK
		data = response.Data
	case entities.UnknownApplication:
		code = http.StatusNotFound
		data = []byte(response.Message)
	case entities.AuthenticationError:
		code = http.StatusUnauthorized
		data = []byte(response.Message)
	case entities.NotImplimentedError:
		code = http.StatusNotImplemented
		data = []byte(response.Message)
	}

	// Write actual response to output stream
	w.WriteHeader(code)
	w.Write(data)
	if code != http.StatusOK {
		msg := fmt.Sprintf("App Error %v -> Http %v : %v", response.ResultType, code, response.Message)
		rest.Registry.Logger.Log("ERROR", msg)
	}
}
